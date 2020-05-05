package model

import (
	"os"
	"runtime"
)

const (
	KeycloakImage         = "RELATED_IMAGE_KEYCLOAK"
	RHSSOImageOpenJ9      = "RELATED_IMAGE_RHSSO_OPENJ9"
	RHSSOImageOpenJDK     = "RELATED_IMAGE_RHSSO_OPENJDK"
	RHSSOImage            = "RELATED_IMAGE_RHSSO"
	KeycloakInitContainer = "RELATED_IMAGE_KEYCLOAK_INIT_CONTAINER"
	RHMIBackupContainer   = "RELATED_IMAGE_RHMI_BACKUP_CONTAINER"
	PostgresqlImage       = "RELATED_IMAGE_POSTGRESQL"

	DefaultKeycloakImage         = "quay.io/keycloak/keycloak:9.0.2"
	DefaultRHSSOImageOpenJ9      = "registry.redhat.io/rh-sso-7/sso74-openshift-rhel8:7.4-1"
	DefaultRHSSOImageOpenJDK     = "registry.redhat.io/rh-sso-7/sso74-openshift-rhel8:7.4-1"
	DefaultKeycloakInitContainer = "quay.io/keycloak/keycloak-init-container:master"
	DefaultRHMIBackupContainer   = "quay.io/integreatly/backup-container:1.0.14"
	DefaultPostgresqlImage       = "registry.access.redhat.com/rhscl/postgresql-10-rhel7:1"
)

var Images = NewImageManager()

type ImageManager struct {
	Images map[string]string
}

func NewImageManager() ImageManager {
	ret := ImageManager{}
	ret.Images = map[string]string{
		KeycloakImage:         ret.getImage(KeycloakImage, DefaultKeycloakImage),
		RHSSOImage:            ret.getRHSSOImage(),
		RHSSOImageOpenJ9:      ret.getImage(RHSSOImageOpenJ9, DefaultRHSSOImageOpenJ9),
		RHSSOImageOpenJDK:     ret.getImage(RHSSOImageOpenJDK, DefaultRHSSOImageOpenJDK),
		KeycloakInitContainer: ret.getImage(KeycloakInitContainer, DefaultKeycloakInitContainer),
		RHMIBackupContainer:   ret.getImage(RHMIBackupContainer, DefaultRHMIBackupContainer),
		PostgresqlImage:       ret.getImage(PostgresqlImage, DefaultPostgresqlImage),
	}
	return ret
}

func (p *ImageManager) getImage(environmentalVariable string, defaultValue string) string {
	env := os.Getenv(environmentalVariable)
	if env == "" {
		return defaultValue
	}
	return env
}

func (p *ImageManager) getRHSSOImage() string {
	// Full list of archs might be found here:
	// https://github.com/golang/go/blob/release-branch.go1.10/src/go/build/syslist.go#L8
	switch arch := runtime.GOARCH; arch {
	case "ppc64":
	case "ppc64le":
	case "s390x":
	case "s390":
		return p.getImage(RHSSOImageOpenJ9, DefaultRHSSOImageOpenJ9)
	default:
		return p.getImage(RHSSOImageOpenJDK, DefaultRHSSOImageOpenJDK)
	}
	panic("Unknown architecture")
}
