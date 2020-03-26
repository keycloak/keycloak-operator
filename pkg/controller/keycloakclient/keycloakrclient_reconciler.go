package keycloakclient

import (
	"fmt"

	kc "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/common"
	"github.com/keycloak/keycloak-operator/pkg/model"
)

type Reconciler interface {
	Reconcile(cr *kc.KeycloakClient) error
}

type KeycloakClientReconciler struct { // nolint
	Keycloak kc.Keycloak
}

func NewKeycloakClientReconciler(keycloak kc.Keycloak) *KeycloakClientReconciler {
	return &KeycloakClientReconciler{
		Keycloak: keycloak,
	}
}

func (i *KeycloakClientReconciler) Reconcile(state *common.ClientState, cr *kc.KeycloakClient) common.DesiredClusterState {
	desired := common.DesiredClusterState{}

	desired.AddAction(i.pingKeycloak())
	if cr.DeletionTimestamp != nil { //nolint
		desired.AddAction(i.getDeletedClientState(state, cr))
		return desired
	}

	if state.Client == nil {
		desired.AddAction(i.getCreatedClientState(state, cr))
	} else {
		desired.AddAction(i.getUpdatedClientState(state, cr))
	}

	if state.ClientSecret == nil {
		desired.AddAction(i.getCreatedClientSecretState(state, cr))
	} else {
		desired.AddAction(i.getUpdatedClientSecretState(state, cr))
	}

	return desired
}

func (i *KeycloakClientReconciler) pingKeycloak() common.ClusterAction {
	return common.PingAction{
		Msg: "check if keycloak is available",
	}
}

func (i *KeycloakClientReconciler) getDeletedClientState(state *common.ClientState, cr *kc.KeycloakClient) common.ClusterAction {
	return common.DeleteClientAction{
		Ref:   cr,
		Realm: state.Realm.Spec.Realm.Realm,
		Msg:   fmt.Sprintf("removing client %v/%v", cr.Namespace, cr.Spec.Client.ClientID),
	}
}

func (i *KeycloakClientReconciler) getCreatedClientState(state *common.ClientState, cr *kc.KeycloakClient) common.ClusterAction {
	return common.CreateClientAction{
		Ref:   cr,
		Realm: state.Realm.Spec.Realm.Realm,
		Msg:   fmt.Sprintf("create client %v/%v", cr.Namespace, cr.Spec.Client.ID),
	}
}

func (i *KeycloakClientReconciler) getUpdatedClientSecretState(state *common.ClientState, cr *kc.KeycloakClient) common.ClusterAction {
	return common.GenericUpdateAction{
		Ref: model.ClientSecretReconciled(cr, state.ClientSecret),
		Msg: fmt.Sprintf("update client secret %v/%v", cr.Namespace, cr.Spec.Client.ID),
	}
}

func (i *KeycloakClientReconciler) getUpdatedClientState(state *common.ClientState, cr *kc.KeycloakClient) common.ClusterAction {
	return common.UpdateClientAction{
		Ref:   cr,
		Realm: state.Realm.Spec.Realm.Realm,
		Msg:   fmt.Sprintf("update client %v/%v", cr.Namespace, cr.Spec.Client.ID),
	}
}

func (i *KeycloakClientReconciler) getCreatedClientSecretState(state *common.ClientState, cr *kc.KeycloakClient) common.ClusterAction {
	return common.GenericCreateAction{
		Ref: model.ClientSecret(cr),
		Msg: fmt.Sprintf("create client secret %v/%v", cr.Namespace, cr.Spec.Client.ID),
	}
}
