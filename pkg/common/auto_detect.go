package common

import (
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	"time"

	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	i8ly "github.com/integr8ly/grafana-operator/pkg/apis/integreatly/v1alpha1"
	kc "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	routev1 "github.com/openshift/api/route/v1"
	"github.com/operator-framework/operator-sdk/pkg/k8sutil"
	"k8s.io/client-go/discovery"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var logger = logf.Log.WithName("autodetect")

// Route kind is not provided by the openshift api
const (
	RouteKind = "Route"
)

// Background represents a procedure that runs in the background, periodically auto-detecting features
type Background struct {
	client     client.Client
	dc         discovery.DiscoveryInterface
	controller controller.Controller
	ticker     *time.Ticker
}

// New creates a new auto-detect runner
func NewAutoDetect(mgr manager.Manager, c controller.Controller) (*Background, error) {
	dc, err := discovery.NewDiscoveryClientForConfig(mgr.GetConfig())
	if err != nil {
		return nil, err
	}

	return &Background{dc: dc, client: mgr.GetClient(), controller: c}, nil
}

// Start initializes the auto-detection process that runs in the background
func (b *Background) Start() {
	// periodically attempts to auto detect all the capabilities for this operator
	b.ticker = time.NewTicker(5 * time.Second)

	done := make(chan bool)
	go func() {
		b.autoDetectCapabilities()
		done <- true
	}()

	go func() {
		for {
			select {
			case <-done:
				logger.Info("finished the first auto-detection")
			case <-b.ticker.C:
				b.autoDetectCapabilities()
			}
		}
	}()
}

// Stop causes the background process to stop auto detecting capabilities
func (b *Background) Stop() {
	b.ticker.Stop()
}

func (b *Background) autoDetectCapabilities() {
	b.detectMonitoringResources()
	b.detectRoute()
}

func (b *Background) detectRoute() {
	resourceExists, _ := k8sutil.ResourceExists(b.dc, routev1.SchemeGroupVersion.String(), RouteKind)
	b.tryWatch(resourceExists, RouteKind, &routev1.Route{})
}

func (b *Background) detectMonitoringResources() {
	// detect the PrometheusRule resource type exist on the cluster
	resourceExists, _ := k8sutil.ResourceExists(b.dc, monitoringv1.SchemeGroupVersion.String(), monitoringv1.PrometheusRuleKind)
	b.tryWatch(resourceExists, monitoringv1.PrometheusRuleKind, &monitoringv1.PrometheusRule{})

	// detect the ServiceMonitor resource type exist on the cluster
	resourceExists, _ = k8sutil.ResourceExists(b.dc, monitoringv1.SchemeGroupVersion.String(), monitoringv1.ServiceMonitorsKind)
	b.tryWatch(resourceExists, monitoringv1.ServiceMonitorsKind, &monitoringv1.ServiceMonitor{})

	// detect the GrafanaDashboard resource type resourceExists on the cluster
	resourceExists, _ = k8sutil.ResourceExists(b.dc, i8ly.SchemeGroupVersion.String(), i8ly.GrafanaDashboardKind)
	b.tryWatch(resourceExists, i8ly.GrafanaDashboardKind, &i8ly.GrafanaDashboard{})
}

func (b *Background) tryWatch(resourceExists bool, kind string, o runtime.Object) error {
	if !resourceExists {
		return nil
	}

	stateManager := GetStateManager()
	watchExists, keyExists := stateManager.GetState(kind).(bool)

	// If no key esists yet, but the resource exists, set up a watch
	// If not no key exists, but no watch exists yet, set up a watch
	if keyExists == false || watchExists == false {
		// Try to set up the actual watch
		err := b.controller.Watch(&source.Kind{Type: o}, &handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    &kc.Keycloak{},
		})

		// Retry on error
		if err != nil {
			logger.Error(err, "error creating watch")
			stateManager.SetState(kind, false)
			return err
		}

		stateManager.SetState(kind, true)
		logger.Info(fmt.Sprintf("'%s' type exists, watch created", kind))
	}

	return nil
}
