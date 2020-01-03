package common

import (
	"k8s.io/apimachinery/pkg/runtime/schema"

	"time"

	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	grafanav1alpha1 "github.com/integr8ly/grafana-operator/v3/pkg/apis/integreatly/v1alpha1"
	routev1 "github.com/openshift/api/route/v1"
	"github.com/operator-framework/operator-sdk/pkg/k8sutil"
	"k8s.io/client-go/discovery"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// Background represents a procedure that runs in the background, periodically auto-detecting features
type Background struct {
	dc                  discovery.DiscoveryInterface
	ticker              *time.Ticker
	SubscriptionChannel chan schema.GroupVersionKind
}

// New creates a new auto-detect runner
func NewAutoDetect(mgr manager.Manager) (*Background, error) {
	dc, err := discovery.NewDiscoveryClientForConfig(mgr.GetConfig())
	if err != nil {
		return nil, err
	}

	// Create a new channel that GVK type will be sent down
	subChan := make(chan schema.GroupVersionKind, 1)

	return &Background{dc: dc, SubscriptionChannel: subChan}, nil
}

// Start initializes the auto-detection process that runs in the background
func (b *Background) Start() {
	// periodically attempts to auto detect all the capabilities for this operator
	b.ticker = time.NewTicker(5 * time.Second)

	go func() {
		b.autoDetectCapabilities()

		for range b.ticker.C {
			b.autoDetectCapabilities()
		}
	}()
}

// Stop causes the background process to stop auto detecting capabilities
func (b *Background) Stop() {
	b.ticker.Stop()
	close(b.SubscriptionChannel)
}

func (b *Background) autoDetectCapabilities() {
	b.detectMonitoringResources()
	b.detectRoute()
}

func (b *Background) detectRoute() {
	resourceExists, _ := k8sutil.ResourceExists(b.dc, routev1.SchemeGroupVersion.String(), RouteKind)
	if resourceExists {
		// Set state that the Route kind exists. Used to determine when a route or an Ingress should be created
		stateManager := GetStateManager()
		stateManager.SetState(RouteKind, true)

		b.SubscriptionChannel <- routev1.SchemeGroupVersion.WithKind(RouteKind)
	}
}

func (b *Background) detectMonitoringResources() {
	// detect the PrometheusRule resource type exist on the cluster
	resourceExists, _ := k8sutil.ResourceExists(b.dc, monitoringv1.SchemeGroupVersion.String(), monitoringv1.PrometheusRuleKind)
	if resourceExists {
		b.SubscriptionChannel <- monitoringv1.SchemeGroupVersion.WithKind(monitoringv1.PrometheusRuleKind)
	}

	// detect the ServiceMonitor resource type exist on the cluster
	resourceExists, _ = k8sutil.ResourceExists(b.dc, monitoringv1.SchemeGroupVersion.String(), monitoringv1.ServiceMonitorsKind)
	if resourceExists {
		b.SubscriptionChannel <- monitoringv1.SchemeGroupVersion.WithKind(monitoringv1.ServiceMonitorsKind)
	}

	// detect the PodMonitor resource type exist on the cluster
	resourceExists, _ = k8sutil.ResourceExists(b.dc, monitoringv1.SchemeGroupVersion.String(), monitoringv1.PodMonitorsKind)
	if resourceExists {
		b.SubscriptionChannel <- monitoringv1.SchemeGroupVersion.WithKind(monitoringv1.PodMonitorsKind)
	}

	// detect the GrafanaDashboard resource type resourceExists on the cluster
	resourceExists, _ = k8sutil.ResourceExists(b.dc, grafanav1alpha1.SchemeGroupVersion.String(), grafanav1alpha1.GrafanaDashboardKind)
	if resourceExists {
		b.SubscriptionChannel <- grafanav1alpha1.SchemeGroupVersion.WithKind(grafanav1alpha1.GrafanaDashboardKind)
	}
}
