package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtil_Test_GetReconciledRHSSOImage_With_No_Image(t *testing.T) {
	// given
	currentImage := ""

	// when
	reconciledImage := GetReconciledRHSSOImage(currentImage)

	// then
	assert.Equal(t, reconciledImage, RHSSOImage)
}

func TestUtil_Test_GetReconciledRHSSOImage_With_Random_Image(t *testing.T) {
	// given
	currentImage := "not/real/image:1.1.2"

	// when
	reconciledImage := GetReconciledRHSSOImage(currentImage)

	// then
	assert.Equal(t, reconciledImage, RHSSOImage)
}

func TestUtil_Test_GetReconciledRHSSOImage_With_No_Change(t *testing.T) {
	// given
	currentImage := RHSSOImage

	// when
	reconciledImage := GetReconciledRHSSOImage(currentImage)

	// then
	assert.Equal(t, reconciledImage, RHSSOImage)
}

func TestUtil_Test_GetReconciledRHSSOImage_With_Lower_Version(t *testing.T) {
	// given
	currentImage := "registry.access.redhat.com/redhat-sso-6/sso62-openshift:1.0-1"

	// when
	reconciledImage := GetReconciledRHSSOImage(currentImage)

	// then
	assert.Equal(t, reconciledImage, RHSSOImage)
}

func TestUtil_Test_GetReconciledRHSSOImage_With_Higher_Major_Version(t *testing.T) {
	// given
	currentImage := "registry.access.redhat.com/redhat-sso-8/sso83-openshift:1.0-15"

	// when
	reconciledImage := GetReconciledRHSSOImage(currentImage)

	// then
	assert.Equal(t, reconciledImage, RHSSOImage)
}

func TestUtil_Test_GetReconciledRHSSOImage_With_Higher_Minor_Version(t *testing.T) {
	// given
	currentImage := "registry.access.redhat.com/redhat-sso-7/sso74-openshift:1.0-15"

	// when
	reconciledImage := GetReconciledRHSSOImage(currentImage)

	// then
	assert.Equal(t, reconciledImage, RHSSOImage)
}

func TestUtil_Test_GetReconciledRHSSOImage_With_Higher_Patch_Version(t *testing.T) {
	// given
	currentImage := RHSSOImage[:len(RHSSOImage)-1] + "11"

	// when
	reconciledImage := GetReconciledRHSSOImage(currentImage)

	// then
	assert.Equal(t, reconciledImage, currentImage)
}

func TestUtil_Test_GetReconciledRHSSOImage_With_Higher_CVE_Patch_Version(t *testing.T) {
	// given
	currentImage := RHSSOImage + ".1"

	// when
	reconciledImage := GetReconciledRHSSOImage(currentImage)

	// then
	assert.Equal(t, reconciledImage, currentImage)
}
