package model

import (
	"testing"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestUtil_Test_GetBackupImage_Without_Override_Image_Set(t *testing.T) {
	// given
	keycloak := &v1alpha1.Keycloak{}

	// when
	returnedImage := GetBackupImage(keycloak)

	// then
	assert.Equal(t, returnedImage, BackupImage)
}

func TestUtil_Test_GetBackupImage_With_Override_Image_Set(t *testing.T) {
	// given
	keycloak := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			ImageOverrides: v1alpha1.KeycloakRelatedImages{
				Backup: "backup-container:1.0",
			},
		},
	}

	// when
	returnedImage := GetBackupImage(keycloak)

	// then
	assert.Equal(t, returnedImage, "backup-container:1.0")
}

func Test_postgresqlAwsBackupCommonContainers_With_Overrided_Image(t *testing.T) {
	// given
	cr := &v1alpha1.KeycloakBackup{}
	keycloak := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			ImageOverrides: v1alpha1.KeycloakRelatedImages{
				Backup: "backup-container:1.0",
			},
		},
	}

	// when
	currentImage := postgresqlAwsBackupCommonContainers(cr, keycloak)[0].Image

	// then
	assert.Equal(t, currentImage, "backup-container:1.0")
}
