package model

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	v13 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func getRHSSOEnv(cr *v1alpha1.Keycloak, dbSecret *v1.Secret) []v1.EnvVar {
	var env = []v1.EnvVar{
		// Database settings
		{
			Name:  "DB_SERVICE_PREFIX_MAPPING",
			Value: PostgresqlServiceName + "=DB",
		},
		{
			Name:  "TX_DATABASE_PREFIX_MAPPING",
			Value: PostgresqlServiceName + "=DB",
		},
		{
			Name:  "DB_JNDI",
			Value: "java:jboss/datasources/KeycloakDS",
		},
		{
			Name:  "DB_SCHEMA",
			Value: "public",
		},
		{
			Name: "DB_USERNAME",
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
		{
			Name:  "DB_DATABASE",
			Value: GetExternalDatabaseName(dbSecret),
		},
		// Discovery settings
		{
			Name:  "JGROUPS_PING_PROTOCOL",
			Value: "dns.DNS_PING",
		},
		{
			Name:  "OPENSHIFT_DNS_PING_SERVICE_NAME",
			Value: KeycloakDiscoveryServiceName + "." + cr.Namespace + ".svc.cluster.local",
		},
		{
			Name: "SSO_ADMIN_USERNAME",
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
			Name: "SSO_ADMIN_PASSWORD",
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
			Name:  "X509_CA_BUNDLE",
			Value: "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt",
		},
	}

	if cr.Spec.ExternalDatabase.Enabled {
		env = append(env, v1.EnvVar{
			Name:  GetServiceEnvVar("SERVICE_HOST"),
			Value: PostgresqlServiceName + "." + cr.Namespace + ".svc.cluster.local",
		})
		env = append(env, v1.EnvVar{
			Name:  GetServiceEnvVar("SERVICE_PORT"),
			Value: fmt.Sprintf("%v", GetExternalDatabasePort(dbSecret)),
		})
	}

	return env
}

func RHSSODeployment(cr *v1alpha1.Keycloak, dbSecret *v1.Secret) *v13.StatefulSet {
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
					Volumes:        KeycloakVolumes(),
					InitContainers: KeycloakExtensionsInitContainers(cr),
					Containers: []v1.Container{
						{
							Name:  KeycloakDeploymentName,
							Image: RHSSOImage,
							Ports: []v1.ContainerPort{
								{
									ContainerPort: KeycloakServicePort,
									Protocol:      "TCP",
								},
								{
									ContainerPort: 8080,
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
							LivenessProbe:  livenessProbe(),
							ReadinessProbe: readinessProbe(),
							Env:            getRHSSOEnv(cr, dbSecret),
							VolumeMounts:   KeycloakVolumeMounts(RhssoExtensionPath),
						},
					},
				},
			},
		},
	}
}

func RHSSODeploymentSelector(cr *v1alpha1.Keycloak) client.ObjectKey {
	return client.ObjectKey{
		Name:      KeycloakDeploymentName,
		Namespace: cr.Namespace,
	}
}

func RHSSODeploymentReconciled(cr *v1alpha1.Keycloak, currentState *v13.StatefulSet, dbSecret *v1.Secret) *v13.StatefulSet {
	currentImage := GetCurrentKeycloakImage(currentState)

	reconciled := currentState.DeepCopy()
	reconciled.ResourceVersion = currentState.ResourceVersion
	reconciled.Spec.Replicas = SanitizeNumberOfReplicas(cr.Spec.Instances, false)
	reconciled.Spec.Template.Spec.Volumes = KeycloakVolumes()
	reconciled.Spec.Template.Spec.Containers = []v1.Container{
		{
			Name:  KeycloakDeploymentName,
			Image: GetReconciledRHSSOImage(currentImage),
			Ports: []v1.ContainerPort{
				{
					ContainerPort: KeycloakServicePort,
					Protocol:      "TCP",
				},
				{
					ContainerPort: 8080,
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
			VolumeMounts:   KeycloakVolumeMounts(RhssoExtensionPath),
			LivenessProbe:  livenessProbe(),
			ReadinessProbe: readinessProbe(),
			Env:            getRHSSOEnv(cr, dbSecret),
		},
	}
	reconciled.Spec.Template.Spec.InitContainers = KeycloakExtensionsInitContainers(cr)

	return reconciled
}

// We allow the patch version of an image for RH-SSO to be increased outside of the operator on the cluster
func GetReconciledRHSSOImage(currentImage string) string {
	currentImageRepo, currentImageMajor, currentImageMinor, currentImagePatch := GetImageRepoAndVersion(currentImage)
	RHSSOImageRepo, RHSSOImageMajor, RHSSOImageMinor, RHSSOImagePatch := GetImageRepoAndVersion(RHSSOImage)

	// Since the minor version of the RHSSO image should always be 0-X, we can ignore all before the -
	currentImageMinorStrings := strings.Split(currentImageMinor, "-")
	if len(currentImageMinorStrings) > 1 {
		currentImageMinor = currentImageMinorStrings[1]
	}
	RHSSOImageMinorStrings := strings.Split(RHSSOImageMinor, "-")
	if len(RHSSOImageMinor) > 1 {
		RHSSOImageMinor = RHSSOImageMinorStrings[1]
	}

	// Need to convert the minor and patch version strings to ints for a > comparison.
	// If we are unable to convert to an int, default to the operator image
	currentImageMinorInt, err := strconv.Atoi(currentImageMinor)
	if err != nil {
		return RHSSOImage
	}

	RHSSOImageMinorInt, err := strconv.Atoi(RHSSOImageMinor)
	if err != nil {
		return RHSSOImage
	}

	// The patch version usually doesn't exist so we can ignore the error
	currentImagePatchInt, _ := strconv.Atoi(currentImagePatch)
	RHSSOImagePatchInt, _ := strconv.Atoi(RHSSOImagePatch)

	// All tags in the RHSSO image repo are technically patch versions
	// Image repo should match, the "major" version should match which is always "1". If the minor or patch tag versions have increased on the cluster, we will allow it ot reconcile this image
	// E.g. registry.access.redhat.com/redhat-sso-7/sso73-openshift:1.0-15 or registry.access.redhat.com/redhat-sso-7/sso73-openshift:1.0-15.1567588155
	if currentImageRepo == RHSSOImageRepo && currentImageMajor == RHSSOImageMajor && (currentImageMinorInt > RHSSOImageMinorInt || currentImagePatchInt > RHSSOImagePatchInt) {
		return currentImage
	}

	return RHSSOImage
}
