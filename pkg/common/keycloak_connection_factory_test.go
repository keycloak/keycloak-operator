package common

import (
	"testing"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestDefaultKeycloakConnectionFactory_useClientCredentialsLoginOnAllFields(t *testing.T) {
	//given
	keycloakCR := &v1alpha1.Keycloak{}
	RealmCR := &v1alpha1.KeycloakRealm{}
	ClientCR := &v1alpha1.KeycloakClient{
		Status: v1alpha1.KeycloakClientStatus{
			Ready: true,
		},
	}
	factory := DefaultKeycloakConnectionFactory{
		keycloakCR: keycloakCR,
		realmCR:    RealmCR,
		clientCR:   ClientCR,
	}

	//when
	loginDecision, err := factory.loginDecision()

	//then
	assert.NoError(t, err)
	assert.Equal(t, UsingClientCredentials, loginDecision)
}

func TestDefaultKeycloakConnectionFactory_useAdminCredentialsLoginOnClientNotReady(t *testing.T) {
	//given
	keycloakCR := &v1alpha1.Keycloak{}
	RealmCR := &v1alpha1.KeycloakRealm{}
	ClientCR := &v1alpha1.KeycloakClient{}
	factory := DefaultKeycloakConnectionFactory{
		keycloakCR: keycloakCR,
		realmCR:    RealmCR,
		clientCR:   ClientCR,
	}

	//when
	loginDecision, err := factory.loginDecision()

	//then
	assert.NoError(t, err)
	assert.Equal(t, UsingAdminUsernameAndPassword, loginDecision)
}

func TestDefaultKeycloakConnectionFactory_useAdminLoginOnMissingClient(t *testing.T) {
	//given
	keycloakCR := &v1alpha1.Keycloak{}
	RealmCR := &v1alpha1.KeycloakRealm{}
	factory := DefaultKeycloakConnectionFactory{
		keycloakCR: keycloakCR,
		realmCR:    RealmCR,
	}

	//when
	loginDecision, err := factory.loginDecision()

	//then
	assert.NoError(t, err)
	assert.Equal(t, UsingAdminUsernameAndPassword, loginDecision)
}

func TestDefaultKeycloakConnectionFactory_throwAnErrorOnMissingRealm(t *testing.T) {
	//given
	RealmCR := &v1alpha1.KeycloakRealm{}
	ClientCR := &v1alpha1.KeycloakClient{}
	factory := DefaultKeycloakConnectionFactory{
		realmCR:  RealmCR,
		clientCR: ClientCR,
	}

	//when
	_, err := factory.loginDecision()

	//then
	assert.Error(t, err)
}
