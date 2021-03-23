package keycloakuser

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/keycloak/keycloak-operator/pkg/common"

	"k8s.io/client-go/tools/record"

	kc "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	ControllerName    = "controller_keycloakuser"
	RequeueDelayError = 5 * time.Second
)

var log = logf.Log.WithName("controller_keycloakuser")

// Add creates a new KeycloakUser Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	return &ReconcileKeycloakUser{
		client:   mgr.GetClient(),
		scheme:   mgr.GetScheme(),
		context:  ctx,
		cancel:   cancel,
		recorder: mgr.GetEventRecorderFor(ControllerName),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("keycloakuser-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource KeycloakUser
	err = c.Watch(&source.Kind{Type: &kc.KeycloakUser{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Pods and requeue the owner KeycloakUser
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &kc.KeycloakUser{},
	})
	if err != nil {
		return err
	}

	// Make sure to watch the credential secrets
	err = c.Watch(&source.Kind{Type: &corev1.Secret{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &kc.KeycloakUser{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileKeycloakUser implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileKeycloakUser{}

// ReconcileKeycloakUser reconciles a KeycloakUser object
type ReconcileKeycloakUser struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client   client.Client
	scheme   *runtime.Scheme
	context  context.Context
	cancel   context.CancelFunc
	recorder record.EventRecorder
}

// Reconcile reads that state of the cluster for a KeycloakUser object and makes changes based on the state read
// and what is in the KeycloakUser.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileKeycloakUser) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling KeycloakUser")

	// Fetch the KeycloakUser instance
	instance := &kc.KeycloakUser{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
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

	// If no selector is set we can't figure out which realm instance this user should
	// be added to. Skip reconcile until a selector has been set.
	if instance.Spec.RealmSelector == nil {
		log.Info(fmt.Sprintf("user %v/%v has no realm selector and will be ignored", instance.Namespace, instance.Name))
		return reconcile.Result{Requeue: false}, nil
	}

	// Find the realms that this user should be added to based on the label selector
	realms, err := common.GetMatchingRealms(r.context, r.client, instance.Spec.RealmSelector)
	if err != nil {
		return reconcile.Result{}, err
	}

	log.Info(fmt.Sprintf("found %v matching realm(s) for user %v/%v", len(realms.Items), instance.Namespace, instance.Name))

	for _, realm := range realms.Items {
		if realm.Spec.Unmanaged {
			return r.ManageError(instance, errors.Errorf("users cannot be created for unmanaged keycloak realms"))
		}

		keycloaks, err := common.GetMatchingKeycloaks(r.context, r.client, realm.Spec.InstanceSelector)
		if err != nil {
			return r.ManageError(instance, err)
		}

		for _, keycloak := range keycloaks.Items {
			if keycloak.Spec.Unmanaged {
				return r.ManageError(instance, errors.Errorf("users cannot be created for unmanaged keycloak instances"))
			}

			// Get an authenticated keycloak api client for the instance
			keycloakFactory := common.LocalConfigKeycloakFactory{}
			authenticated, err := keycloakFactory.AuthenticatedClient(keycloak)
			if err != nil {
				return r.ManageError(instance, err)
			}

			// Compute the current state of the realm
			userState := common.NewUserState(keycloak)

			log.Info(fmt.Sprintf("read state for keycloak %v/%v, realm %v/%v",
				keycloak.Namespace,
				keycloak.Name,
				instance.Namespace,
				realm.Spec.Realm.Realm))

			err = userState.Read(authenticated, r.client, instance, realm)
			if err != nil {
				return r.ManageError(instance, err)
			}
			reconciler := NewKeycloakuserReconciler(keycloak, realm)
			desiredState := reconciler.Reconcile(userState, instance)

			actionRunner := common.NewClusterAndKeycloakActionRunner(r.context, r.client, r.scheme, instance, authenticated)
			err = actionRunner.RunAll(desiredState)
			if err != nil {
				return r.ManageError(instance, err)
			}
		}
	}

	return reconcile.Result{Requeue: false}, r.manageSuccess(instance, instance.DeletionTimestamp != nil)
}

func (r *ReconcileKeycloakUser) manageSuccess(user *kc.KeycloakUser, deleted bool) error {
	user.Status.Phase = kc.UserPhaseReconciled

	err := r.client.Status().Update(r.context, user)
	if err != nil {
		log.Error(err, "unable to update status")
	}

	// Finalizer already set?
	finalizerExists := false
	for _, finalizer := range user.Finalizers {
		if finalizer == kc.UserFinalizer {
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
		user.Finalizers = append(user.Finalizers, kc.UserFinalizer)
		log.Info(fmt.Sprintf("added finalizer to keycloak user %v/%v", user.Namespace, user.Name))
		return r.client.Update(r.context, user)
	}

	// Otherwise remove the finalizer
	newFinalizers := []string{}
	for _, finalizer := range user.Finalizers {
		if finalizer == kc.UserFinalizer {
			log.Info(fmt.Sprintf("removed finalizer from keycloak user %v/%v", user.Namespace, user.Name))
			continue
		}
		newFinalizers = append(newFinalizers, finalizer)
	}

	user.Finalizers = newFinalizers
	return r.client.Update(r.context, user)
}

func (r *ReconcileKeycloakUser) ManageError(user *kc.KeycloakUser, issue error) (reconcile.Result, error) {
	r.recorder.Event(user, "Warning", "ProcessingError", issue.Error())

	user.Status.Phase = kc.UserPhaseFailing
	user.Status.Message = issue.Error()

	err := r.client.Status().Update(r.context, user)
	if err != nil {
		log.Error(err, "unable to update status")
	}

	return reconcile.Result{
		RequeueAfter: RequeueDelayError,
	}, nil
}
