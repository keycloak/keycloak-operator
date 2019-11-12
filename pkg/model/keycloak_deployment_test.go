package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtil_Test_GetReconciledKeycloakImage_With_No_Image(t *testing.T) {
	// given
	currentImage := ""

	// when
	reconciledImage := GetReconciledKeycloakImage(currentImage)

	// then
	assert.Equal(t, reconciledImage, KeycloakImage)
}

func TestUtil_Test_GetReconciledKeycloakImage_With_Random_Image(t *testing.T) {
	// given
	currentImage := "not/real/image:1.1.1"

	// when
	reconciledImage := GetReconciledKeycloakImage(currentImage)

	// then
	assert.Equal(t, reconciledImage, KeycloakImage)
}

func TestUtil_Test_GetReconciledKeycloakImage_With_No_Change(t *testing.T) {
	// given
	currentImage := KeycloakImage

	// when
	reconciledImage := GetReconciledKeycloakImage(currentImage)

	// then
	assert.Equal(t, reconciledImage, KeycloakImage)
}

func TestUtil_Test_GetReconciledKeycloakImage_With_Lower_Version(t *testing.T) {
	// given
	currentImage := "quay.io/keycloak/keycloak:6.0.0"

	// when
	reconciledImage := GetReconciledKeycloakImage(currentImage)

	// then
	assert.Equal(t, reconciledImage, KeycloakImage)
}

func TestUtil_Test_GetReconciledKeycloakImage_With_Higher_Major_Version(t *testing.T) {
	// given
	currentImage := "quay.io/keycloak/keycloak:100.0.1"

	// when
	reconciledImage := GetReconciledKeycloakImage(currentImage)

	// then
	assert.Equal(t, reconciledImage, KeycloakImage)
}

func TestUtil_Test_GetReconciledKeycloakImage_With_Higher_Minor_Version(t *testing.T) {
	// given
	currentImage := "quay.io/keycloak/keycloak:100.100.1"

	// when
	reconciledImage := GetReconciledKeycloakImage(currentImage)

	// then
	assert.Equal(t, reconciledImage, KeycloakImage)
}

func TestUtil_Test_GetReconciledKeycloakImage_With_Higher_Patch_Version(t *testing.T) {
	// given
	currentImage := KeycloakImage[:len(KeycloakImage)-1] + "100"

	// when
	reconciledImage := GetReconciledKeycloakImage(currentImage)

	// then
	assert.Equal(t, reconciledImage, currentImage)
}
