package keycloak

import (
	kc "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/common"
	"github.com/keycloak/keycloak-operator/pkg/model/keycloak"
)

type Reconciler struct {
	clusterState *common.ClusterState
	runner       common.ActionRunner
}

func NewKeycloakReconciler(state *common.ClusterState, runner common.ActionRunner) *Reconciler {
	return &Reconciler{
		clusterState: state,
		runner:       runner,
	}
}

func (i *Reconciler) Reconcile(cr *kc.Keycloak) error {
	// Create the desired cluster state as a list of modifications to the
	// current state
	desiredState := i.buildDesiredState(cr)

	// Run all the modifications (actions)
	return i.runner.RunAll(desiredState)
}

func (i *Reconciler) buildDesiredState(cr *kc.Keycloak) common.DesiredClusterState {
	desired := common.DesiredClusterState{}
	desired = append(desired, i.getKeycloakServiceDesiredState(cr))
	return desired
}

func (i *Reconciler) getKeycloakServiceDesiredState(cr *kc.Keycloak) common.ClusterAction {
	service := keycloak.Service(cr)

	if i.clusterState.KeycloakService == nil {
		return common.GenericCreateAction{
			Ref: service,
			Msg: "create keycloak service",
		}
	}

	return common.ServiceUpdateAction{
		Ref:             service,
		Msg:             "update keycloak service",
		ClusterIP:       i.clusterState.KeycloakService.Spec.ClusterIP,
		ResourceVersion: i.clusterState.KeycloakService.ResourceVersion,
	}
}
