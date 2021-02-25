package model

import (
	"testing"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestKeycloakRoute_testHost(t *testing.T) {
	//given
	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			ExternalAccess: v1alpha1.KeycloakExternalAccess{
				Enabled: true,
			},
		},
	}

	//when
	route := KeycloakRoute(cr)

	//then
	assert.Equal(t, "", route.Spec.Host)
}

func TestKeycloakRoute_testHostNotTakenIntoAccount(t *testing.T) {
	//given
	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			ExternalAccess: v1alpha1.KeycloakExternalAccess{
				Enabled: true,
				Host:    "host-override",
			},
		},
	}

	//when
	route := KeycloakRoute(cr)

	//then
	assert.Equal(t, "", route.Spec.Host)
}
