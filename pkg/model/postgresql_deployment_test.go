package model

import (
	"testing"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestUtil_Test_GetPostgresqlImage_Without_Override_Image_Set(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{}

	// when
	returnedImage := GetPostgresqlImage(cr)

	// then
	assert.Equal(t, returnedImage, PostgresqlImage)
}

func TestUtil_Test_GetPostgresqlImage_With_Wrong_Override_Image_Key(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			ImageOverrides: v1alpha1.KeycloakRelatedImages{
				Keycloak: "postgres:1.0",
			},
		},
	}

	// when
	returnedImage := GetPostgresqlImage(cr)

	// then
	assert.Equal(t, returnedImage, PostgresqlImage)
}

func TestUtil_Test_GetPostgresqlImage_With_Override_Image_Set(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			ImageOverrides: v1alpha1.KeycloakRelatedImages{
				Postgresql: "postgres:1.0",
			},
		},
	}

	// when
	returnedImage := GetPostgresqlImage(cr)

	// then
	assert.Equal(t, returnedImage, "postgres:1.0")
}

func Test_PostgresqlDeployment_With_Overrided_Image(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			ImageOverrides: v1alpha1.KeycloakRelatedImages{
				Postgresql: "postgres:1.0",
			},
		},
	}

	// when
	currentImage := PostgresqlDeployment(cr).Spec.Template.Spec.Containers[0].Image

	// then
	assert.Equal(t, currentImage, "postgres:1.0")
}

func Test_PostgresqlDeploymentReconciled_With_Overrided_Image(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{}
	cr2 := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			ImageOverrides: v1alpha1.KeycloakRelatedImages{
				Postgresql: "postgres:1.0",
			},
		},
	}

	// when
	currentState := PostgresqlDeployment(cr)
	currentImage := currentState.Spec.Template.Spec.Containers[0].Image

	reconciledState := PostgresqlDeploymentReconciled(cr2, currentState)
	reconciledImage := reconciledState.Spec.Template.Spec.Containers[0].Image

	// then
	assert.Equal(t, currentImage, PostgresqlImage)
	assert.Equal(t, reconciledImage, "postgres:1.0")
}
