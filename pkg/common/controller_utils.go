package common

import (
	"context"
	"fmt"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// These kinds are not provided by the openshift api
const (
	RouteKind                 = "Route"
	JobKind                   = "Job"
	CronJobKind               = "CronJob"
	SecretKind                = "Secret"
	StatefulSetKind           = "StatefulSet"
	ServiceKind               = "Service"
	IngressKind               = "Ingress"
	DeploymentKind            = "Deployment"
	PersistentVolumeClaimKind = "PersistentVolumeClaim"
	PodDisruptionBudgetKind   = "PodDisruptionBudget"
	OpenShiftAPIServerKind    = "OpenShiftAPIServer"
)

func WatchSecondaryResource(c controller.Controller, controllerName string, resourceKind string, objectTypetoWatch runtime.Object, cr runtime.Object) error {
	stateManager := GetStateManager()
	stateFieldName := GetStateFieldName(controllerName, resourceKind)

	// Avoid watching non-existing resources and watch duplication
	watchExists, _ := stateManager.GetState(stateFieldName).(bool)
	keyExists, _ := stateManager.GetState(resourceKind).(bool)
	if !keyExists || watchExists {
		return nil
	}

	// Set up the actual watch
	err := c.Watch(&source.Kind{Type: objectTypetoWatch}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    cr,
	})

	// Retry on error
	if err != nil {
		log.Error(err, "error creating watch")
		stateManager.SetState(stateFieldName, false)
		return err
	}

	stateManager.SetState(stateFieldName, true)
	log.Info(fmt.Sprintf("Watch created for '%s' resource in '%s'", resourceKind, controllerName))
	return nil
}

func GetStateFieldName(controllerName string, kind string) string {
	return controllerName + "-watch-" + kind
}

// Try to get a list of keycloak instances that match the selector specified on the realm
func GetMatchingKeycloaks(ctx context.Context, c client.Client, labelSelector *v1.LabelSelector) (v1alpha1.KeycloakList, error) {
	var list v1alpha1.KeycloakList
	opts := []client.ListOption{
		client.MatchingLabels(labelSelector.MatchLabels),
	}

	err := c.List(ctx, &list, opts...)
	return list, err
}

// Try to get a list of keycloak instances that match the selector specified on the realm
func GetMatchingRealms(ctx context.Context, c client.Client, labelSelector *v1.LabelSelector) (v1alpha1.KeycloakRealmList, error) {
	var list v1alpha1.KeycloakRealmList
	opts := []client.ListOption{
		client.MatchingLabels(labelSelector.MatchLabels),
	}

	err := c.List(ctx, &list, opts...)
	return list, err
}
