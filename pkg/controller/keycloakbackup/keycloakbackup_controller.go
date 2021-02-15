package keycloakbackup

import (
	"context"
	"fmt"
	"time"

	kc "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/common"
	errors "github.com/pkg/errors"
	v1 "k8s.io/api/batch/v1"
	"k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
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

const (
	RequeueDelay      = 30 * time.Second
	RequeueDelayError = 5 * time.Second
	ControllerName    = "keycloakbackup-controller"
)

var log = logf.Log.WithName("controller_keycloakbackup")

func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	return &ReconcileKeycloakBackup{
		client:   mgr.GetClient(),
		scheme:   mgr.GetScheme(),
		cancel:   cancel,
		context:  ctx,
		recorder: mgr.GetEventRecorderFor(ControllerName),
	}
}

func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("keycloakbackup-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource KeycloakBackup
	err = c.Watch(&source.Kind{Type: &kc.KeycloakBackup{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	if err := common.WatchSecondaryResource(c, ControllerName, common.JobKind, &v1.Job{}, &kc.KeycloakBackup{}); err != nil {
		return err
	}

	if err := common.WatchSecondaryResource(c, ControllerName, common.CronJobKind, &v1beta1.CronJob{}, &kc.KeycloakBackup{}); err != nil {
		return err
	}

	if err := common.WatchSecondaryResource(c, ControllerName, common.PersistentVolumeClaimKind, &corev1.PersistentVolumeClaim{}, &kc.KeycloakBackup{}); err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileKeycloakBackup implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileKeycloakBackup{}

// ReconcileKeycloakBackup reconciles a KeycloakBackup object
type ReconcileKeycloakBackup struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client   client.Client
	scheme   *runtime.Scheme
	context  context.Context
	cancel   context.CancelFunc
	recorder record.EventRecorder
}

// Reconcile reads that state of the cluster for a KeycloakBackup object and makes changes based on the state read
func (r *ReconcileKeycloakBackup) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling KeycloakBackup")

	// Fetch the KeycloakBackup instance
	instance := &kc.KeycloakBackup{}
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

	// If no selector is set we can't figure out which Keycloak instance this backup should
	// be added for. Skip reconcile until a selector has been set.
	if instance.Spec.InstanceSelector == nil {
		log.Info(fmt.Sprintf("backup %v/%v has no instance selector and will be ignored", instance.Namespace, instance.Name))
		return reconcile.Result{Requeue: false}, nil
	}

	keycloaks, err := common.GetMatchingKeycloaks(r.context, r.client, instance.Spec.InstanceSelector)
	if err != nil {
		return r.ManageError(instance, err)
	}

	// backups without instances to backup are treated as errors
	if len(keycloaks.Items) == 0 {
		return r.ManageError(instance, errors.Errorf("no instance to backup for %v/%v", instance.Namespace, instance.Name))
	}

	log.Info(fmt.Sprintf("found %v matching keycloak(s) for backup %v/%v", len(keycloaks.Items), instance.Namespace, instance.Name))

	var currentState *common.BackupState
	for _, keycloak := range keycloaks.Items {
		if keycloak.Spec.Unmanaged {
			return r.ManageError(instance, errors.Errorf("backups cannot be created for unmanaged keycloak instances"))
		}

		currentState = common.NewBackupState(keycloak)
		err = currentState.Read(r.context, instance, r.client)
		if err != nil {
			return r.ManageError(instance, err)
		}
		reconciler := NewKeycloakBackupReconciler(keycloak)
		desiredState := reconciler.Reconcile(currentState, instance)
		actionRunner := common.NewClusterActionRunner(r.context, r.client, r.scheme, instance)
		err = actionRunner.RunAll(desiredState)
		if err != nil {
			return r.ManageError(instance, err)
		}
	}

	return r.ManageSuccess(instance, currentState)
}

func (r *ReconcileKeycloakBackup) ManageError(instance *kc.KeycloakBackup, issue error) (reconcile.Result, error) {
	r.recorder.Event(instance, "Warning", "ProcessingError", issue.Error())

	instance.Status.Message = issue.Error()
	instance.Status.Ready = false
	instance.Status.Phase = kc.BackupPhaseFailing

	err := r.client.Status().Update(r.context, instance)
	if err != nil {
		log.Error(err, "unable to update status")
	}

	return reconcile.Result{
		RequeueAfter: RequeueDelayError,
		Requeue:      true,
	}, nil
}

func (r *ReconcileKeycloakBackup) ManageSuccess(instance *kc.KeycloakBackup, currentState *common.BackupState) (reconcile.Result, error) {
	resourcesReady, err := currentState.IsResourcesReady()
	if err != nil {
		return r.ManageError(instance, err)
	}
	instance.Status.Ready = resourcesReady
	instance.Status.Message = ""

	if resourcesReady {
		instance.Status.Phase = kc.BackupPhaseCreated
	} else {
		instance.Status.Phase = kc.BackupPhaseReconciling
	}

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
