package common

import (
	"time"

	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	grafanav1alpha1 "github.com/integr8ly/grafana-operator/v3/pkg/apis/integreatly/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/k8sutil"
	routev1 "github.com/openshift/api/route/v1"
	"k8s.io/client-go/discovery"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// Background represents a procedure that runs in the background, periodically auto-detecting features
type Background struct {
	dc     discovery.DiscoveryInterface
	ticker *time.Ticker
}

// New creates a new auto-detect runner
func NewAutoDetect(mgr manager.Manager) (*Background, error) {
	dc, err := discovery.NewDiscoveryClientForConfig(mgr.GetConfig())
	if err != nil {
		return nil, err
	}

	return &Background{dc: dc}, nil
}

// Start initializes the auto-detection process that runs in the background
func (b *Background) Start() {
	b.autoDetectCapabilities()
	// periodically attempts to auto detect all the capabilities for this operator
	b.ticker = time.NewTicker(5 * time.Second)

	go func() {
		for range b.ticker.C {
			b.autoDetectCapabilities()
		}
	}()
}

// Stop causes the background process to stop auto detecting capabilities
func (b *Background) Stop() {
	b.ticker.Stop()
}

func (b *Background) autoDetectCapabilities() {
	b.detectOpenshift()
	b.detectMonitoringResources()
	b.detectRoute()
}

func (b *Background) detectRoute() {
	resourceExists, _ := k8sutil.ResourceExists(b.dc, routev1.SchemeGroupVersion.String(), RouteKind)
	if resourceExists {
		// Set state that the Route kind exists. Used to determine when a route or an Ingress should be created
		stateManager := GetStateManager()
		stateManager.SetState(RouteKind, true)
	}
}

func (b *Background) detectMonitoringResources() {
	// detect the PrometheusRule resource type exist on the cluster
	stateManager := GetStateManager()
	resourceExists, _ := k8sutil.ResourceExists(b.dc, monitoringv1.SchemeGroupVersion.String(), monitoringv1.PrometheusRuleKind)
	stateManager.SetState(monitoringv1.PrometheusRuleKind, resourceExists)

	// detect the ServiceMonitor resource type exist on the cluster
	resourceExists, _ = k8sutil.ResourceExists(b.dc, monitoringv1.SchemeGroupVersion.String(), monitoringv1.ServiceMonitorsKind)
	stateManager.SetState(monitoringv1.ServiceMonitorsKind, resourceExists)

	// detect the GrafanaDashboard resource type resourceExists on the cluster
	resourceExists, _ = k8sutil.ResourceExists(b.dc, grafanav1alpha1.SchemeGroupVersion.String(), grafanav1alpha1.GrafanaDashboardKind)
	stateManager.SetState(monitoringv1.ServiceMonitorsKind, resourceExists)
}

func (b *Background) detectOpenshift() {
	apiGroupVersion := "operator.openshift.io/v1"
	kind := OpenShiftAPIServerKind
	stateManager := GetStateManager()
	isOpenshift, _ := k8sutil.ResourceExists(b.dc, apiGroupVersion, kind)
	if isOpenshift {
		// Set state that its Openshift (helps to differentiate between openshift and kubernetes)
		stateManager.SetState(OpenShiftAPIServerKind, true)
	} else {
		stateManager.SetState(OpenShiftAPIServerKind, false)
	}
}
