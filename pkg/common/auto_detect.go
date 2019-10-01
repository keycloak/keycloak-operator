package common

import (
	"time"

	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	integreatlyv1alpha1 "github.com/integr8ly/grafana-operator/pkg/apis/integreatly/v1alpha1"
	keycloakv1alpha1 "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/k8sutil"
	"k8s.io/client-go/discovery"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
	_ "github.com/openshift/api/route/v1"
)

var logAuto = logf.Log.WithName("autodetect")

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
				logAuto.Info("finished the first auto-detection")
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
}

func (b *Background) detectMonitoringResources() {
	stateManager := GetStateManager()
	// detect the PrometheusRule resource type exist on the cluster
	resourceExists, _ := k8sutil.ResourceExists(b.dc, monitoringv1.SchemeGroupVersion.String(), monitoringv1.PrometheusRuleKind)
	prometheusRuleExistsState, keyExists := stateManager.GetState(monitoringv1.PrometheusRuleKind).(bool)

	// If the resource exists and we have not set the flag. We do not want to set up the watch a second time
	if resourceExists && (!keyExists || (keyExists && !prometheusRuleExistsState)) {
		stateManager.SetState(monitoringv1.PrometheusRuleKind, true)

		err := watchPrometheusRule(b.controller)
		if err != nil {
			stateManager.SetState(monitoringv1.PrometheusRuleKind, false)
		}
		logAuto.Info("PrometheusRule resource type found on cluster. Secondary watch setup")
	}

	// detect the ServiceMonitor resource type exist on the cluster
	resourceExists, _ = k8sutil.ResourceExists(b.dc, monitoringv1.SchemeGroupVersion.String(), monitoringv1.ServiceMonitorsKind)
	serviceMonitorExistsState, keyExists := stateManager.GetState(monitoringv1.ServiceMonitorsKind).(bool)

	// If the resource exists and we have not set the flag
	if resourceExists && (!keyExists || (keyExists && !serviceMonitorExistsState)) {
		stateManager.SetState(monitoringv1.ServiceMonitorsKind, true)

		err := watchServiceMonitor(b.controller)
		if err != nil {
			stateManager.SetState(monitoringv1.ServiceMonitorsKind, false)
		}
		logAuto.Info("ServiceMonitor resource type found on cluster. Secondary watch setup")
	}

	// detect the GrafanaDashboard resource type resourceExists on the cluster
	resourceExists, _ = k8sutil.ResourceExists(b.dc, integreatlyv1alpha1.SchemeGroupVersion.String(), integreatlyv1alpha1.GrafanaDashboardKind)
	GrafanaDashboardExistsState, keyExists := stateManager.GetState(integreatlyv1alpha1.GrafanaDashboardKind).(bool)

	// If the resource exists and we have not set the flag
	if resourceExists && (!keyExists || (keyExists && !GrafanaDashboardExistsState)) {
		stateManager.SetState(integreatlyv1alpha1.GrafanaDashboardKind, true)

		err := watchGrafanaDashboard(b.controller)
		if err != nil {
			stateManager.SetState(integreatlyv1alpha1.GrafanaDashboardKind, false)
		}
		logAuto.Info("GrafanaDashboard resource type found on cluster. Secondary watch setup")
	}
}

// Setup watch for prometheus rule resource if the resource type exists on the cluster
func watchPrometheusRule(c controller.Controller) error {
	// Watch for changes to secondary resource PrometheusRule and requeue the owner Keycloak
	err := c.Watch(&source.Kind{Type: &monitoringv1.PrometheusRule{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &keycloakv1alpha1.Keycloak{},
	})
	if err != nil {
		return err
	}

	return nil
}

// Setup watch for service monitor resource if the resource type exists on the cluster
func watchServiceMonitor(c controller.Controller) error {
	// Watch for changes to secondary resource ServiceMonitor and requeue the owner Keycloak
	err := c.Watch(&source.Kind{Type: &monitoringv1.ServiceMonitor{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &keycloakv1alpha1.Keycloak{},
	})
	if err != nil {
		return err
	}

	return nil
}

// Setup watch for grafana dashboard resource if the resource type exists on the cluster
func watchGrafanaDashboard(c controller.Controller) error {
	err := c.Watch(&source.Kind{Type: &integreatlyv1alpha1.GrafanaDashboard{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &keycloakv1alpha1.Keycloak{},
	})
	if err != nil {
		return err
	}

	return nil
}
