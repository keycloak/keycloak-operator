package keycloakrealm

import (
	"testing"

	v12 "k8s.io/api/core/v1"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/common"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getDummyRealm() *v1alpha1.KeycloakRealm {
	return &v1alpha1.KeycloakRealm{
		Spec: v1alpha1.KeycloakRealmSpec{
			InstanceSelector: &v1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "keycloak",
				},
			},
			RealmOverrides: []*v1alpha1.RedirectorIdentityProviderOverride{
				{
					IdentityProvider: "openshift-v4",
					ForFlow:          "browser",
				},
			},
			Realm: &v1alpha1.KeycloakAPIRealm{
				ID:                        "dummy",
				Realm:                     "dummy",
				Enabled:                   true,
				DisplayName:               "dummy",
				EventsEnabled:             &[]bool{true}[0],
				AdminEventsEnabled:        &[]bool{true}[0],
				AdminEventsDetailsEnabled: &[]bool{true}[0],
				Users: []*v1alpha1.KeycloakAPIUser{
					{
						ID:        "dummy",
						UserName:  "dummy",
						FirstName: "dummy",
						LastName:  "dummy",
						Enabled:   true,
						Credentials: []v1alpha1.KeycloakCredential{
							{
								Type:      "password",
								Value:     "password",
								Temporary: false,
							},
						},
					},
				},
			},
		},
	}
}

func getDummyState() *common.RealmState {
	return &common.RealmState{
		Realm:            nil,
		RealmUserSecrets: nil,
		Context:          nil,
		Keycloak:         nil,
	}
}

func TestKeycloakRealmReconciler_Reconcile(t *testing.T) {
	// given
	keycloak := v1alpha1.Keycloak{}
	reconciler := NewKeycloakRealmReconciler(keycloak)

	realm := getDummyRealm()
	state := getDummyState()

	// when
	desiredState := reconciler.Reconcile(state, realm)

	// then
	// 0 - check keycloak available
	// 1 - create realm
	// 2 - create user credential secret
	// 3 - configure browser redirector
	assert.IsType(t, &common.PingAction{}, desiredState[0])
	assert.IsType(t, &common.CreateRealmAction{}, desiredState[1])
	assert.IsType(t, &common.GenericCreateAction{}, desiredState[2])
	assert.IsType(t, &common.ConfigureRealmAction{}, desiredState[3])

	assert.True(t, *realm.Spec.Realm.EventsEnabled)
	assert.True(t, *realm.Spec.Realm.AdminEventsEnabled)
	assert.True(t, *realm.Spec.Realm.AdminEventsDetailsEnabled)

	state.Realm = realm

	// Second round: realm is already created
	desiredState = reconciler.Reconcile(state, realm)
	assert.IsType(t, &common.PingAction{}, desiredState[0])

	// The user credential secret still needs to be created because we
	// did not set it in the state
	assert.IsType(t, &common.GenericCreateAction{}, desiredState[1])
}

func TestKeycloakRealmReconciler_ReconcileRealmDelete(t *testing.T) {
	// given
	keycloak := v1alpha1.Keycloak{}
	reconciler := NewKeycloakRealmReconciler(keycloak)

	realm := getDummyRealm()
	state := getDummyState()
	realm.DeletionTimestamp = &v1.Time{}

	// when
	desiredState := reconciler.Reconcile(state, realm)

	// then
	// 0 - check keycloak available
	// 1 - delete realm
	assert.IsType(t, &common.PingAction{}, desiredState[0])
	assert.IsType(t, &common.DeleteRealmAction{}, desiredState[1])
}

func TestKeycloakRealmReconciler_ReconcileCredentials(t *testing.T) {
	// given
	keycloak := v1alpha1.Keycloak{}
	reconciler := NewKeycloakRealmReconciler(keycloak)

	realm := getDummyRealm()
	state := getDummyState()

	// reset user credentials to force the operator to create a password
	realm.Spec.Realm.Users[0].Credentials = []v1alpha1.KeycloakCredential{}

	// when
	desiredState := reconciler.Reconcile(state, realm)

	// then
	// 0 - check keycloak available
	// 1 - create realm
	// 2 - create user credential secret
	// 3 - ensure a password is assigned automatically
	assert.IsType(t, &common.PingAction{}, desiredState[0])
	assert.IsType(t, &common.CreateRealmAction{}, desiredState[1])
	assert.IsType(t, &common.GenericCreateAction{}, desiredState[2])
}

func TestKeycloakRealmReconciler_Update(t *testing.T) {
	// given
	keycloak := v1alpha1.Keycloak{}
	reconciler := NewKeycloakRealmReconciler(keycloak)

	realm := getDummyRealm()
	state := getDummyState()

	// reset user credentials to force the operator to create a password
	state.Realm = realm
	state.RealmUserSecrets = make(map[string]*v12.Secret)
	state.RealmUserSecrets[realm.Spec.Realm.Users[0].UserName] = &v12.Secret{}

	// when
	desiredState := reconciler.Reconcile(state, realm)

	// then
	// 0 - check keycloak available
	// 1 - no other action added
	assert.IsType(t, &common.PingAction{}, desiredState[0])
	assert.Len(t, desiredState, 1)
}
