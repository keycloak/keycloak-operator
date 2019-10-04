package keycloak

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	integreatlyv1alpha1 "github.com/integr8ly/grafana-operator/pkg/apis/integreatly/v1alpha1"
	routev1 "github.com/openshift/api/route/v1"

	kc "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	keycloakv1alpha1 "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/common"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_keycloak")

const (
	RequeueDelaySeconds = 30
	ControllerName      = "keycloak-controller"
)

// Add creates a new Keycloak Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager, autodetectChannel chan schema.GroupVersionKind) error {
	return add(mgr, newReconciler(mgr), autodetectChannel)
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	client := mgr.GetClient()

	return &ReconcileKeycloak{
		client:  client,
		scheme:  mgr.GetScheme(),
		context: ctx,
		cancel:  cancel,
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler, autodetectChannel chan schema.GroupVersionKind) error {
	// Create a new controller
	c, err := controller.New(ControllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Keycloak
	err = c.Watch(&source.Kind{Type: &keycloakv1alpha1.Keycloak{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Setting up a listener for events on the channel from autodetect
	go func() {
		for gvk := range autodetectChannel {
			// Check if this channel event was for the PrometheusRule resource type
			if gvk.String() == monitoringv1.SchemeGroupVersion.WithKind(monitoringv1.PrometheusRuleKind).String() {
				watchSecondaryResource(c, gvk, &monitoringv1.PrometheusRule{}) // nolint
			}

			// Check if this channel event was for the ServiceMonitor resource type
			if gvk.String() == monitoringv1.SchemeGroupVersion.WithKind(monitoringv1.ServiceMonitorsKind).String() {
				watchSecondaryResource(c, gvk, &monitoringv1.ServiceMonitor{}) // nolint
			}

			// Check if this channel event was for the GrafanaDashboard resource type
			if gvk.String() == integreatlyv1alpha1.SchemeGroupVersion.WithKind(integreatlyv1alpha1.GrafanaDashboardKind).String() {
				watchSecondaryResource(c, gvk, &integreatlyv1alpha1.GrafanaDashboard{}) // nolint
			}

			// Check if this channel event was for the Route resource type
			if gvk.String() == routev1.SchemeGroupVersion.WithKind(common.RouteKind).String() {
				_ = watchSecondaryResource(c, gvk, &routev1.Route{})
			}
		}
	}()

	return nil
}

// blank assignment to verify that ReconcileKeycloak implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileKeycloak{}

// ReconcileKeycloak reconciles a Keycloak object
type ReconcileKeycloak struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client  client.Client
	scheme  *runtime.Scheme
	context context.Context
	cancel  context.CancelFunc
}

// Reconcile reads that state of the cluster for a Keycloak object and makes changes based on the state read
// and what is in the Keycloak.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileKeycloak) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Keycloak")

	// Fetch the Keycloak instance
	instance := &keycloakv1alpha1.Keycloak{}
	err := r.client.Get(r.context, request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Read current state
	currentState := common.NewClusterState()
	err = currentState.Read(r.context, instance, r.client)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Get Action to reconcile current state into desired state
	reconciler := NewKeycloakReconciler()
	desiredState := reconciler.Reconcile(currentState, instance)

	// Run the actions to reach the desired state
	actionRunner := common.NewClusterActionRunner(r.context, r.client, r.scheme, instance)
	err = actionRunner.RunAll(desiredState)
	if err != nil {
		return reconcile.Result{}, err
	}

	log.Info("desired cluster state met")
	return reconcile.Result{RequeueAfter: RequeueDelaySeconds * time.Second}, nil
}

func watchSecondaryResource(c controller.Controller, gvk schema.GroupVersionKind, o runtime.Object) error {
	stateManager := common.GetStateManager()
	stateFieldName := getStateFieldName(gvk.Kind)

	// Check to see if the watch exists for this resource type already for this controller, if it does, we return so we don't set up another watch
	watchExists, keyExists := stateManager.GetState(stateFieldName).(bool)
	if keyExists || watchExists {
		return nil
	}

	// Set up the actual watch
	err := c.Watch(&source.Kind{Type: o}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &kc.Keycloak{},
	})

	// Retry on error
	if err != nil {
		log.Error(err, "error creating watch")
		stateManager.SetState(stateFieldName, false)
		return err
	}

	stateManager.SetState(stateFieldName, true)
	log.Info(fmt.Sprintf("Watch created for '%s' resource in '%s'", gvk.Kind, ControllerName))
	return nil
}

func getStateFieldName(kind string) string {
	return ControllerName + "-watch-" + kind
}
