package keycloakclient

import (
	"context"
	"fmt"
	"time"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	kc "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/common"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_keycloakclient")

const (
	ClientFinalizer   = "client.cleanup"
	RequeueDelayError = 5 * time.Second
	ControllerName    = "keycloakclient-controller"
)

// Add creates a new KeycloakClient Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	return &ReconcileKeycloakClient{
		client:   mgr.GetClient(),
		scheme:   mgr.GetScheme(),
		cancel:   cancel,
		context:  ctx,
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

	// Watch for changes to primary resource KeycloakClient
	err = c.Watch(&source.Kind{Type: &kc.KeycloakClient{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Make sure to watch the credential secrets
	err = c.Watch(&source.Kind{Type: &corev1.Secret{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &kc.KeycloakClient{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileKeycloakClient implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileKeycloakClient{}

// ReconcileKeycloakClient reconciles a KeycloakClient object
type ReconcileKeycloakClient struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client   client.Client
	scheme   *runtime.Scheme
	context  context.Context
	cancel   context.CancelFunc
	recorder record.EventRecorder
}

// Reconcile reads that state of the cluster for a KeycloakClient object and makes changes based on the state read
// and what is in the KeycloakClient.Spec
func (r *ReconcileKeycloakClient) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling KeycloakClient")

	// Fetch the KeycloakClient instance
	instance := &kc.KeycloakClient{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
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

	r.adjustCrDefaults(instance)

	// The client may be applicable to multiple keycloak instances,
	// process all of them
	realms, err := common.GetMatchingRealms(r.context, r.client, instance.Spec.RealmSelector)
	if err != nil {
		return r.ManageError(instance, err)
	}
	log.Info(fmt.Sprintf("found %v matching realm(s) for client %v/%v", len(realms.Items), instance.Namespace, instance.Name))
	for _, realm := range realms.Items {
		keycloaks, err := common.GetMatchingKeycloaks(r.context, r.client, realm.Spec.InstanceSelector)
		if err != nil {
			return r.ManageError(instance, err)
		}
		log.Info(fmt.Sprintf("found %v matching keycloak(s) for realm %v/%v", len(keycloaks.Items), realm.Namespace, realm.Name))

		for _, keycloak := range keycloaks.Items {
			// Get an authenticated keycloak api client for the instance
			keycloakFactory := common.LocalConfigKeycloakFactory{}
			authenticated, err := keycloakFactory.AuthenticatedClient(keycloak)
			if err != nil {
				return r.ManageError(instance, err)
			}

			// Compute the current state of the realm
			log.Info(fmt.Sprintf("got authenticated client for keycloak at %v", authenticated.Endpoint()))
			clientState := common.NewClientState(r.context, realm.DeepCopy())

			log.Info(fmt.Sprintf("read client state for keycloak %v/%v, realm %v/%v, client %v/%v",
				keycloak.Namespace,
				keycloak.Name,
				realm.Namespace,
				realm.Name,
				instance.Namespace,
				instance.Name))

			err = clientState.Read(r.context, instance, authenticated, r.client)
			if err != nil {
				return r.ManageError(instance, err)
			}

			// Figure out the actions to keep the realms up to date with
			// the desired state
			reconciler := NewKeycloakClientReconciler(keycloak)
			desiredState := reconciler.Reconcile(clientState, instance)
			actionRunner := common.NewClusterAndKeycloakActionRunner(r.context, r.client, r.scheme, instance, authenticated)

			// Run all actions to keep the realms updated
			err = actionRunner.RunAll(desiredState)
			if err != nil {
				return r.ManageError(instance, err)
			}
		}
	}

	return reconcile.Result{Requeue: false}, r.manageSuccess(instance, instance.DeletionTimestamp != nil)
}

// Fills the CR with default values. Nils are not acceptable for Kubernetes.
func (r *ReconcileKeycloakClient) adjustCrDefaults(cr *kc.KeycloakClient) {
	if cr.Spec.Client.Attributes == nil {
		cr.Spec.Client.Attributes = make(map[string]string)
	}
	if cr.Spec.Client.Access == nil {
		cr.Spec.Client.Access = make(map[string]bool)
	}
}

func (r *ReconcileKeycloakClient) manageSuccess(client *kc.KeycloakClient, deleted bool) error {
	client.Status.Ready = true
	client.Status.Message = ""
	client.Status.Phase = v1alpha1.PhaseReconciling

	err := r.client.Status().Update(r.context, client)
	if err != nil {
		log.Error(err, "unable to update status")
	}

	// Finalizer already set?
	finalizerExists := false
	for _, finalizer := range client.Finalizers {
		if finalizer == ClientFinalizer {
			finalizerExists = true
			break
		}
	}

	// Resource created and finalizer exists: nothing to do
	if !deleted && finalizerExists {
		return nil
	}

	// Resource created and finalizer does not exist: add finalizer
	if !deleted && !finalizerExists {
		client.Finalizers = append(client.Finalizers, ClientFinalizer)
		log.Info(fmt.Sprintf("added finalizer to keycloak client %v/%v",
			client.Namespace,
			client.Spec.Client.ClientID))

		return r.client.Update(r.context, client)
	}

	// Otherwise remove the finalizer
	newFinalizers := []string{}
	for _, finalizer := range client.Finalizers {
		if finalizer == ClientFinalizer {
			log.Info(fmt.Sprintf("removed finalizer from keycloak client %v/%v",
				client.Namespace,
				client.Spec.Client.ClientID))

			continue
		}
		newFinalizers = append(newFinalizers, finalizer)
	}

	client.Finalizers = newFinalizers
	return r.client.Update(r.context, client)
}

func (r *ReconcileKeycloakClient) ManageError(realm *kc.KeycloakClient, issue error) (reconcile.Result, error) {
	r.recorder.Event(realm, "Warning", "ProcessingError", issue.Error())

	realm.Status.Message = issue.Error()
	realm.Status.Ready = false
	realm.Status.Phase = v1alpha1.PhaseFailing

	err := r.client.Status().Update(r.context, realm)
	if err != nil {
		log.Error(err, "unable to update status")
	}

	return reconcile.Result{
		RequeueAfter: RequeueDelayError,
		Requeue:      true,
	}, nil
}
