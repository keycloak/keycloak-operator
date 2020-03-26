package keycloakclient

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/common"
	"github.com/keycloak/keycloak-operator/pkg/model"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	v13 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestKeycloakClientReconciler_Test_Creating_Client(t *testing.T) {
	// given
	keycloakCr := v1alpha1.Keycloak{}
	cr := &v1alpha1.KeycloakClient{
		ObjectMeta: v13.ObjectMeta{
			Name:      "test",
			Namespace: "test",
		},
		Spec: v1alpha1.KeycloakClientSpec{
			RealmSelector: &v13.LabelSelector{
				MatchLabels: map[string]string{"application": "sso"},
			},
			Client: &v1alpha1.KeycloakAPIClient{
				ClientID: "test",
				Secret:   "test",
			},
		},
	}

	currentState := &common.ClientState{
		Realm: &v1alpha1.KeycloakRealm{
			Spec: v1alpha1.KeycloakRealmSpec{
				Realm: &v1alpha1.KeycloakAPIRealm{
					Realm: "test",
				},
			},
		},
	}

	// when
	reconciler := NewKeycloakClientReconciler(keycloakCr)
	desiredState := reconciler.Reconcile(currentState, cr)

	// then
	assert.IsType(t, common.PingAction{}, desiredState[0])
	assert.IsType(t, common.CreateClientAction{}, desiredState[1])
	assert.IsType(t, common.GenericCreateAction{}, desiredState[2])
	assert.IsType(t, model.ClientSecret(cr), desiredState[2].(common.GenericCreateAction).Ref)
}

func TestKeycloakClientReconciler_Test_PartialUpdate_Client(t *testing.T) {
	// given
	keycloakCr := v1alpha1.Keycloak{}
	cr := &v1alpha1.KeycloakClient{
		ObjectMeta: v13.ObjectMeta{
			Name:      "test",
			Namespace: "test",
		},
		Spec: v1alpha1.KeycloakClientSpec{
			RealmSelector: &v13.LabelSelector{
				MatchLabels: map[string]string{"application": "sso"},
			},
			Client: &v1alpha1.KeycloakAPIClient{
				ClientID: "test",
				Secret:   "test",
			},
		},
	}

	currentState := &common.ClientState{
		Realm: &v1alpha1.KeycloakRealm{
			Spec: v1alpha1.KeycloakRealmSpec{
				Realm: &v1alpha1.KeycloakAPIRealm{
					Realm: "test",
				},
			},
		},
		Client: &v1alpha1.KeycloakAPIClient{
			Name: "dummy",
		},
	}

	// when
	reconciler := NewKeycloakClientReconciler(keycloakCr)
	desiredState := reconciler.Reconcile(currentState, cr)

	// then
	assert.IsType(t, common.PingAction{}, desiredState[0])
	assert.IsType(t, common.UpdateClientAction{}, desiredState[1])

	// client secret still needs to be created even if the client exists
	assert.IsType(t, common.GenericCreateAction{}, desiredState[2])
	assert.IsType(t, model.ClientSecret(cr), desiredState[2].(common.GenericCreateAction).Ref)
}

func TestKeycloakClientReconciler_Test_Delete_Client(t *testing.T) {
	// given
	keycloakCr := v1alpha1.Keycloak{}
	cr := &v1alpha1.KeycloakClient{
		ObjectMeta: v13.ObjectMeta{
			Name:      "test",
			Namespace: "test",
			DeletionTimestamp: &v13.Time{
				Time: time.Now(),
			},
		},
		Spec: v1alpha1.KeycloakClientSpec{
			RealmSelector: &v13.LabelSelector{
				MatchLabels: map[string]string{"application": "sso"},
			},
			Client: &v1alpha1.KeycloakAPIClient{
				ClientID: "test",
				Secret:   "test",
			},
		},
	}

	currentState := &common.ClientState{
		Realm: &v1alpha1.KeycloakRealm{
			Spec: v1alpha1.KeycloakRealmSpec{
				Realm: &v1alpha1.KeycloakAPIRealm{
					Realm: "test",
				},
			},
		},
	}

	// when
	reconciler := NewKeycloakClientReconciler(keycloakCr)
	desiredState := reconciler.Reconcile(currentState, cr)

	// then
	assert.IsType(t, common.PingAction{}, desiredState[0])
	assert.IsType(t, common.DeleteClientAction{}, desiredState[1])
}

func TestKeycloakClientReconciler_Test_Update_Client(t *testing.T) {
	// given
	keycloakCr := v1alpha1.Keycloak{}
	cr := &v1alpha1.KeycloakClient{
		ObjectMeta: v13.ObjectMeta{
			Name:      "test",
			Namespace: "test",
		},
		Spec: v1alpha1.KeycloakClientSpec{
			RealmSelector: &v13.LabelSelector{
				MatchLabels: map[string]string{"application": "sso"},
			},
			Client: &v1alpha1.KeycloakAPIClient{
				ClientID: "test",
				Secret:   "test",
			},
		},
	}

	currentState := &common.ClientState{
		Client:       &v1alpha1.KeycloakAPIClient{},
		ClientSecret: &v1.Secret{},
		Realm: &v1alpha1.KeycloakRealm{
			Spec: v1alpha1.KeycloakRealmSpec{
				Realm: &v1alpha1.KeycloakAPIRealm{
					Realm: "test",
				},
			},
		},
	}

	// when
	reconciler := NewKeycloakClientReconciler(keycloakCr)
	desiredState := reconciler.Reconcile(currentState, cr)

	// then
	assert.IsType(t, common.PingAction{}, desiredState[0])
	assert.IsType(t, common.UpdateClientAction{}, desiredState[1])
	assert.Equal(t, "test", desiredState[1].(common.UpdateClientAction).Realm)
	assert.IsType(t, common.GenericUpdateAction{}, desiredState[2])
	assert.IsType(t, model.ClientSecretReconciled(cr, currentState.ClientSecret), desiredState[2].(common.GenericUpdateAction).Ref)
}

func TestKeycloakClientReconciler_Test_Marshal_Client_directAccessGrantsEnabled(t *testing.T) {
	// given
	cr := &v1alpha1.KeycloakClient{
		ObjectMeta: v13.ObjectMeta{
			Name:      "test",
			Namespace: "test",
		},
		Spec: v1alpha1.KeycloakClientSpec{
			RealmSelector: &v13.LabelSelector{
				MatchLabels: map[string]string{"application": "sso"},
			},
			Client: &v1alpha1.KeycloakAPIClient{
				ClientID:                  "test",
				Secret:                    "test",
				DirectAccessGrantsEnabled: false,
			},
		},
	}

	// when
	b, err := json.Marshal(cr)
	s := string(b)

	// then
	assert.Nil(t, err, "Client couldn't be marshalled")
	assert.True(t, strings.Contains(s, "\"directAccessGrantsEnabled\":false"), "Element directAccessGrantsEnabled should not be omitted if false, as keycloaks default is true")
}
