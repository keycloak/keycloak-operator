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
	if cr.DeletionTimestamp != nil {
		desired.AddAction(i.getDeletedClientState(state, cr))
		return desired
	}

	if state.Client == nil {
		desired.AddAction(i.getCreatedClientState(state, cr))
	} else {
		desired.AddActions(i.getUpdatedClientState(state, cr))
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
		Msg:   fmt.Sprintf("create client %v/%v", cr.Namespace, cr.Spec.Client.ClientID),
	}
}

func (i *KeycloakClientReconciler) getUpdatedClientSecretState(state *common.ClientState, cr *kc.KeycloakClient) common.ClusterAction {
	return common.GenericUpdateAction{
		Ref: model.ClientSecretReconciled(cr, state.ClientSecret),
		Msg: fmt.Sprintf("update client secret %v/%v", cr.Namespace, cr.Spec.Client.ClientID),
	}
}

func (i *KeycloakClientReconciler) getUpdatedClientState(state *common.ClientState, cr *kc.KeycloakClient) []common.ClusterAction {
	var actions []common.ClusterAction
	actions = append(actions, common.UpdateClientAction{
		Ref:   cr,
		Realm: state.Realm.Spec.Realm.Realm,
		Msg:   fmt.Sprintf("update client %v/%v", cr.Namespace, cr.Spec.Client.ClientID),
	})

	// Sync the requested roles
	actions = append(actions, i.getServiceAccountRealmRolesDesiredState(state, cr)...)
	actions = append(actions, i.getServiceAccountClientRolesDesiredState(state, cr)...)

	return actions
}

func (i *KeycloakClientReconciler) getCreatedClientSecretState(state *common.ClientState, cr *kc.KeycloakClient) common.ClusterAction {
	return common.GenericCreateAction{
		Ref: model.ClientSecret(cr),
		Msg: fmt.Sprintf("create client secret %v/%v", cr.Namespace, cr.Spec.Client.ClientID),
	}
}

func (i *KeycloakClientReconciler) getServiceAccountRealmRolesDesiredState(state *common.ClientState, cr *kc.KeycloakClient) []common.ClusterAction {
	var assignRoles []common.ClusterAction
	var removeRoles []common.ClusterAction

	for _, role := range cr.Spec.ServiceAccountRoles.RealmRoles {
		// Is the role available for this user?
		roleRef := state.GetAvailableRealmRole(role)
		if roleRef == nil {
			continue
		}

		// Role requested but not assigned?
		if !containsRole(state.RealmRoles, role) {
			assignRoles = append(assignRoles, &common.AssignRealmRoleAction{
				UserID: cr.Status.RealmServiceAccountIds[state.Realm.Spec.Realm.Realm],
				Ref:    roleRef,
				Realm:  state.Realm.Spec.Realm.Realm,
				Msg:    fmt.Sprintf("assign realm role %v to service account %v", role, state.Client.Name),
			})
		}
	}

	for _, role := range state.RealmRoles {
		// Role assigned but not requested?
		if !containsRoleID(cr.Spec.ServiceAccountRoles.RealmRoles, role.Name) {
			removeRoles = append(removeRoles, &common.RemoveRealmRoleAction{
				UserID: cr.Status.RealmServiceAccountIds[state.Realm.Spec.Realm.Realm],
				Ref:    role,
				Realm:  state.Realm.Spec.Realm.Realm,
				Msg:    fmt.Sprintf("remove realm role %v from service account %v", role.Name, state.Client.Name),
			})
		}
	}

	return append(assignRoles, removeRoles...)
}

func (i *KeycloakClientReconciler) getServiceAccountClientRolesDesiredState(state *common.ClientState, cr *kc.KeycloakClient) []common.ClusterAction {
	actions := []common.ClusterAction{}

	for _, client := range state.Clients {
		actions = append(actions, i.syncRolesForClient(state, cr, client.ClientID)...)
	}

	return actions
}

func (i *KeycloakClientReconciler) syncRolesForClient(state *common.ClientState, cr *kc.KeycloakClient, clientID string) []common.ClusterAction {
	var assignRoles []common.ClusterAction
	var removeRoles []common.ClusterAction

	for _, role := range cr.Spec.ServiceAccountRoles.ClientRoles[clientID] {
		// Is the role available for this user?
		roleRef := state.GetAvailableClientRole(role, clientID)
		if roleRef == nil {
			continue
		}

		// Valid client?
		client := state.GetClientByID(clientID)
		if client == nil {
			continue
		}

		// Role requested but not assigned?
		if !containsRole(state.ClientRoles[clientID], role) {
			assignRoles = append(assignRoles, &common.AssignClientRoleAction{
				UserID:   cr.Status.RealmServiceAccountIds[state.Realm.Spec.Realm.Realm],
				ClientID: client.ID,
				Ref:      roleRef,
				Realm:    state.Realm.Spec.Realm.Realm,
				Msg:      fmt.Sprintf("assign role %v of client %v to client %v", role, clientID, state.Client.Name),
			})
		}
	}

	for _, role := range state.ClientRoles[clientID] {
		// Valid client?
		client := state.GetClientByID(clientID)
		if client == nil {
			continue
		}

		// Role assigned but not requested?
		if !containsRoleID(cr.Spec.ServiceAccountRoles.ClientRoles[clientID], role.Name) {
			removeRoles = append(removeRoles, &common.RemoveClientRoleAction{
				UserID:   cr.Status.RealmServiceAccountIds[state.Realm.Spec.Realm.Realm],
				ClientID: client.ID,
				Ref:      role,
				Realm:    state.Realm.Spec.Realm.Realm,
				Msg:      fmt.Sprintf("remove role %v of client %v from client %v", role.Name, clientID, state.Client.Name),
			})
		}
	}

	return append(assignRoles, removeRoles...)
}

func containsRole(list []*kc.KeycloakUserRole, id string) bool {
	for _, item := range list {
		if item.ID == id {
			return true
		}
	}
	return false
}

func containsRoleID(list []string, id string) bool {
	for _, item := range list {
		if item == id {
			return true
		}
	}
	return false
}
