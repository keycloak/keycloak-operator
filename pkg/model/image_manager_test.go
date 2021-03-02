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
	assert.Equal(t, DefaultKeycloakImage, imageChooser[KeycloakImage].Image)
	assert.Equal(t, DefaultRHSSOImageOpenJ9, imageChooser[RHSSOImageOpenJ9].Image)
	assert.Equal(t, DefaultRHSSOImageOpenJDK, imageChooser[RHSSOImageOpenJDK].Image)
	assert.Equal(t, DefaultRHSSOImageOpenJDK, imageChooser[RHSSOImage].Image)
	assert.Equal(t, DefaultKeycloakInitContainer, imageChooser[KeycloakInitContainer].Image)
	assert.Equal(t, DefaultRHSSOInitContainer, imageChooser[RHSSOInitContainer].Image)
	assert.Equal(t, DefaultRHMIBackupContainer, imageChooser[RHMIBackupContainer].Image)
}

func TestImageManager_test_defining_image_using_environment_variable(t *testing.T) {
	//given
	os.Setenv(KeycloakImage, "test")

	//when
	imageChooser := NewImageManager()
	os.Unsetenv(KeycloakImage)

	//then
	assert.Equal(t, "test", imageChooser[KeycloakImage].Image)
}

func TestImageManager_test_defining_image_pull_secrets_using_environment_variables(t *testing.T) {
	//given
	os.Setenv(keycloakImageIPS, "keycloakImagePullSecret")
	os.Setenv(keycloakInitContainerIPS, "keycloakInitContainerImagePullSecret")
	os.Setenv(rhmiBackupContainerIPS, "rhmiBackupContainerImagePullSecret")
	os.Setenv(postgresqlImageIPS, "postgresqlImagePullSecret")

	//when
	imageChooser := NewImageManager()
	os.Unsetenv(keycloakImageIPS)
	os.Unsetenv(keycloakInitContainerIPS)
	os.Unsetenv(rhmiBackupContainerIPS)
	os.Unsetenv(postgresqlImageIPS)

	//then
	assert.Equal(t, "keycloakImagePullSecret", imageChooser[KeycloakImage].ImagePullSecret.Name)
	assert.Equal(t, "keycloakInitContainerImagePullSecret", imageChooser[KeycloakInitContainer].ImagePullSecret.Name)
	assert.Equal(t, "rhmiBackupContainerImagePullSecret", imageChooser[RHMIBackupContainer].ImagePullSecret.Name)
	assert.Equal(t, "postgresqlImagePullSecret", imageChooser[PostgresqlImage].ImagePullSecret.Name)
}

func TestImageManager_test_overriding_rhsso_image_using_environment_variable(t *testing.T) {
	//given
	os.Setenv(RHSSOImage, "test")

	//when
	imageChooser := NewImageManager()
	os.Unsetenv(RHSSOImage)

	//then
	assert.Equal(t, "test", imageChooser[RHSSOImage].Image)
}

func TestImageManager_test_overriding_multiple_images_using_environment_variables(t *testing.T) {
	//given
	os.Setenv(RHSSOImage, "RHSSOImage")
	os.Setenv(RHSSOImageOpenJ9, "RHSSOImageOpenJ9")
	os.Setenv(RHSSOImageOpenJDK, "RHSSOImageOpenJDK")

	//when
	imageChooser := NewImageManager()
	getRHSSOImageResult := getRHSSOImage()
	os.Unsetenv(RHSSOImage)
	os.Unsetenv(RHSSOImageOpenJ9)
	os.Unsetenv(RHSSOImageOpenJDK)

	//then
	assert.Equal(t, "RHSSOImage", getRHSSOImageResult)
	assert.Equal(t, "RHSSOImage", imageChooser[RHSSOImage].Image)
	assert.Equal(t, "RHSSOImageOpenJ9", imageChooser[RHSSOImageOpenJ9].Image)
	assert.Equal(t, "RHSSOImageOpenJDK", imageChooser[RHSSOImageOpenJDK].Image)
}
