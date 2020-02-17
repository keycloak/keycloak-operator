package model

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	v13 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetServiceEnvVar(suffix string) string {
	serviceName := strings.ToUpper(PostgresqlServiceName)
	serviceName = strings.ReplaceAll(serviceName, "-", "_")
	return fmt.Sprintf("%v_%v", serviceName, suffix)
}

func getKeycloakEnv(cr *v1alpha1.Keycloak, dbSecret *v1.Secret) []v1.EnvVar {
	return []v1.EnvVar{
		// Database settings
		{
			Name:  "DB_VENDOR",
			Value: "POSTGRES",
		},
		{
			Name:  "DB_SCHEMA",
			Value: "public",
		},
		{
			Name:  "DB_ADDR",
			Value: PostgresqlServiceName + "." + cr.Namespace + ".svc.cluster.local",
		},
		{
			Name:  "DB_DATABASE",
			Value: GetExternalDatabaseName(dbSecret),
		},
		{
			Name: "DB_USER",
			ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: DatabaseSecretName,
					},
					Key: DatabaseSecretUsernameProperty,
				},
			},
		},
		{
			Name: "DB_PASSWORD",
			ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: DatabaseSecretName,
					},
					Key: DatabaseSecretPasswordProperty,
				},
			},
		},
		// Discovery settings
		{
			Name:  "NAMESPACE",
			Value: cr.Namespace,
		},
		{
			Name:  "JGROUPS_DISCOVERY_PROTOCOL",
			Value: "dns.DNS_PING",
		},
		{
			Name:  "JGROUPS_DISCOVERY_PROPERTIES",
			Value: "dns_query=" + KeycloakDiscoveryServiceName + "." + cr.Namespace + ".svc.cluster.local",
		},
		// Cache settings
		{
			Name:  "CACHE_OWNERS_COUNT",
			Value: "2",
		},
		{
			Name:  "CACHE_OWNERS_AUTH_SESSIONS_COUNT",
			Value: "2",
		},
		{
			Name: "KEYCLOAK_USER",
			ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: "credential-" + cr.Name,
					},
					Key: AdminUsernameProperty,
				},
			},
		},
		{
			Name: "KEYCLOAK_PASSWORD",
			ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: "credential-" + cr.Name,
					},
					Key: AdminPasswordProperty,
				},
			},
		},
		{
			Name:  GetServiceEnvVar("SERVICE_HOST"),
			Value: PostgresqlServiceName + "." + cr.Namespace + ".svc.cluster.local",
		},
		{
			Name:  GetServiceEnvVar("SERVICE_PORT"),
			Value: fmt.Sprintf("%v", GetExternalDatabasePort(dbSecret)),
		},
	}
}

