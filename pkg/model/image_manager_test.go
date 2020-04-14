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
