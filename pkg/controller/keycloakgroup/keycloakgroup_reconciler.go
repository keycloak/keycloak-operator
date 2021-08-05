package keycloakgroup

import (
	"fmt"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/common"
)

type Reconciler interface {
	Reconcile(cr *v1alpha1.KeycloakGroup) error
}

type KeycloakgroupReconciler struct { // nolint
	Realm    v1alpha1.KeycloakRealm
	Keycloak v1alpha1.Keycloak
}

func NewKeycloakgroupReconciler(keycloak v1alpha1.Keycloak, realm v1alpha1.KeycloakRealm) *KeycloakgroupReconciler {
	return &KeycloakgroupReconciler{
		Realm:    realm,
		Keycloak: keycloak,
	}
}

func (i *KeycloakgroupReconciler) Reconcile(state *common.GroupState, cr *v1alpha1.KeycloakGroup) common.DesiredClusterState {
	if cr.DeletionTimestamp != nil {
		return i.reconcileGroupDelete(state, cr)
	}
	return i.reconcileGroup(state, cr)
}

func (i *KeycloakgroupReconciler) reconcileGroup(state *common.GroupState, cr *v1alpha1.KeycloakGroup) common.DesiredClusterState {
	desired := common.DesiredClusterState{}

	desired.AddAction(i.getKeycloakDesiredState())
	desired.AddActions(i.getKeycloakGroupDesiredState(state, cr))

	return desired
}

func (i *KeycloakgroupReconciler) reconcileGroupDelete(state *common.GroupState, cr *v1alpha1.KeycloakGroup) common.DesiredClusterState {
	desired := common.DesiredClusterState{}

	desired.AddAction(i.getKeycloakDesiredState())

	// If there is an error retrieving the userID then the user has
	// probably been deleted in the Admin UI. Nothing to do for us
	// here
	if state.Group != nil {
		desired.AddAction(&common.DeleteGroupAction{
			ID:    state.Group.ID,
			Realm: i.Realm.Spec.Realm.Realm,
			Msg:   fmt.Sprintf("delete group %v", cr.Spec.Group.Name),
		})
	}

	return desired
}

// Always make sure keycloak is able to respond
func (i *KeycloakgroupReconciler) getKeycloakDesiredState() common.ClusterAction {
	return &common.PingAction{
		Msg: "check if keycloak is available",
	}
}

// Always make sure keycloak is able to respond
func (i *KeycloakgroupReconciler) getKeycloakGroupDesiredState(state *common.GroupState, cr *v1alpha1.KeycloakGroup) []common.ClusterAction {
	var actions []common.ClusterAction

	if state.Group == nil {
		actions = append(actions, &common.CreateGroupAction{
			Ref:   cr,
			Realm: i.Realm.Spec.Realm.Realm,
			Msg:   fmt.Sprintf("create group %v", cr.Spec.Group.Name),
		})
	} else {
		actions = append(actions, &common.UpdateGroupAction{
			Ref:   cr,
			Realm: i.Realm.Spec.Realm.Realm,
			Msg:   fmt.Sprintf("update group %v", cr.Spec.Group.Name),
		})

		//// Sync the requested roles
		actions = append(actions, i.getGroupRealmRolesDesiredState(state, cr)...)
		actions = append(actions, i.getGroupClientRolesDesiredState(state, cr)...)
	}

	return actions
}

func (i *KeycloakgroupReconciler) getGroupRealmRolesDesiredState(state *common.GroupState, cr *v1alpha1.KeycloakGroup) []common.ClusterAction {
	var assignRoles []common.ClusterAction
	var removeRoles []common.ClusterAction

	for _, role := range cr.Spec.Group.RealmRoles {
		// Is the role available for this user?
		roleRef := state.GetAvailableRealmRole(role)
		if roleRef == nil {
			continue
		}

		// Role requested but not assigned?
		if !containsRole(state.RealmRoles, role) {
			assignRoles = append(assignRoles, &common.AssignGroupRealmRoleAction{
				GroupID: state.Group.ID,
				Ref:     roleRef,
				Realm:   i.Realm.Spec.Realm.Realm,
				Msg:     fmt.Sprintf("assign realm role %v to group %v", role, state.Group.Name),
			})
		}
	}

	for _, role := range state.RealmRoles {
		// Role assigned but not requested?
		if !containsRoleID(cr.Spec.Group.RealmRoles, role.Name) {
			removeRoles = append(removeRoles, &common.RemoveGroupRealmRoleAction{
				GroupID: state.Group.ID,
				Ref:     role,
				Realm:   i.Realm.Spec.Realm.Realm,
				Msg:     fmt.Sprintf("remove realm role %v from group %v", role.Name, state.Group.Name),
			})
		}
	}

	return append(assignRoles, removeRoles...)
}

//
func (i *KeycloakgroupReconciler) getGroupClientRolesDesiredState(state *common.GroupState, cr *v1alpha1.KeycloakGroup) []common.ClusterAction {
	actions := []common.ClusterAction{}

	for _, client := range state.Clients {
		actions = append(actions, i.syncRolesForClient(state, cr, client.ClientID)...)
	}

	return actions
}

func (i *KeycloakgroupReconciler) syncRolesForClient(state *common.GroupState, cr *v1alpha1.KeycloakGroup, clientID string) []common.ClusterAction {
	var assignRoles []common.ClusterAction
	var removeRoles []common.ClusterAction

	for _, role := range cr.Spec.Group.ClientRoles[clientID] {
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
			assignRoles = append(assignRoles, &common.AssignGroupClientRoleAction{
				GroupID:  state.Group.ID,
				ClientID: client.ID,
				Ref:      roleRef,
				Realm:    i.Realm.Spec.Realm.Realm,
				Msg:      fmt.Sprintf("assign role %v of client %v to group %v", role, clientID, state.Group.Name),
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
		if !containsRoleID(cr.Spec.Group.ClientRoles[clientID], role.Name) {
			removeRoles = append(removeRoles, &common.RemoveGroupClientRoleAction{
				GroupID:  state.Group.ID,
				ClientID: client.ID,
				Ref:      role,
				Realm:    i.Realm.Spec.Realm.Realm,
				Msg:      fmt.Sprintf("remove role %v of client %v from group %v", role.Name, clientID, state.Group.Name),
			})
		}
	}

	return append(assignRoles, removeRoles...)
}

func containsRole(list []*v1alpha1.KeycloakUserRole, id string) bool {
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
