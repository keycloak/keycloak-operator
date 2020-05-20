package model

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImageManager_test_default_images(t *testing.T) {
	//when
	imageChooser := NewImageManager()

	//then
	assert.Equal(t, DefaultKeycloakImage, imageChooser.Images[KeycloakImage])
	assert.Equal(t, DefaultRHSSOImageOpenJ9, imageChooser.Images[RHSSOImageOpenJ9])
	assert.Equal(t, DefaultRHSSOImageOpenJDK, imageChooser.Images[RHSSOImageOpenJDK])
	assert.Equal(t, DefaultRHSSOImageOpenJDK, imageChooser.Images[RHSSOImage])
	assert.Equal(t, DefaultKeycloakInitContainer, imageChooser.Images[KeycloakInitContainer])
	assert.Equal(t, DefaultRHSSOInitContainer, imageChooser.Images[RHSSOInitContainer])
	assert.Equal(t, DefaultRHMIBackupContainer, imageChooser.Images[RHMIBackupContainer])
}

func TestImageManager_test_defining_image_using_environment_variable(t *testing.T) {
	//given
	os.Setenv(KeycloakImage, "test")

	//when
	imageChooser := NewImageManager()
	os.Unsetenv(KeycloakImage)

	//then
	assert.Equal(t, "test", imageChooser.Images[KeycloakImage])
}

func TestImageManager_test_overriding_rhsso_image_using_environment_variable(t *testing.T) {
	//given
	os.Setenv(RHSSOImage, "test")

	//when
	imageChooser := NewImageManager()
	os.Unsetenv(RHSSOImage)

	//then
	assert.Equal(t, "test", imageChooser.Images[RHSSOImage])
}

func TestImageManager_test_overriding_multiple_images_using_environment_variables(t *testing.T) {
	//given
	os.Setenv(RHSSOImage, "RHSSOImage")
	os.Setenv(RHSSOImageOpenJ9, "RHSSOImageOpenJ9")
	os.Setenv(RHSSOImageOpenJDK, "RHSSOImageOpenJDK")

	//when
	imageChooser := NewImageManager()
	getRHSSOImageResult := imageChooser.getRHSSOImage()
	os.Unsetenv(RHSSOImage)
	os.Unsetenv(RHSSOImageOpenJ9)
	os.Unsetenv(RHSSOImageOpenJDK)

	//then
	assert.Equal(t, "RHSSOImage", getRHSSOImageResult)
	assert.Equal(t, "RHSSOImage", imageChooser.Images[RHSSOImage])
	assert.Equal(t, "RHSSOImageOpenJ9", imageChooser.Images[RHSSOImageOpenJ9])
	assert.Equal(t, "RHSSOImageOpenJDK", imageChooser.Images[RHSSOImageOpenJDK])
}
