package model

import (
	"os"
	"runtime"

	v1 "k8s.io/api/core/v1"
)

const (
	KeycloakImage         = "RELATED_IMAGE_KEYCLOAK"
	RHSSOImageOpenJ9      = "RELATED_IMAGE_RHSSO_OPENJ9"
	RHSSOImageOpenJDK     = "RELATED_IMAGE_RHSSO_OPENJDK"
	RHSSOImage            = "RELATED_IMAGE_RHSSO"
	KeycloakInitContainer = "RELATED_IMAGE_KEYCLOAK_INIT_CONTAINER"
	RHSSOInitContainer    = "RELATED_IMAGE_RHSSO_INIT_CONTAINER"
	RHMIBackupContainer   = "RELATED_IMAGE_RHMI_BACKUP_CONTAINER"
	PostgresqlImage       = "RELATED_IMAGE_POSTGRESQL"

	DefaultKeycloakImage         = "quay.io/keycloak/keycloak:latest"
	DefaultRHSSOImageOpenJ9      = "registry.redhat.io/rh-sso-7/sso74-openj9-openshift-rhel8:7.4"
	DefaultRHSSOImageOpenJDK     = "registry.redhat.io/rh-sso-7/sso74-openshift-rhel8:7.4"
	DefaultKeycloakInitContainer = "quay.io/keycloak/keycloak-init-container:master"
	DefaultRHSSOInitContainer    = "registry.redhat.io/rh-sso-7-tech-preview/sso74-init-container-rhel8:7.4"
	DefaultRHMIBackupContainer   = "quay.io/integreatly/backup-container:1.0.16"
	DefaultPostgresqlImage       = "registry.access.redhat.com/rhscl/postgresql-10-rhel7:1"

	keycloakImageIPS         = "KEYCLOAK_IMAGE_PULL_SECRET"
	keycloakInitContainerIPS = "KEYCLOAK_INIT_CONTAINER_IMAGE_PULL_SECRET"
	rhmiBackupContainerIPS   = "RHMI_BACKUP_CONTAINER_IMAGE_PULL_SECRET"
	postgresqlImageIPS       = "POSTGRESQL_IMAGE_PULL_SECRET"
)

var Images = NewImageManager()

type ImageManager map[string]Image

type Image struct {
	Image           string
	ImagePullSecret v1.LocalObjectReference
}

func NewImageManager() ImageManager {
	ret := ImageManager{
		KeycloakImage: {
			Image:           getImage(KeycloakImage, DefaultKeycloakImage),
			ImagePullSecret: getImagePullSecret(keycloakImageIPS),
		},
		RHSSOImage:        {Image: getRHSSOImage()},
		RHSSOImageOpenJ9:  {Image: getImage(RHSSOImageOpenJ9, DefaultRHSSOImageOpenJ9)},
		RHSSOImageOpenJDK: {Image: getImage(RHSSOImageOpenJDK, DefaultRHSSOImageOpenJDK)},
		KeycloakInitContainer: {
			Image:           getImage(KeycloakInitContainer, DefaultKeycloakInitContainer),
			ImagePullSecret: getImagePullSecret(keycloakInitContainerIPS),
		},
		RHSSOInitContainer: {Image: getImage(RHSSOInitContainer, DefaultRHSSOInitContainer)},
		RHMIBackupContainer: {
			Image:           getImage(RHMIBackupContainer, DefaultRHMIBackupContainer),
			ImagePullSecret: getImagePullSecret(rhmiBackupContainerIPS),
		},
		PostgresqlImage: {
			Image:           getImage(PostgresqlImage, DefaultPostgresqlImage),
			ImagePullSecret: getImagePullSecret(postgresqlImageIPS),
		},
	}
	return ret
}

func getImage(environmentalVariable string, defaultValue string) string {
	env := os.Getenv(environmentalVariable)
	if env == "" {
		return defaultValue
	}
	return env
}

func getRHSSOImage() string {
	defaultImage := getDefaultRHSSOImageForCurrentArchitecture()
	return getImage(RHSSOImage, defaultImage)
}

func getDefaultRHSSOImageForCurrentArchitecture() string {
	// Full list of archs might be found here:
	// https://github.com/golang/go/blob/release-branch.go1.10/src/go/build/syslist.go#L8
	switch arch := runtime.GOARCH; arch {
	case "ppc64", "ppc64le", "s390x", "s390":
		return getImage(RHSSOImageOpenJ9, DefaultRHSSOImageOpenJ9)
	default:
		return getImage(RHSSOImageOpenJDK, DefaultRHSSOImageOpenJDK)
	}
}

func getImagePullSecret(environmentalVariable string) v1.LocalObjectReference {
	return v1.LocalObjectReference{Name: os.Getenv(environmentalVariable)}
}
