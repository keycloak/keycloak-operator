package model

import (
	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	v1 "k8s.io/api/core/v1"
)

func postgresqlAwsBackupCommonContainers(cr *v1alpha1.KeycloakBackup) []v1.Container {
	return []v1.Container{
		{
			Name:    cr.Name,
			Image:   Images.Images[RHMIBackupContainer],
			Command: []string{"/opt/intly/tools/entrypoint.sh", "-c", "postgres", "-n", cr.Namespace, "-b", "s3", "-e", ""},
			Env: []v1.EnvVar{
				{
					Name:  "BACKEND_SECRET_NAME",
					Value: cr.Spec.AWS.CredentialsSecretName,
				},
				{
					Name:  "BACKEND_SECRET_NAMESPACE",
					Value: cr.Namespace,
				},
				{
					Name:  "ENCRYPTION_SECRET_NAME",
					Value: cr.Spec.AWS.EncryptionKeySecretName,
				},
				{
					Name:  "COMPONENT_SECRET_NAME",
					Value: DatabaseSecretName,
				},
				{
					Name:  "COMPONENT_SECRET_NAMESPACE",
					Value: cr.Namespace,
				},
				{
					Name:  "PRODUCT_NAME",
					Value: "rhsso",
				},
			},
		},
	}
}
