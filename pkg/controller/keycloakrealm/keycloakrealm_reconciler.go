package keycloakrealm

import (
	"fmt"

	kc "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/common"
	"github.com/keycloak/keycloak-operator/pkg/model"
)

type Reconciler interface {
	Reconcile(cr *kc.KeycloakRealm) error
}

type KeycloakRealmReconciler struct { // nolint
	Keycloak kc.Keycloak
}

func NewKeycloakRealmReconciler(keycloak kc.Keycloak) *KeycloakRealmReconciler {
	return &KeycloakRealmReconciler{
		Keycloak: keycloak,
	}
}

// Auto generate a password if the user didn't specify one
// It will be written to the secret
func ensureCredentials(users []*kc.KeycloakApiUser) {
	for _, user := range users {
		if len(user.Credentials) == 0 {
			user.Credentials = []kc.KeycloakCredential{
				{
					Type:      "password",
					Value:     model.RandStringRunes(10),
					Temporary: false,
				},
			}
		}
	}
}

func (i *KeycloakRealmReconciler) Reconcile(state *common.RealmState, cr *kc.KeycloakRealm) common.DesiredClusterState {
	if cr.DeletionTimestamp == nil {
		return i.ReconcileRealmCreate(state, cr)
	}
	return i.ReconcileRealmDelete(state, cr)
}

func (i *KeycloakRealmReconciler) ReconcileRealmCreate(state *common.RealmState, cr *kc.KeycloakRealm) common.DesiredClusterState {
	desired := common.DesiredClusterState{}

	desired.AddAction(i.getKeycloakDesiredState())
	desired.AddAction(i.getDesiredRealmState(state, cr))

	for _, user := range cr.Spec.Users {
		desired.AddAction(i.getDesiredUserSate(state, cr, user))
	}

	return desired
}

func (i *KeycloakRealmReconciler) ReconcileRealmDelete(state *common.RealmState, cr *kc.KeycloakRealm) common.DesiredClusterState {
	desired := common.DesiredClusterState{}
	desired.AddAction(i.getKeycloakDesiredState())
	desired.AddAction(i.getDesiredRealmState(state, cr))
	return desired
}

// Always make sure keycloak is able to respond
func (i *KeycloakRealmReconciler) getKeycloakDesiredState() common.ClusterAction {
	return &common.PingAction{
		Msg: "check if keycloak is available",
	}
}

func (i *KeycloakRealmReconciler) getDesiredRealmState(state *common.RealmState, cr *kc.KeycloakRealm) common.ClusterAction {
	if cr.DeletionTimestamp != nil {
		return &common.DeleteRealmAction{
			Ref: cr,
			Msg: fmt.Sprintf("removing realm %v/%v", cr.Namespace, cr.Spec.Realm),
		}
	}

	// Ensure that all users have credentials, if not provided then
	// automatically create a password. The user can later find the
	// credentials in the output secret
	ensureCredentials(cr.Spec.Users)

	if state.Realm == nil {
		return &common.CreateRealmAction{
			Ref: cr,
			Msg: fmt.Sprintf("create realm %v/%v", cr.Namespace, cr.Spec.Realm),
		}
	}

	return nil
}

func (i *KeycloakRealmReconciler) getDesiredUserSate(state *common.RealmState, cr *kc.KeycloakRealm, user *kc.KeycloakApiUser) common.ClusterAction {
	val, ok := state.RealmUserSecrets[user.UserName]
	if !ok || val == nil {
		return &common.GenericCreateAction{
			Ref: model.RealmCredentialSecret(cr, user, &i.Keycloak),
			Msg: fmt.Sprintf("create credential secret for user %v in realm %v/%v", user.UserName, cr.Namespace, cr.Spec.Realm),
		}
	}

	return nil
}
