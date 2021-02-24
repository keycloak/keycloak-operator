package model

import (
	"fmt"

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
			Value: "/var/run/secrets/kubernetes.io/serviceaccount/*.crt",
		},
		{
			Name:  "STATISTICS_ENABLED",
			Value: "TRUE",
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

	if len(cr.Spec.KeycloakDeploymentSpec.Experimental.Env) > 0 {
		// We override Keycloak pre-defined envs with what user specified. Not the other way around.
		env = MergeEnvs(cr.Spec.KeycloakDeploymentSpec.Experimental.Env, env)
	}

	return env
}

func RHSSODeployment(cr *v1alpha1.Keycloak, dbSecret *v1.Secret) *v13.StatefulSet {
	rhssoStatefulSet := &v13.StatefulSet{
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
					Volumes:        KeycloakVolumes(cr),
					InitContainers: KeycloakExtensionsInitContainers(cr),
					Affinity:       KeycloakPodAffinity(cr),
					Containers: []v1.Container{
						{
							Name:  KeycloakDeploymentName,
							Image: Images.Images[RHSSOImage],
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
							LivenessProbe:   livenessProbe(),
							ReadinessProbe:  readinessProbe(),
							Env:             getRHSSOEnv(cr, dbSecret),
							Args:            cr.Spec.KeycloakDeploymentSpec.Experimental.Args,
							Command:         cr.Spec.KeycloakDeploymentSpec.Experimental.Command,
							VolumeMounts:    KeycloakVolumeMounts(cr, RhssoExtensionPath),
							Resources:       getResources(cr),
							ImagePullPolicy: "Always",
						},
					},
				},
			},
		},
	}

	if cr.Spec.KeycloakDeploymentSpec.Experimental.Affinity != nil {
		rhssoStatefulSet.Spec.Template.Spec.Affinity = cr.Spec.KeycloakDeploymentSpec.Experimental.Affinity
	} else if cr.Spec.MultiAvailablityZones.Enabled {
		rhssoStatefulSet.Spec.Template.Spec.Affinity = KeycloakPodAffinity(cr)
	}
	return rhssoStatefulSet
}

func RHSSODeploymentSelector(cr *v1alpha1.Keycloak) client.ObjectKey {
	return client.ObjectKey{
		Name:      KeycloakDeploymentName,
		Namespace: cr.Namespace,
	}
}

func RHSSODeploymentReconciled(cr *v1alpha1.Keycloak, currentState *v13.StatefulSet, dbSecret *v1.Secret) *v13.StatefulSet {
	reconciled := currentState.DeepCopy()
	reconciled.ResourceVersion = currentState.ResourceVersion
	reconciled.Spec.Replicas = SanitizeNumberOfReplicas(cr.Spec.Instances, false)
	reconciled.Spec.Template.Spec.Volumes = KeycloakVolumes(cr)
	reconciled.Spec.Template.Spec.Containers = []v1.Container{
		{
			Name:    KeycloakDeploymentName,
			Image:   Images.Images[RHSSOImage],
			Args:    cr.Spec.KeycloakDeploymentSpec.Experimental.Args,
			Command: cr.Spec.KeycloakDeploymentSpec.Experimental.Command,
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
			VolumeMounts:    KeycloakVolumeMounts(cr, RhssoExtensionPath),
			LivenessProbe:   livenessProbe(),
			ReadinessProbe:  readinessProbe(),
			Env:             getRHSSOEnv(cr, dbSecret),
			Resources:       getResources(cr),
			ImagePullPolicy: "Always",
		},
	}
	reconciled.Spec.Template.Spec.InitContainers = KeycloakExtensionsInitContainers(cr)
	if cr.Spec.KeycloakDeploymentSpec.Experimental.Affinity != nil {
		reconciled.Spec.Template.Spec.Affinity = cr.Spec.KeycloakDeploymentSpec.Experimental.Affinity
	}

	return reconciled
}
