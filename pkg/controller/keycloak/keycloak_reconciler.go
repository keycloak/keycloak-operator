package keycloak

import (
	kc "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/common"
	"github.com/keycloak/keycloak-operator/pkg/model/keycloak"
)

type Reconciler interface {
	Reconcile(clusterState *common.ClusterState, cr *kc.Keycloak) (common.DesiredClusterState, error)
}

type KeycloakReconciler struct { // nolint
}

func NewKeycloakReconciler() *KeycloakReconciler {
	return &KeycloakReconciler{}
}

func (i *KeycloakReconciler) Reconcile(clusterState *common.ClusterState, cr *kc.Keycloak) (common.DesiredClusterState, error) {
	desired := common.DesiredClusterState{}
	desired = append(desired, i.getKeycloakServiceDesiredState(clusterState, cr))
	return desired, nil
}

func (i *KeycloakReconciler) getKeycloakServiceDesiredState(clusterState *common.ClusterState, cr *kc.Keycloak) common.ClusterAction {
	service := keycloak.Service(cr)

	if clusterState.KeycloakService == nil {
		return common.GenericCreateAction{
			Ref: service,
			Msg: "create keycloak service",
		}
	}

	// This part may change in the future once we have more resources to reconcile.
	// Perhaps there should be another method, like `keycloak.Service(cr, clusterState)`?
	service.Spec.ClusterIP = clusterState.KeycloakService.Spec.ClusterIP
	service.ResourceVersion = clusterState.KeycloakService.ResourceVersion
	return common.GenericUpdateAction{
		Ref: service,
		Msg: "update keycloak service",
	}
}
