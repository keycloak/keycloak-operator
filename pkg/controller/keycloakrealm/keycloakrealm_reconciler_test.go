package keycloakrealm

import (
	"testing"

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
			KeycloakApiRealm: &v1alpha1.KeycloakApiRealm{
				ID:          "dummy",
				Realm:       "dummy",
				Enabled:     true,
				DisplayName: "dummy",
				Users: []*v1alpha1.KeycloakApiUser{
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
	assert.IsType(t, &common.PingAction{}, desiredState[0])
	assert.IsType(t, &common.CreateRealmAction{}, desiredState[1])
	assert.IsType(t, &common.GenericCreateAction{}, desiredState[2])

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
	realm.Spec.Users[0].Credentials = []v1alpha1.KeycloakCredential{}

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
	assert.Len(t, realm.Spec.Users[0].Credentials, 1)
}
