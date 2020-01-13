package keycloak

import (
	"context"
	"fmt"
	"time"

	v1beta12 "k8s.io/api/policy/v1beta1"

	"github.com/keycloak/keycloak-operator/pkg/model"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	grafanav1alpha1 "github.com/integr8ly/grafana-operator/v3/pkg/apis/integreatly/v1alpha1"
	routev1 "github.com/openshift/api/route/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	kc "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	keycloakv1alpha1 "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/common"

	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_keycloak")

const (
	RequeueDelaySeconds      = 30 * time.Second
	RequeueDelayErrorSeconds = 5 * time.Second
	ControllerName           = "keycloak-controller"
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
		client:   client,
		scheme:   mgr.GetScheme(),
		context:  ctx,
		cancel:   cancel,
		recorder: mgr.GetEventRecorderFor(ControllerName),
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

	if err := common.WatchSecondaryResource(c, ControllerName, common.SecretKind, &corev1.Secret{}, &kc.Keycloak{}); err != nil {
		return err
	}

	if err := common.WatchSecondaryResource(c, ControllerName, common.StatefulSetKind, &appsv1.StatefulSet{}, &kc.Keycloak{}); err != nil {
		return err
	}

	if err := common.WatchSecondaryResource(c, ControllerName, common.ServiceKind, &corev1.Service{}, &kc.Keycloak{}); err != nil {
		return err
	}

	if err := common.WatchSecondaryResource(c, ControllerName, common.IngressKind, &v1beta1.Ingress{}, &kc.Keycloak{}); err != nil {
		return err
	}

	if err := common.WatchSecondaryResource(c, ControllerName, common.DeploymentKind, &appsv1.Deployment{}, &kc.Keycloak{}); err != nil {
		return err
	}

	if err := common.WatchSecondaryResource(c, ControllerName, common.PersistentVolumeClaimKind, &corev1.PersistentVolumeClaim{}, &kc.Keycloak{}); err != nil {
		return err
	}

	if err := common.WatchSecondaryResource(c, ControllerName, common.PodDisruptionBudgetKind, &v1beta12.PodDisruptionBudget{}, &kc.Keycloak{}); err != nil {
		return err
	}

	// Setting up a listener for events on the channel from autodetect
	go func() {
		for gvk := range autodetectChannel {
			// Check if this channel event was for the PrometheusRule resource type
			if gvk.String() == monitoringv1.SchemeGroupVersion.WithKind(monitoringv1.PrometheusRuleKind).String() {
				common.WatchSecondaryResource(c, ControllerName, gvk.Kind, &monitoringv1.PrometheusRule{}, &kc.Keycloak{}) // nolint
			}

			// Check if this channel event was for the ServiceMonitor resource type
			if gvk.String() == monitoringv1.SchemeGroupVersion.WithKind(monitoringv1.ServiceMonitorsKind).String() {
				common.WatchSecondaryResource(c, ControllerName, gvk.Kind, &monitoringv1.ServiceMonitor{}, &kc.Keycloak{}) // nolint
			}

			// Check if this channel event was for the PodMonitor resource type
			if gvk.String() == monitoringv1.SchemeGroupVersion.WithKind(monitoringv1.PodMonitorsKind).String() {
				common.WatchSecondaryResource(c, ControllerName, gvk.Kind, &monitoringv1.PodMonitor{}, &kc.Keycloak{}) // nolint
			}

			// Check if this channel event was for the GrafanaDashboard resource type
			if gvk.String() == grafanav1alpha1.SchemeGroupVersion.WithKind(grafanav1alpha1.GrafanaDashboardKind).String() {
				common.WatchSecondaryResource(c, ControllerName, gvk.Kind, &grafanav1alpha1.GrafanaDashboard{}, &kc.Keycloak{}) // nolint
			}

			// Check if this channel event was for the Route resource type
			if gvk.String() == routev1.SchemeGroupVersion.WithKind(common.RouteKind).String() {
				_ = common.WatchSecondaryResource(c, ControllerName, gvk.Kind, &routev1.Route{}, &kc.Keycloak{})
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
	client   client.Client
	scheme   *runtime.Scheme
	context  context.Context
	cancel   context.CancelFunc
	recorder record.EventRecorder
}

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
		return r.ManageError(instance, err)
	}

	// Get Action to reconcile current state into desired state
	reconciler := NewKeycloakReconciler()
	desiredState := reconciler.Reconcile(currentState, instance)

	// Perform migration if needed
	migrator := NewDefaultMigrator()
	desiredState, err = migrator.Migrate(instance, currentState, desiredState)
	if err != nil {
		return r.ManageError(instance, err)
	}

	// Run the actions to reach the desired state
	actionRunner := common.NewClusterActionRunner(r.context, r.client, r.scheme, instance)
	err = actionRunner.RunAll(desiredState)
	if err != nil {
		return r.ManageError(instance, err)
	}

	return r.ManageSuccess(instance, currentState)
}

func (r *ReconcileKeycloak) ManageError(instance *v1alpha1.Keycloak, issue error) (reconcile.Result, error) {
	r.recorder.Event(instance, "Warning", "ProcessingError", issue.Error())

	instance.Status.Message = issue.Error()
	instance.Status.Ready = false
	instance.Status.Phase = v1alpha1.PhaseFailing

	err := r.client.Status().Update(r.context, instance)
	if err != nil {
		log.Error(err, "unable to update status")
	}

	return reconcile.Result{
		RequeueAfter: RequeueDelayErrorSeconds,
		Requeue:      true,
	}, nil
}

func (r *ReconcileKeycloak) ManageSuccess(instance *v1alpha1.Keycloak, currentState *common.ClusterState) (reconcile.Result, error) {
	// Check if the resources are ready
	resourcesReady, err := currentState.IsResourcesReady()
	if err != nil {
		return r.ManageError(instance, err)
	}

	instance.Status.Ready = resourcesReady
	instance.Status.Message = ""

	// If resources are ready and we have not errored before now, we are in a reconciling phase
	if resourcesReady {
		instance.Status.Phase = v1alpha1.PhaseReconciling
	} else {
		instance.Status.Phase = v1alpha1.PhaseInitialising
	}

	// Make this keycloaks url public to allow access via the client
	if currentState.KeycloakRoute != nil && currentState.KeycloakRoute.Spec.Host != "" {
		instance.Status.InternalURL = fmt.Sprintf("https://%v", currentState.KeycloakRoute.Spec.Host)
	} else if currentState.KeycloakService != nil && currentState.KeycloakService.Spec.ClusterIP != "" {
		instance.Status.InternalURL = fmt.Sprintf("https://%v.%v.svc:%v",
			currentState.KeycloakService.Name,
			currentState.KeycloakService.Namespace,
			model.KeycloakServicePort)
	}

	// Let the clients know where the admin credentials are stored
	if currentState.KeycloakAdminSecret != nil {
		instance.Status.CredentialSecret = currentState.KeycloakAdminSecret.Name
	}

	err = r.client.Status().Update(r.context, instance)
	if err != nil {
		log.Error(err, "unable to update status")
		return reconcile.Result{
			RequeueAfter: RequeueDelayErrorSeconds,
			Requeue:      true,
		}, nil
	}

	log.Info("desired cluster state met")
	return reconcile.Result{RequeueAfter: RequeueDelaySeconds}, nil
}