func KeycloakDeployment(cr *v1alpha1.Keycloak, dbSecret *v1.Secret) *v13.StatefulSet {
	return &v13.StatefulSet{
		ObjectMeta: v12.ObjectMeta{
			Name:      KeycloakDeploymentName,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app":       ApplicationName,
				"component": KeycloakDeploymentComponent,
			},
		},
		Spec: v13.StatefulSetSpec{
			Replicas: SanitizeNumberOfReplicas(cr.Spec.Instances, true),
			Selector: &v12.LabelSelector{
				MatchLabels: map[string]string{
					"app":       ApplicationName,
					"component": KeycloakDeploymentComponent,
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: v12.ObjectMeta{
					Name:      KeycloakDeploymentName,
					Namespace: cr.Namespace,
					Labels: map[string]string{
						"app":       ApplicationName,
						"component": KeycloakDeploymentComponent,
					},
				},
				Spec: v1.PodSpec{
					InitContainers: KeycloakExtensionsInitContainers(cr),
					Volumes:        KeycloakVolumes(cr),
					Containers: []v1.Container{
						{
							Name:  KeycloakDeploymentName,
							Image: KeycloakImage,
							Ports: []v1.ContainerPort{
								{
									ContainerPort: KeycloakPodPort,
									Protocol:      "TCP",
								},
								{
									ContainerPort: KeycloakPodPortSSL,
									Protocol:      "TCP",
								},
								{
									ContainerPort: 9990,
									Protocol:      "TCP",
								},
								{
									ContainerPort: 8778,
									Protocol:      "TCP",
								},
							},
							VolumeMounts: KeycloakVolumeMounts(cr, KeycloakExtensionPath),
              Env:          getKeycloakEnv(cr, dbSecret),
							LivenessProbe: &v1.Probe{
								InitialDelaySeconds: 60,
								TimeoutSeconds:      1,
								Handler: v1.Handler{
									HTTPGet: &v1.HTTPGetAction{
										Path:   "/auth/realms/master",
										Port:   intstr.FromInt(8080),
										Scheme: "HTTP",
									},
								},
							},
							ReadinessProbe: &v1.Probe{
								TimeoutSeconds:      1,
								InitialDelaySeconds: 10,
								Handler: v1.Handler{
									HTTPGet: &v1.HTTPGetAction{
										Path:   "/auth/realms/master",
										Port:   intstr.FromInt(8080),
										Scheme: "HTTP",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func KeycloakDeploymentSelector(cr *v1alpha1.Keycloak) client.ObjectKey {
	return client.ObjectKey{
		Name:      KeycloakDeploymentName,
		Namespace: cr.Namespace,
	}
}

func KeycloakDeploymentReconciled(cr *v1alpha1.Keycloak, currentState *v13.StatefulSet, dbSecret *v1.Secret) *v13.StatefulSet {
	currentImage := GetCurrentKeycloakImage(currentState)

	reconciled := currentState.DeepCopy()
	reconciled.ResourceVersion = currentState.ResourceVersion
	reconciled.Spec.Replicas = SanitizeNumberOfReplicas(cr.Spec.Instances, false)
	reconciled.Spec.Template.Spec.Volumes = KeycloakVolumes(cr)
	reconciled.Spec.Template.Spec.Containers = []v1.Container{
		{
			Name:  KeycloakDeploymentName,
			Image: GetReconciledKeycloakImage(currentImage),
			Ports: []v1.ContainerPort{
				{
					ContainerPort: KeycloakPodPort,
					Protocol:      "TCP",
				},
				{
					ContainerPort: KeycloakPodPortSSL,
					Protocol:      "TCP",
				},
				{
					ContainerPort: 9990,
					Protocol:      "TCP",
				},
				{
					ContainerPort: 8778,
					Protocol:      "TCP",
				},
			},
			VolumeMounts: KeycloakVolumeMounts(cr, KeycloakExtensionPath),
			LivenessProbe: &v1.Probe{
				InitialDelaySeconds: 60,
				TimeoutSeconds:      1,
				Handler: v1.Handler{
					HTTPGet: &v1.HTTPGetAction{
						Path:   "/auth/realms/master",
						Port:   intstr.FromInt(8080),
						Scheme: "HTTP",
					},
				},
			},
			ReadinessProbe: &v1.Probe{
				TimeoutSeconds:      1,
				InitialDelaySeconds: 10,
				Handler: v1.Handler{
					HTTPGet: &v1.HTTPGetAction{
						Path:   "/auth/realms/master",
						Port:   intstr.FromInt(8080),
						Scheme: "HTTP",
					},
				},
			},
			Env: getKeycloakEnv(cr, dbSecret),
		},
	}
	reconciled.Spec.Template.Spec.InitContainers = KeycloakExtensionsInitContainers(cr)
	return reconciled
}

func KeycloakVolumeMounts(cr *v1alpha1.Keycloak, extensionsPath string) []v1.VolumeMount {

	volumeMounts := []v1.VolumeMount{
		{
			Name:      "keycloak-extensions",
			ReadOnly:  false,
			MountPath: extensionsPath,
		},
	}

	if !cr.Spec.ServingCertDisabled {
		volumeMounts = append(volumeMounts, v1.VolumeMount{
			Name:
			ServingCertSecretName,
			MountPath: "/etc/x509/https",
		})
	}

	return volumeMounts
}

func KeycloakVolumes(cr *v1alpha1.Keycloak) []v1.Volume {

	volumes := []v1.Volume{
		{
			Name: "keycloak-extensions",
			VolumeSource: v1.VolumeSource{
				EmptyDir: &v1.EmptyDirVolumeSource{},
			},
		},
	}

	if !cr.Spec.ServingCertDisabled {
		volumes = append(volumes, v1.Volume{
			Name: ServingCertSecretName,
			VolumeSource: v1.VolumeSource{
				Secret: &v1.SecretVolumeSource{
					SecretName: ServingCertSecretName,
					Optional:   &[]bool{true}[0],
				},
			},
		})
	}

	return volumes
}

// We allow the patch version of an image for keycloak to be increased outside of the operator on the cluster
func GetReconciledKeycloakImage(currentImage string) string {
	currentImageRepo, currentImageMajor, currentImageMinor, currentImagePatch := GetImageRepoAndVersion(currentImage)
	keycloakImageRepo, keycloakImageMajor, keycloakImageMinor, keycloakImagePatch := GetImageRepoAndVersion(KeycloakImage)

	// Need to convert the patch version strings to ints for a > comparison.
	currentImagePatchInt, err := strconv.Atoi(currentImagePatch)
	// If we are unable to convert to an int, always default to the operator image
	if err != nil {
		return KeycloakImage
	}

	// Need to convert the patch version strings to ints for a > comparison.
	keycloakImagePatchInt, err := strconv.Atoi(keycloakImagePatch)
	// If we are unable to convert to an int, always default to the operator image
	if err != nil {
		return KeycloakImage
	}

	// Check the repos, major and minor versions match. Check the cluster patch version is greater. If so, return and reconcile with the current cluster image
	// E.g. quay.io/keycloak/keycloak:7.0.1
	if currentImageRepo == keycloakImageRepo && currentImageMajor == keycloakImageMajor && currentImageMinor == keycloakImageMinor && currentImagePatchInt > keycloakImagePatchInt {
		return currentImage
	}

	return KeycloakImage
}
