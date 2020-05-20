package model

import (
	"os"
	"testing"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestProfileManager_no_profiles_set(t *testing.T) {
	//when
	profileManager := NewProfileManager()

	//then
	assert.Equal(t, 0, len(profileManager.Profiles))
}

func TestProfileManager_no_profiles_with_nil_cr(t *testing.T) {
	//given
	profileManager := NewProfileManager()

	//when
	isRHSSO := profileManager.IsRHSSO(nil)

	//then
	assert.False(t, isRHSSO)
}

func TestProfileManager_profile_with_proper_cr(t *testing.T) {
	//given
	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			Profile: RHSSOProfile,
		},
	}
	profileManager := NewProfileManager()

	//when
	isRHSSO := profileManager.IsRHSSO(cr)

	//then
	assert.True(t, isRHSSO)
}

func TestProfileManager_profile_with_proper_cr_and_no_rhsso_profile(t *testing.T) {
	//given
	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			Profile: "SomeOtherProfile",
		},
	}
	profileManager := NewProfileManager()

	//when
	isRHSSO := profileManager.IsRHSSO(cr)

	//then
	assert.False(t, isRHSSO)
}

func TestProfileManager_profile_with_environmental_variables_set(t *testing.T) {
	//given
	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			Profile: "SomeOtherProfile",
		},
	}

	//when
	os.Setenv(ProfileEnvironmentalVariable, RHSSOProfile)
	profileManager := NewProfileManager()
	os.Unsetenv(ProfileEnvironmentalVariable)
	isRHSSO := profileManager.IsRHSSO(cr)

	//then
	assert.True(t, isRHSSO)
}

func TestProfileManager_RHSSO_image_with_cr(t *testing.T) {
	//given
	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			Profile: RHSSOProfile,
		},
	}

	//when
	profileManager := NewProfileManager()
	image := profileManager.GetKeycloakOrRHSSOImage(cr)

	//then
	assert.Equal(t, DefaultRHSSOImageOpenJDK, image)
}

func TestProfileManager_RHSSO_image_with_environmental_variable(t *testing.T) {
	//given
	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{},
	}

	//when
	os.Setenv(ProfileEnvironmentalVariable, RHSSOProfile)
	profileManager := NewProfileManager()
	os.Unsetenv(ProfileEnvironmentalVariable)
	image := profileManager.GetKeycloakOrRHSSOImage(cr)

	//then
	assert.Equal(t, DefaultRHSSOImageOpenJDK, image)
}

func TestProfileManager_get_keycloak_image_with_no_profile(t *testing.T) {
	//given
	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{},
	}

	//when
	profileManager := NewProfileManager()
	image := profileManager.GetKeycloakOrRHSSOImage(cr)

	//then
	assert.Equal(t, DefaultKeycloakImage, image)
}

func TestProfileManager_get_init_container_image_with_no_profile(t *testing.T) {
	//given
	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{},
	}

	//when
	profileManager := NewProfileManager()
	image := profileManager.GetInitContainerImage(cr)

	//then
	assert.Equal(t, DefaultKeycloakInitContainer, image)
}

func TestProfileManager_get_init_container_image_with_with_RHSSO_profile(t *testing.T) {
	//given
	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{},
	}

	//when
	os.Setenv(ProfileEnvironmentalVariable, RHSSOProfile)
	profileManager := NewProfileManager()
	os.Unsetenv(ProfileEnvironmentalVariable)
	image := profileManager.GetInitContainerImage(cr)

	//then
	assert.Equal(t, DefaultRHSSOInitContainer, image)
}
