package keycloakuser

import (
	"fmt"

	"github.com/keycloak/keycloak-operator/pkg/model"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/common"
)

type Reconciler interface {
	Reconcile(cr *v1alpha1.KeycloakUser) error
}

type KeycloakuserReconciler struct { // nolint
	Realm    v1alpha1.KeycloakRealm
	Keycloak v1alpha1.Keycloak
}

func NewKeycloakuserReconciler(keycloak v1alpha1.Keycloak, realm v1alpha1.KeycloakRealm) *KeycloakuserReconciler {
	return &KeycloakuserReconciler{
		Realm:    realm,
		Keycloak: keycloak,
	}
}

func (i *KeycloakuserReconciler) Reconcile(state *common.UserState, cr *v1alpha1.KeycloakUser) common.DesiredClusterState {
	if cr.DeletionTimestamp != nil {
		return i.reconcileUserDelete(state, cr)
	}
	return i.reconcileUser(state, cr)
}

func (i *KeycloakuserReconciler) reconcileUser(state *common.UserState, cr *v1alpha1.KeycloakUser) common.DesiredClusterState {
	desired := common.DesiredClusterState{}

	desired.AddAction(i.getKeycloakDesiredState())
	desired.AddActions(i.getKeycloakUserDesiredState(state, cr))
	desired.AddAction(i.getUserSecretDesiredState(state, cr))

	return desired
}

func (i *KeycloakuserReconciler) reconcileUserDelete(state *common.UserState, cr *v1alpha1.KeycloakUser) common.DesiredClusterState {
	desired := common.DesiredClusterState{}

	desired.AddAction(i.getKeycloakDesiredState())

	// If there is an error retrieving the userID then the user has
	// probably been deleted in the Admin UI. Nothing to do for us
	// here
	if state.User != nil {
		desired.AddAction(&common.DeleteUserAction{
			ID:    state.User.ID,
			Realm: i.Realm.Spec.Realm.Realm,
			Msg:   fmt.Sprintf("delete user %v", cr.Spec.User.UserName),
		})
	}

	return desired
}

// Always make sure keycloak is able to respond
func (i *KeycloakuserReconciler) getKeycloakDesiredState() common.ClusterAction {
	return &common.PingAction{
		Msg: "check if keycloak is available",
	}
}

// Always make sure keycloak is able to respond
func (i *KeycloakuserReconciler) getKeycloakUserDesiredState(state *common.UserState, cr *v1alpha1.KeycloakUser) []common.ClusterAction {
	var actions []common.ClusterAction

	if state.User == nil {
		actions = append(actions, &common.CreateUserAction{
			Ref:   cr,
			Realm: i.Realm.Spec.Realm.Realm,
			Msg:   fmt.Sprintf("create user %v", cr.Spec.User.UserName),
		})
	} else {
		actions = append(actions, &common.UpdateUserAction{
			Ref:   cr,
			Realm: i.Realm.Spec.Realm.Realm,
			Msg:   fmt.Sprintf("update user %v", cr.Spec.User.UserName),
		})

		// Sync the requested roles
		actions = append(actions, i.getUserRealmRolesDesiredState(state, cr)...)
		actions = append(actions, i.getUserClientRolesDesiredState(state, cr)...)
	}

	return actions
}

func (i *KeycloakuserReconciler) getUserRealmRolesDesiredState(state *common.UserState, cr *v1alpha1.KeycloakUser) []common.ClusterAction {
	var assignRoles []common.ClusterAction
	var removeRoles []common.ClusterAction

	for _, role := range cr.Spec.User.RealmRoles {
		// Is the role available for this user?
		roleRef := state.GetAvailableRealmRole(role)
		if roleRef == nil {
			continue
		}

		// Role requested but not assigned?
		if !containsRole(state.RealmRoles, role) {
			assignRoles = append(assignRoles, &common.AssignRealmRoleAction{
				UserID: state.User.ID,
				Ref:    roleRef,
				Realm:  i.Realm.Spec.Realm.Realm,
				Msg:    fmt.Sprintf("assign realm role %v to user %v", role, state.User.UserName),
			})
		}
	}

	for _, role := range state.RealmRoles {
		// Role assigned but not requested?
		if !containsRoleID(cr.Spec.User.RealmRoles, role.Name) {
			removeRoles = append(removeRoles, &common.RemoveRealmRoleAction{
				UserID: state.User.ID,
				Ref:    role,
				Realm:  i.Realm.Spec.Realm.Realm,
				Msg:    fmt.Sprintf("remove realm role %v from user %v", role.Name, state.User.UserName),
			})
		}
	}

	return append(assignRoles, removeRoles...)
}

func (i *KeycloakuserReconciler) getUserClientRolesDesiredState(state *common.UserState, cr *v1alpha1.KeycloakUser) []common.ClusterAction {
	actions := []common.ClusterAction{}

	for _, client := range state.Clients {
		actions = append(actions, i.syncRolesForClient(state, cr, client.ClientID)...)
	}

	return actions
}

func (i *KeycloakuserReconciler) getUserSecretDesiredState(state *common.UserState, cr *v1alpha1.KeycloakUser) common.ClusterAction {
	// Only ever create the secret, because we can't know when the
	// users change their credentials in keycloak. Also the owner
	// reference ensures that it gets deleted once the User CR is
	// deleted
	if state.Secret == nil {
		return &common.GenericCreateAction{
			Ref: model.RealmCredentialSecret(&i.Realm, &cr.Spec.User, &i.Keycloak),
			Msg: fmt.Sprintf("create credential secret for user %v in realm %v/%v",
				cr.Spec.User.UserName,
				cr.Namespace,
				i.Realm.Spec.Realm.Realm),
		}
	}
	return nil
}

func (i *KeycloakuserReconciler) syncRolesForClient(state *common.UserState, cr *v1alpha1.KeycloakUser, clientID string) []common.ClusterAction {
	var assignRoles []common.ClusterAction
	var removeRoles []common.ClusterAction

	for _, role := range cr.Spec.User.ClientRoles[clientID] {
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
				UserID:   state.User.ID,
				ClientID: client.ID,
				Ref:      roleRef,
				Realm:    i.Realm.Spec.Realm.Realm,
				Msg:      fmt.Sprintf("assign role %v of client %v to user %v", role, clientID, state.User.UserName),
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
		if !containsRoleID(cr.Spec.User.ClientRoles[clientID], role.Name) {
			removeRoles = append(removeRoles, &common.RemoveClientRoleAction{
				UserID:   state.User.ID,
				ClientID: client.ID,
				Ref:      role,
				Realm:    i.Realm.Spec.Realm.Realm,
				Msg:      fmt.Sprintf("remove role %v of client %v from user %v", role.Name, clientID, state.User.UserName),
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
