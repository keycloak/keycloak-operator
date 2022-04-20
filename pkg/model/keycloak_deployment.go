package model

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	v13 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	LivenessProbeInitialDelay  = 30
	ReadinessProbeInitialDelay = 40
	//10s (curl) + 10s (curl) + 2s (just in case)
	ProbeTimeoutSeconds         = 22
	ProbeTimeBetweenRunsSeconds = 30
	ProbeFailureThreshold       = 10
)

func GetServiceEnvVar(suffix string) string {
	serviceName := strings.ToUpper(PostgresqlServiceName)
	serviceName = strings.ReplaceAll(serviceName, "-", "_")
	return fmt.Sprintf("%v_%v", serviceName, suffix)
}

func getResources(cr *v1alpha1.Keycloak) v1.ResourceRequirements {
	requirements := v1.ResourceRequirements{}
	requirements.Limits = v1.ResourceList{}
	requirements.Requests = v1.ResourceList{}

	cpu, err := resource.ParseQuantity(cr.Spec.KeycloakDeploymentSpec.Resources.Requests.Cpu().String())
	if err == nil && cpu.String() != "0" {
		requirements.Requests[v1.ResourceCPU] = cpu
	}

	memory, err := resource.ParseQuantity(cr.Spec.KeycloakDeploymentSpec.Resources.Requests.Memory().String())
	if err == nil && memory.String() != "0" {
		requirements.Requests[v1.ResourceMemory] = memory
	}

	cpu, err = resource.ParseQuantity(cr.Spec.KeycloakDeploymentSpec.Resources.Limits.Cpu().String())
	if err == nil && cpu.String() != "0" {
		requirements.Limits[v1.ResourceCPU] = cpu
	}
	memory, err = resource.ParseQuantity(cr.Spec.KeycloakDeploymentSpec.Resources.Limits.Memory().String())
	if err == nil && memory.String() != "0" {
		requirements.Limits[v1.ResourceMemory] = memory
	}

	return requirements
}

