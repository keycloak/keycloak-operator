package model

import (
	"strings"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	v1 "k8s.io/api/core/v1"
)

func KeycloakExtensionsInitContainers(cr *v1alpha1.Keycloak) []v1.Container {
	return []v1.Container{
		{
			Name:  "extensions-init",
			Image: Profiles.GetInitContainerImage(cr),
			Env: []v1.EnvVar{
				{
					Name:  KeycloakExtensionEnvVar,
					Value: strings.Join(cr.Spec.Extensions, ","),
				},
			},
			VolumeMounts: []v1.VolumeMount{
				{
					Name:      "keycloak-extensions",
					ReadOnly:  false,
					MountPath: KeycloakExtensionsInitContainerPath,
				},
			},
			TerminationMessagePath:   "/dev/termination-log",
			TerminationMessagePolicy: "File",
			ImagePullPolicy:          "Always",
		},
	}
}
