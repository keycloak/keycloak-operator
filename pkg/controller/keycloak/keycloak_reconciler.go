package keycloak

import (
	"github.com/go-logr/logr"
	kc "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/common"
	"github.com/keycloak/keycloak-operator/pkg/model/keycloak"
)

type KeycloakReconciler struct {
	logger       logr.Logger
	clusterState *common.ClusterState
	runner       common.ActionRunner
}

func NewKeycloakReconciler(state *common.ClusterState, runner common.ActionRunner) *KeycloakReconciler {
	return &KeycloakReconciler{
		clusterState: state,
		runner:       runner,
	}
}

func (i *KeycloakReconciler) Reconcile(cr *kc.Keycloak) error {
	// Create the desired cluster state as a list of modifications to the
	// current state
	desiredState := i.buildDesiredState(cr)

	// Run all the modifications (actions)
	return i.runner.RunAll(desiredState)
}

func (i *KeycloakReconciler) buildDesiredState(cr *kc.Keycloak) common.DesiredClusterState {
	desired := common.DesiredClusterState{}
	desired = append(desired, i.getKeycloakServiceDesiredState(cr))
	return desired
}

func (i *KeycloakReconciler) getKeycloakServiceDesiredState(cr *kc.Keycloak) common.ClusterAction {
	service := keycloak.KeycloakService(cr)

	if i.clusterState.KeycloakService == nil {
		return common.GenericCreateAction{
			Ref: service,
			Msg: "create keycloak service",
		}
	} else {
		return common.ServiceUpdateAction{
			Ref:             service,
			Msg:             "update keycloak service",
			ClusterIP:       i.clusterState.KeycloakService.Spec.ClusterIP,
			ResourceVersion: i.clusterState.KeycloakService.ResourceVersion,
		}
	}
}
