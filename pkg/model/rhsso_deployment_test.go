package model

import (
	"testing"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
)

func TestUtil_Test_GetRHSSOImage_Without_Override_Image_Set(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{}

	// when
	returnedImage := GetRHSSOImage(cr)

	// then
	assert.Equal(t, returnedImage, RHSSOImage)
}

func TestUtil_Test_GetRHSSOImage_With_Wrong_Override_Image_Key(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			ImageOverrides: v1alpha1.KeycloakRelatedImages{
				Keycloak: "sso-openshift:1.0",
			},
		},
	}

	// when
	returnedImage := GetRHSSOImage(cr)

	// then
	assert.Equal(t, returnedImage, RHSSOImage)
}

func TestUtil_Test_GetRHSSOImage_With_Override_Image_Set(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			ImageOverrides: v1alpha1.KeycloakRelatedImages{
				RHSSO: "sso-openshift:1.0",
			},
		},
	}

	// when
	returnedImage := GetRHSSOImage(cr)

	// then
	assert.Equal(t, returnedImage, "sso-openshift:1.0")
}

func Test_RHSSODeployment_With_Overrided_Image(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			ImageOverrides: v1alpha1.KeycloakRelatedImages{
				RHSSO: "sso-openshift:1.0",
			},
		},
	}
	dbSecret := &v1.Secret{}
	currentImage := RHSSOImage

	// when
	for _, ele := range RHSSODeployment(cr, dbSecret).Spec.Template.Spec.Containers {
		if ele.Name == KeycloakDeploymentName {
			currentImage = ele.Image
		}
	}

	// then
	assert.Equal(t, currentImage, "sso-openshift:1.0")
}

func Test_RHSSODeploymentReconciled_With_Overrided_Image(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{}
	cr2 := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			ImageOverrides: v1alpha1.KeycloakRelatedImages{
				RHSSO: "sso-openshift:1.0",
			},
		},
	}
	dbSecret := &v1.Secret{}
	currentImage := RHSSOImage
	reconciledImage := RHSSOImage

	// when
	currentState := RHSSODeployment(cr, dbSecret)
	for _, ele := range currentState.Spec.Template.Spec.Containers {
		if ele.Name == KeycloakDeploymentName {
			currentImage = ele.Image
		}
	}

	reconciledState := RHSSODeploymentReconciled(cr2, currentState, dbSecret)
	for _, ele := range reconciledState.Spec.Template.Spec.Containers {
		if ele.Name == KeycloakDeploymentName {
			reconciledImage = ele.Image
		}
	}

	// then
	assert.Equal(t, currentImage, RHSSOImage)
	assert.Equal(t, reconciledImage, "sso-openshift:1.0")
}
