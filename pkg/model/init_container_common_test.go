package model

import (
	"testing"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestUtil_Test_GetKeycloakInitContainerImage_Without_Override_Image_Set(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{}

	// when
	returnedImage := GetKeycloakInitContainerImage(cr)

	// then
	assert.Equal(t, returnedImage, KeycloakInitContainerImage)
}

func TestUtil_Test_GetKeycloakInitContainerImage_With_Wrong_Override_Image_Key(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			ImageOverrides: v1alpha1.KeycloakRelatedImages{
				Keycloak: "keycloak-init-container:1.0",
			},
		},
	}

	// when
	returnedImage := GetKeycloakInitContainerImage(cr)

	// then
	assert.Equal(t, returnedImage, KeycloakInitContainerImage)
}

func TestUtil_Test_GetKeycloakInitContainerImage_With_Override_Image_Set(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			ImageOverrides: v1alpha1.KeycloakRelatedImages{
				InitContainer: "keycloak-init-container:1.0",
			},
		},
	}

	// when
	returnedImage := GetKeycloakInitContainerImage(cr)

	// then
	assert.Equal(t, returnedImage, "keycloak-init-container:1.0")
}

func Test_KeycloakExtensionsInitContainers_With_Overrided_Image(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			ImageOverrides: v1alpha1.KeycloakRelatedImages{
				InitContainer: "keycloak-init-container:1.0",
			},
		},
	}

	// when
	currentImage := KeycloakExtensionsInitContainers(cr)[0].Image

	// then
	assert.Equal(t, currentImage, "keycloak-init-container:1.0")
}
