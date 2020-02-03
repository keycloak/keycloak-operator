package model

import (
	"testing"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
)

func TestUtil_Test_GetKeycloakImage_Without_Override_Image_Set(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{}

	// when
	returnedImage := GetKeycloakImage(cr)

	// then
	assert.Equal(t, returnedImage, KeycloakImage)
}

func TestUtil_Test_GetKeycloakImage_With_Wrong_Override_Image_Key(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			ImageOverrides: v1alpha1.KeycloakRelatedImages{
				RHSSO: "quay.io/keycloak/keycloak:1.0.0",
			},
		},
	}

	// when
	returnedImage := GetKeycloakImage(cr)

	// then
	assert.Equal(t, returnedImage, KeycloakImage)
}

func TestUtil_Test_GetKeycloakImage_With_Override_Image_Set(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			ImageOverrides: v1alpha1.KeycloakRelatedImages{
				Keycloak: "quay.io/keycloak/keycloak:1.0.0",
			},
		},
	}

	// when
	returnedImage := GetKeycloakImage(cr)

	// then
	assert.Equal(t, returnedImage, "quay.io/keycloak/keycloak:1.0.0")
}

func Test_KeycloakDeployment_With_Overrided_Image(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			ImageOverrides: v1alpha1.KeycloakRelatedImages{
				Keycloak: "quay.io/keycloak/keycloak:1.0.0",
			},
		},
	}
	dbSecret := &v1.Secret{}
	currentImage := KeycloakImage

	// when
	for _, ele := range KeycloakDeployment(cr, dbSecret).Spec.Template.Spec.Containers {
		if ele.Name == KeycloakDeploymentName {
			currentImage = ele.Image
		}
	}

	// then
	assert.Equal(t, currentImage, "quay.io/keycloak/keycloak:1.0.0")
}

func Test_KeycloakDeploymentReconciled_With_Overrided_Image(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{}
	cr2 := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			ImageOverrides: v1alpha1.KeycloakRelatedImages{
				Keycloak: "quay.io/keycloak/keycloak:1.0.0",
			},
		},
	}
	dbSecret := &v1.Secret{}
	currentImage := KeycloakImage
	reconciledImage := KeycloakImage

	// when
	currentState := KeycloakDeployment(cr, dbSecret)
	for _, ele := range currentState.Spec.Template.Spec.Containers {
		if ele.Name == KeycloakDeploymentName {
			currentImage = ele.Image
		}
	}

	reconciledState := KeycloakDeploymentReconciled(cr2, currentState, dbSecret)
	for _, ele := range reconciledState.Spec.Template.Spec.Containers {
		if ele.Name == KeycloakDeploymentName {
			reconciledImage = ele.Image
		}
	}

	// then
	assert.Equal(t, currentImage, KeycloakImage)
	assert.Equal(t, reconciledImage, "quay.io/keycloak/keycloak:1.0.0")
}
