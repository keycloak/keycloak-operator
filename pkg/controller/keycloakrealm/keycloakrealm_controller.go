package keycloakrealm

import (
	"context"
	"fmt"

	kc "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/common"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
	config2 "sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_keycloakrealm")

const (
	RealmFinalizer = "realm.cleanup"
)

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new KeycloakRealm Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager, _ chan schema.GroupVersionKind) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	return &ReconcileKeycloakRealm{
		client:  mgr.GetClient(),
		scheme:  mgr.GetScheme(),
		cancel:  cancel,
		context: ctx,
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("keycloakrealm-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource KeycloakRealm
	err = c.Watch(&source.Kind{Type: &kc.KeycloakRealm{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Make sure to watch the credential secrets
	err = c.Watch(&source.Kind{Type: &corev1.Secret{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &kc.KeycloakRealm{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileKeycloakRealm implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileKeycloakRealm{}

// ReconcileKeycloakRealm reconciles a KeycloakRealm object
type ReconcileKeycloakRealm struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client  client.Client
	scheme  *runtime.Scheme
	context context.Context
	cancel  context.CancelFunc
}

// Reconcile reads that state of the cluster for a KeycloakRealm object and makes changes based on the state read
// and what is in the KeycloakRealm.Spec
func (r *ReconcileKeycloakRealm) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling KeycloakRealm")

	// Fetch the KeycloakRealm instance
	instance := &kc.KeycloakRealm{}
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

	// If no selector is set we can't figure out which Keycloak instance this realm should
	// be added to. Skip reconcile until a selector has been set.
	if instance.Spec.InstanceSelector == nil {
		log.Info(fmt.Sprintf("realm %v/%v has no instance selector and will be ignored", instance.Namespace, instance.Name))
		return reconcile.Result{Requeue: false}, nil
	}

	keycloaks, err := r.getMatchingKeycloaks(instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	log.Info(fmt.Sprintf("found %v matching keycloak(s) for realm %v/%v", len(keycloaks.Items), instance.Namespace, instance.Name))

	// The realm may be applicable to multiple keycloak instances,
	// process all of them
	for _, keycloak := range keycloaks.Items {
		// Get an authenticated keycloak api client for the instance
		authenticated, err := r.getAuthenticatedClient(keycloak, instance)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Compute the current state of the realm
		log.Info(fmt.Sprintf("got authenticated client for keycloak at %v", keycloak.Status.InternalURL))
		realmState := common.NewRealmState(r.context, keycloak)

		log.Info(fmt.Sprintf("read state for keycloak %v/%v, realm %v/%v",
			keycloak.Namespace,
			keycloak.Name,
			instance.Namespace,
			instance.Spec.Realm))

		err = realmState.Read(instance, authenticated, r.client)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Figure out the actions to keep the realms up to date with
		// the desired state
		reconciler := NewKeycloakRealmReconciler(keycloak)
		desiredState := reconciler.Reconcile(realmState, instance)
		actionRunner := common.NewRealmActionRunner(r.context, r.client, r.scheme, instance, authenticated)

		// Run all actions to keep the realms updated
		err = actionRunner.RunAll(desiredState)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{Requeue: false}, r.manageSuccess(instance, instance.DeletionTimestamp != nil)
}

func (r *ReconcileKeycloakRealm) getAuthenticatedClient(kc kc.Keycloak, realm *kc.KeycloakRealm) (common.KeycloakInterface, error) {
	config, err := config2.GetConfig()
	if err != nil {
		return nil, err
	}

	secretClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	factory := common.KeycloakFactory{
		SecretClient: secretClient.CoreV1().Secrets(kc.Namespace),
	}

	return factory.AuthenticatedClient(kc)
}

// Try to get a list of keycloak instances that match the selector specified on the realm
func (r *ReconcileKeycloakRealm) getMatchingKeycloaks(realm *kc.KeycloakRealm) (kc.KeycloakList, error) {
	var list kc.KeycloakList
	opts := &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(realm.Spec.InstanceSelector.MatchLabels),
	}

	err := r.client.List(r.context, opts, &list)
	if err != nil {
		return list, err
	}

	return list, nil
}

func (r *ReconcileKeycloakRealm) manageSuccess(realm *kc.KeycloakRealm, deleted bool) error {
	// Finalizer already set?
	finalizerExists := false
	for _, finalizer := range realm.Finalizers {
		if finalizer == RealmFinalizer {
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
		realm.Finalizers = append(realm.Finalizers, RealmFinalizer)
		log.Info(fmt.Sprintf("added finalizer to keycloak realm %v/%v", realm.Namespace, realm.Spec.Realm))
		return r.client.Update(r.context, realm)
	}

	// Otherwise remove the finalizer
	newFinalizers := []string{}
	for _, finalizer := range realm.Finalizers {
		if finalizer == RealmFinalizer {
			log.Info(fmt.Sprintf("removed finalizer from keycloak realm %v/%v", realm.Namespace, realm.Spec.Realm))
			continue
		}
		newFinalizers = append(newFinalizers, finalizer)
	}

	realm.Finalizers = newFinalizers
	return r.client.Update(r.context, realm)
}
