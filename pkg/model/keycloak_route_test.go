package model

import (
	"testing"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	v1 "github.com/openshift/api/route/v1"
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

func TestKeycloakRoute_testHostReconciled(t *testing.T) {
	//given
	currentState := &v1.Route{
		Spec: v1.RouteSpec{},
	}

	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			ExternalAccess: v1alpha1.KeycloakExternalAccess{
				Enabled: true,
			},
		},
	}

	//when
	reconciledRoute := KeycloakRouteReconciled(cr, currentState)

	//then
	assert.Equal(t, "", reconciledRoute.Spec.Host)
}

func TestKeycloakRoute_testHostOverride(t *testing.T) {
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
	assert.Equal(t, "host-override", route.Spec.Host)
}

func TestKeycloakRoute_testHostOverrideReconciled(t *testing.T) {
	//given
	currentState := &v1.Route{
		Spec: v1.RouteSpec{},
	}

	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			ExternalAccess: v1alpha1.KeycloakExternalAccess{
				Enabled: true,
				Host:    "host-override",
			},
		},
	}

	//when
	reconciledRoute := KeycloakRouteReconciled(cr, currentState)

	//then
	assert.Equal(t, "host-override", reconciledRoute.Spec.Host)
}
