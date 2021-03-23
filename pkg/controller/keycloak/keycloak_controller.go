package keycloak

import (
	"context"
	"fmt"
	"time"

	"github.com/keycloak/keycloak-operator/version"

	v1beta12 "k8s.io/api/policy/v1beta1"

	"github.com/keycloak/keycloak-operator/pkg/model"

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
	"github.com/pkg/errors"

	"k8s.io/api/extensions/v1beta1"
	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
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
	RequeueDelay      = 30 * time.Second
	RequeueDelayError = 5 * time.Second
	ControllerName    = "keycloak-controller"
)

// Add creates a new Keycloak Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
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
func add(mgr manager.Manager, r reconcile.Reconciler) error {
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

	if err := common.WatchSecondaryResource(c, ControllerName, monitoringv1.PrometheusRuleKind, &monitoringv1.PrometheusRule{}, &kc.Keycloak{}); err != nil {
		return err
	}

	if err := common.WatchSecondaryResource(c, ControllerName, monitoringv1.ServiceMonitorsKind, &monitoringv1.ServiceMonitor{}, &kc.Keycloak{}); err != nil {
		return err
	}

	if err := common.WatchSecondaryResource(c, ControllerName, grafanav1alpha1.GrafanaDashboardKind, &grafanav1alpha1.GrafanaDashboard{}, &kc.Keycloak{}); err != nil {
		return err
	}

	if err := common.WatchSecondaryResource(c, ControllerName, common.RouteKind, &routev1.Route{}, &kc.Keycloak{}); err != nil {
		return err
	}

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
		if kubeerrors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}
	currentState := common.NewClusterState()

	if instance.Spec.Unmanaged {
		return r.ManageSuccess(instance, currentState)
	}

	if instance.Spec.External.Enabled {
		return r.ManageError(instance, errors.Errorf("if external.enabled is true, unmanaged also needs to be true"))
	}

	if instance.Spec.ExternalAccess.Host != "" {
		isOpenshift, _ := common.GetStateManager().GetState(common.OpenShiftAPIServerKind).(bool)
		if isOpenshift {
			return r.ManageError(instance, errors.Errorf("Setting Host in External Access on OpenShift is prohibited"))
		}
	}

	// Read current state
	err = currentState.Read(r.context, instance, r.client)
	if err != nil {
		return r.ManageError(instance, err)
	}

	// Get Action to reconcile current state into desired state
	reconciler := NewKeycloakReconciler()
	desiredState := reconciler.Reconcile(currentState, instance)

	// Perform migration if needed
	migrator, err := GetMigrator(instance)
	if err != nil {
		return r.ManageError(instance, err)
	}
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

	r.setVersion(instance)

	err := r.client.Status().Update(r.context, instance)
	if err != nil {
		log.Error(err, "unable to update status")
	}

	return reconcile.Result{
		RequeueAfter: RequeueDelayError,
		Requeue:      true,
	}, nil
}

func (r *ReconcileKeycloak) ManageSuccess(instance *v1alpha1.Keycloak, currentState *common.ClusterState) (reconcile.Result, error) {
	// Check if the resources are ready
	resourcesReady, err := currentState.IsResourcesReady(instance)
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

	if currentState.KeycloakService != nil && currentState.KeycloakService.Spec.ClusterIP != "" {
		instance.Status.InternalURL = fmt.Sprintf("https://%v.%v.svc:%v",
			currentState.KeycloakService.Name,
			currentState.KeycloakService.Namespace,
			model.KeycloakServicePort)
	}

	if instance.Spec.External.URL != "" { //nolint
		instance.Status.ExternalURL = instance.Spec.External.URL
	} else if currentState.KeycloakRoute != nil && currentState.KeycloakRoute.Spec.Host != "" {
		instance.Status.ExternalURL = fmt.Sprintf("https://%v", currentState.KeycloakRoute.Spec.Host)
	} else if currentState.KeycloakIngress != nil && currentState.KeycloakIngress.Spec.Rules[0].Host != "" {
		instance.Status.ExternalURL = fmt.Sprintf("https://%v", currentState.KeycloakIngress.Spec.Rules[0].Host)
	}

	// Let the clients know where the admin credentials are stored
	if currentState.KeycloakAdminSecret != nil {
		instance.Status.CredentialSecret = currentState.KeycloakAdminSecret.Name
	}

	r.setVersion(instance)

	err = r.client.Status().Update(r.context, instance)
	if err != nil {
		log.Error(err, "unable to update status")
		return reconcile.Result{
			RequeueAfter: RequeueDelayError,
			Requeue:      true,
		}, nil
	}

	log.Info("desired cluster state met")
	return reconcile.Result{RequeueAfter: RequeueDelay}, nil
}

func (r *ReconcileKeycloak) setVersion(instance *v1alpha1.Keycloak) {
	instance.Status.Version = version.Version
}