func getKeycloakEnv(cr *v1alpha1.Keycloak, dbSecret *v1.Secret) []v1.EnvVar {
	env := []v1.EnvVar{
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
			Value: PostgresqlServiceName + "." + cr.Namespace,
		},
		{
			Name:  "DB_PORT",
			Value: fmt.Sprintf("%v", GetExternalDatabasePort(dbSecret)),
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
			Value: "dns_query=" + KeycloakDiscoveryServiceName + "." + cr.Namespace,
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
			Name:  "X509_CA_BUNDLE",
			Value: "/var/run/secrets/kubernetes.io/serviceaccount/*.crt",
		},
		{
			Name:  "PROXY_ADDRESS_FORWARDING",
			Value: "true",
		},
		{
			Name:  "KEYCLOAK_STATISTICS",
			Value: "all",
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

	env = KeycloakSslEnvVariables(dbSecret, env)

	return env
}

func KeycloakSslEnvVariables(dbSecret *v1.Secret, env []v1.EnvVar) []v1.EnvVar {
	if dbSecret != nil {
		sslMode := string(dbSecret.Data[DatabaseSecretSslModeProperty])

		if sslMode != "" {
			dbParams := ""
			// is the deployment already having JDBC_PARAMS set ?
			for _, element := range env {
				if element.Name == KeycloakDatabaseConnectionParamsProperty {
					dbParams = element.Value + "&"
					break
				}
			}
			// append env variable
			env = append(env, v1.EnvVar{
				Name:  KeycloakDatabaseConnectionParamsProperty,
				Value: dbParams + "sslmode=" + sslMode + "&sslrootcert=" + KeycloakCertificatePath + "/root.crt",
			})
		}
	}
	return env
}

func KeycloakDeployment(cr *v1alpha1.Keycloak, dbSecret *v1.Secret, dbSSLSecret *v1.Secret) *v13.StatefulSet {
	labels := map[string]string{
		"app":       ApplicationName,
		"component": KeycloakDeploymentComponent,
	}
	podLabels := AddPodLabels(cr, labels)
	podAnnotations := cr.Spec.KeycloakDeploymentSpec.PodAnnotations
	keycloakStatefulset := &v13.StatefulSet{
		ObjectMeta: v12.ObjectMeta{
			Name:        KeycloakDeploymentName,
			Namespace:   cr.Namespace,
			Labels:      podLabels,
			Annotations: podAnnotations,
		},
		Spec: v13.StatefulSetSpec{
			Replicas: SanitizeNumberOfReplicas(cr.Spec.Instances, true),
			Selector: &v12.LabelSelector{
				MatchLabels: labels,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: v12.ObjectMeta{
					Name:        KeycloakDeploymentName,
					Namespace:   cr.Namespace,
					Labels:      podLabels,
					Annotations: podAnnotations,
				},
				Spec: v1.PodSpec{
					InitContainers: KeycloakExtensionsInitContainers(cr),
					Volumes:        KeycloakVolumes(cr, dbSSLSecret),
					Containers: []v1.Container{
						{
							Name:  KeycloakDeploymentName,
							Image: Images.Images[KeycloakImage],
							Ports: []v1.ContainerPort{
								{
									ContainerPort: KeycloakServicePort,
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
							VolumeMounts:   KeycloakVolumeMounts(cr, KeycloakExtensionPath, dbSSLSecret, KeycloakCertificatePath),
							LivenessProbe:  livenessProbe(),
							ReadinessProbe: readinessProbe(),
							Env:            getKeycloakEnv(cr, dbSecret),
							Args:           cr.Spec.KeycloakDeploymentSpec.Experimental.Args,
							Command:        cr.Spec.KeycloakDeploymentSpec.Experimental.Command,
							Resources:      getResources(cr),
						},
					},
					ServiceAccountName: getServiceAccountName(cr),
				},
			},
		},
	}

	if cr.Spec.KeycloakDeploymentSpec.Experimental.Affinity != nil {
		keycloakStatefulset.Spec.Template.Spec.Affinity = cr.Spec.KeycloakDeploymentSpec.Experimental.Affinity
	} else if cr.Spec.MultiAvailablityZones.Enabled {
		keycloakStatefulset.Spec.Template.Spec.Affinity = KeycloakPodAffinity(cr)
	}
	return keycloakStatefulset
}

func KeycloakDeploymentSelector(cr *v1alpha1.Keycloak) client.ObjectKey {
	return client.ObjectKey{
		Name:      KeycloakDeploymentName,
		Namespace: cr.Namespace,
	}
}

func KeycloakDeploymentReconciled(cr *v1alpha1.Keycloak, currentState *v13.StatefulSet, dbSecret *v1.Secret, dbSSLSecret *v1.Secret) *v13.StatefulSet {
	reconciled := currentState.DeepCopy()

	reconciled.ObjectMeta.Labels = AddPodLabels(cr, reconciled.ObjectMeta.Labels)
	reconciled.ObjectMeta.Annotations = AddPodAnnotations(cr, reconciled.ObjectMeta.Annotations)
	reconciled.Spec.Template.ObjectMeta.Labels = AddPodLabels(cr, reconciled.Spec.Template.ObjectMeta.Labels)
	reconciled.Spec.Template.ObjectMeta.Annotations = AddPodAnnotations(cr, reconciled.Spec.Template.ObjectMeta.Annotations)

	reconciled.ResourceVersion = currentState.ResourceVersion
	reconciled.Spec.Replicas = SanitizeNumberOfReplicas(cr.Spec.Instances, false)
	reconciled.Spec.Template.Spec.Volumes = KeycloakVolumes(cr, dbSSLSecret)
	reconciled.Spec.Template.Spec.Containers = []v1.Container{
		{
			Name:    KeycloakDeploymentName,
			Image:   Images.Images[KeycloakImage],
			Args:    cr.Spec.KeycloakDeploymentSpec.Experimental.Args,
			Command: cr.Spec.KeycloakDeploymentSpec.Experimental.Command,
			Ports: []v1.ContainerPort{
				{
					ContainerPort: KeycloakServicePort,
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
			VolumeMounts:   KeycloakVolumeMounts(cr, KeycloakExtensionPath, dbSSLSecret, KeycloakCertificatePath),
			LivenessProbe:  livenessProbe(),
			ReadinessProbe: readinessProbe(),
			Env:            getKeycloakEnv(cr, dbSecret),
			Resources:      getResources(cr),
		},
	}
	reconciled.Spec.Template.Spec.InitContainers = KeycloakExtensionsInitContainers(cr)
	if cr.Spec.KeycloakDeploymentSpec.Experimental.Affinity != nil {
		reconciled.Spec.Template.Spec.Affinity = cr.Spec.KeycloakDeploymentSpec.Experimental.Affinity
	}

	return reconciled
}

func KeycloakVolumeMounts(cr *v1alpha1.Keycloak, extensionsPath string, dbSSLSecret *v1.Secret, certificatePath string) []v1.VolumeMount {
	mountedVolumes := []v1.VolumeMount{
		{
			Name:      ServingCertSecretName,
			MountPath: "/etc/x509/https",
		},
		{
			Name:      "keycloak-extensions",
			ReadOnly:  false,
			MountPath: extensionsPath,
		},
		{
			Name:      KeycloakProbesName,
			MountPath: "/probes",
		},
	}

	if dbSSLSecret != nil {
		mountedVolumes = append(mountedVolumes, v1.VolumeMount{
			Name:      DatabaseSecretSslCert + "-vol",
			ReadOnly:  true,
			MountPath: certificatePath,
		})
	}

	mountedVolumes = addVolumeMountsFromKeycloakCR(cr, mountedVolumes)

	return mountedVolumes
}

func addVolumeMountsFromKeycloakCR(cr *v1alpha1.Keycloak, mountedVolumes []v1.VolumeMount) []v1.VolumeMount {
	if cr.Spec.KeycloakDeploymentSpec.Experimental.Volumes.Items != nil {
		for _, v := range cr.Spec.KeycloakDeploymentSpec.Experimental.Volumes.Items {
			volumeMapMount := v1.VolumeMount{
				Name:      v.Name,
				MountPath: v.MountPath,
			}
			mountedVolumes = append(mountedVolumes, volumeMapMount)
		}
	}
	return mountedVolumes
}

func KeycloakVolumes(cr *v1alpha1.Keycloak, dbSSLSecret *v1.Secret) []v1.Volume {
	volumes := []v1.Volume{
		{
			Name: ServingCertSecretName,
			VolumeSource: v1.VolumeSource{
				Secret: &v1.SecretVolumeSource{
					SecretName: ServingCertSecretName,
					Optional:   &[]bool{true}[0],
				},
			},
		},
		{
			Name: "keycloak-extensions",
			VolumeSource: v1.VolumeSource{
				EmptyDir: &v1.EmptyDirVolumeSource{},
			},
		},
		{
			Name: KeycloakProbesName,
			VolumeSource: v1.VolumeSource{
				ConfigMap: &v1.ConfigMapVolumeSource{
					LocalObjectReference: v1.LocalObjectReference{
						Name: KeycloakProbesName,
					},
					DefaultMode: &[]int32{0555}[0],
				},
			},
		},
	}
	if dbSSLSecret != nil {
		volumes = append(volumes, v1.Volume{
			Name: DatabaseSecretSslCert + "-vol",
			VolumeSource: v1.VolumeSource{
				Secret: &v1.SecretVolumeSource{
					SecretName: DatabaseSecretSslCert,
					Optional:   &[]bool{false}[0],
				},
			},
		})
	}

	volumes = addVolumesFromKeycloakCR(cr, volumes)

	return volumes
}

func addVolumesFromKeycloakCR(cr *v1alpha1.Keycloak, volumes []v1.Volume) []v1.Volume {
	if cr.Spec.KeycloakDeploymentSpec.Experimental.Volumes.Items != nil {
		for _, v := range cr.Spec.KeycloakDeploymentSpec.Experimental.Volumes.Items {
			var sources []v1.VolumeProjection
			if v.ConfigMaps != nil {
				for _, name := range v.ConfigMaps {
					sources = append(sources, v1.VolumeProjection{
						ConfigMap: &v1.ConfigMapProjection{
							LocalObjectReference: v1.LocalObjectReference{
								Name: name,
							},
							Items: v.Items,
						},
					})
				}
			}
			if v.Secrets != nil {
				for _, name := range v.Secrets {
					sources = append(sources, v1.VolumeProjection{
						Secret: &v1.SecretProjection{
							LocalObjectReference: v1.LocalObjectReference{
								Name: name,
							},
							Items: v.Items,
						},
					})
				}
			}

			mapVolume := v1.Volume{
				Name: v.Name,
				VolumeSource: v1.VolumeSource{
					Projected: &v1.ProjectedVolumeSource{
						Sources:     sources,
						DefaultMode: cr.Spec.KeycloakDeploymentSpec.Experimental.Volumes.DefaultMode,
					},
				},
			}
			volumes = append(volumes, mapVolume)
		}
	}
	return volumes
}

func livenessProbe() *v1.Probe {
	return &v1.Probe{
		Handler: v1.Handler{
			Exec: &v1.ExecAction{
				Command: []string{
					"/bin/sh",
					"-c",
					"/probes/" + LivenessProbeProperty,
				},
			},
		},
		InitialDelaySeconds: LivenessProbeInitialDelay,
		TimeoutSeconds:      ProbeTimeoutSeconds,
		PeriodSeconds:       ProbeTimeBetweenRunsSeconds,
		FailureThreshold:    ProbeFailureThreshold,
	}
}

func readinessProbe() *v1.Probe {
	return &v1.Probe{
		Handler: v1.Handler{
			Exec: &v1.ExecAction{
				Command: []string{
					"/bin/sh",
					"-c",
					"/probes/" + ReadinessProbeProperty,
				},
			},
		},
		InitialDelaySeconds: ReadinessProbeInitialDelay,
		TimeoutSeconds:      ProbeTimeoutSeconds,
		PeriodSeconds:       ProbeTimeBetweenRunsSeconds,
		FailureThreshold:    ProbeFailureThreshold,
	}
}

func KeycloakPodAffinity(cr *v1alpha1.Keycloak) *v1.Affinity {
	return &v1.Affinity{
		PodAntiAffinity: &v1.PodAntiAffinity{
			PreferredDuringSchedulingIgnoredDuringExecution: []v1.WeightedPodAffinityTerm{
				{
					Weight: 100,
					PodAffinityTerm: v1.PodAffinityTerm{
						LabelSelector: &v12.LabelSelector{
							MatchExpressions: []v12.LabelSelectorRequirement{
								{
									Key:      "app",
									Operator: "In",
									Values: []string{
										ApplicationName,
									},
								},
							},
						},
						TopologyKey: "topology.kubernetes.io/zone",
					},
				},
				{
					Weight: 90,
					PodAffinityTerm: v1.PodAffinityTerm{
						LabelSelector: &v12.LabelSelector{
							MatchExpressions: []v12.LabelSelectorRequirement{
								{
									Key:      "app",
									Operator: "In",
									Values: []string{
										ApplicationName,
									},
								},
							},
						},
						TopologyKey: "kubernetes.io/hostname",
					},
				},
			},
		},
	}
}

func getServiceAccountName(cr *v1alpha1.Keycloak) string {
	if cr.Spec.KeycloakDeploymentSpec.Experimental.ServiceAccountName == "" {
		return "default"
	}
	return cr.Spec.KeycloakDeploymentSpec.Experimental.ServiceAccountName
}
