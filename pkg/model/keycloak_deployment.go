package model

import (
	"fmt"
	"strings"
	"os"

	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	v13 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	LivenessProbeInitialDelay	int32 = 30
	ReadinessProbeInitialDelay	int32 = 40
	//10s (curl) + 10s (curl) + 2s (just in case)
	ProbeTimeoutSeconds		int32 = 22
	ProbeTimeBetweenRunsSeconds	int32 = 30
	SuccessThresholdCount	int32 = 1
	FailureThresholdCount	int32 = 3
)

func GetServiceEnvVar(suffix string) string {
	serviceName := strings.ToUpper(PostgresqlServiceName)
	serviceName = strings.ReplaceAll(serviceName, "-", "_")
	return fmt.Sprintf("%v_%v", serviceName, suffix)
}

func GetContainerPorts() []v1.ContainerPort {
	return []v1.ContainerPort {
		{
			Name: "http",
			ContainerPort: KeycloakServicePort,
			Protocol:      "TCP",
		},
		{
			Name: "http-management",
			ContainerPort: KeycloakManagementPort,
			Protocol:      "TCP",
		},
		{
			Name: "http-monitoring",
			ContainerPort: 8778,
			Protocol:      "TCP",
		},
	}
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
			Name:  "X509_CA_BUNDLE",
			Value: "/var/run/secrets/kubernetes.io/serviceaccount/*.crt",
		},
		{
			Name:  "PROXY_ADDRESS_FORWARDING",
			Value: "true",
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

	for envKey, envValue := range cr.Spec.KeycloakDeploymentSpec.EnvVars {
		env = append(env, v1.EnvVar{
			Name:  envKey,
			Value: envValue,
		})
	}

	return env
}

func getCommand() []string {
	command := make([]string, 3)
	command[0] = "/bin/sh"
	command[1] = "-c"
	command[2] = "START=$(date +%s); while true; do STATUS=$(curl -s -o /dev/null -w '%{http_code}' http://localhost:15020/healthz/ready); if [ ${STATUS} -eq 200 ]; then exec /opt/jboss/tools/docker-entrypoint.sh -b 0.0.0.0; break; else END=$(date +%s); DIFF=$(( $END - $START )); if [ ${DIFF} -gt 300 ]; then curl -X POST http://127.0.0.1:15000/quitquitquit; break; else sleep 1; fi; fi; done;"
	return command
}

func getSecrets(cr *v1alpha1.Keycloak) []v1.LocalObjectReference {
	secret := os.Getenv("RELATED_IMAGE_PULL_SECRET")
	if secret == "" {
		return []v1.LocalObjectReference{}
	}
	return []v1.LocalObjectReference{
		{
			Name: secret,
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
					Annotations: cr.Spec.KeycloakDeploymentSpec.PodAnnotations,
				},
				Spec: v1.PodSpec{
					InitContainers: KeycloakExtensionsInitContainers(cr),
					Volumes:        KeycloakVolumes(),
					Containers: []v1.Container{
						{
							Name:  KeycloakDeploymentName,
							Image: Images.Images[KeycloakImage],
							Ports: GetContainerPorts(),
							VolumeMounts:   KeycloakVolumeMounts(KeycloakExtensionPath),
							LivenessProbe:  livenessProbe(cr),
							ReadinessProbe: readinessProbe(cr),
							Env:            getKeycloakEnv(cr, dbSecret),
							Resources:      getResources(cr),
							//Command:		getCommand(),
						},
					},
					ImagePullSecrets: getSecrets(cr),
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
	reconciled := currentState.DeepCopy()
	reconciled.ResourceVersion = currentState.ResourceVersion
	reconciled.Spec.Replicas = SanitizeNumberOfReplicas(cr.Spec.Instances, false)
	reconciled.Spec.Template.Spec.Volumes = KeycloakVolumes()
	reconciled.Spec.Template.Spec.Containers = []v1.Container{
		{
			Name:  KeycloakDeploymentName,
			Image: Images.Images[KeycloakImage],
			Ports: GetContainerPorts(),
			VolumeMounts:   KeycloakVolumeMounts(KeycloakExtensionPath),
			LivenessProbe:  livenessProbe(cr),
			ReadinessProbe: readinessProbe(cr),
			Env:            getKeycloakEnv(cr, dbSecret),
			Resources:      getResources(cr),
			//Command:		getCommand(),
		},
	}
	reconciled.Spec.Template.Spec.ImagePullSecrets = getSecrets(cr)
	reconciled.Spec.Template.Spec.InitContainers = KeycloakExtensionsInitContainers(cr)
	reconciled.Spec.Template.ObjectMeta.Annotations = cr.Spec.KeycloakDeploymentSpec.PodAnnotations
	return reconciled
}

func KeycloakVolumeMounts(extensionsPath string) []v1.VolumeMount {
	return []v1.VolumeMount{
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
}

func KeycloakVolumes() []v1.Volume {
	return []v1.Volume{
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
}

func livenessProbe(cr *v1alpha1.Keycloak) *v1.Probe {
	initialDelay := LivenessProbeInitialDelay
	if (cr.Spec.KeycloakDeploymentSpec.LivenessProbe.InitialDelaySeconds != 0) {
		initialDelay = cr.Spec.KeycloakDeploymentSpec.LivenessProbe.InitialDelaySeconds
	}

	period := ProbeTimeBetweenRunsSeconds
	if (cr.Spec.KeycloakDeploymentSpec.LivenessProbe.PeriodSeconds != 0) {
		period = cr.Spec.KeycloakDeploymentSpec.LivenessProbe.PeriodSeconds
	}

	timeout := ProbeTimeoutSeconds
	if (cr.Spec.KeycloakDeploymentSpec.LivenessProbe.TimeoutSeconds != 0) {
		timeout = cr.Spec.KeycloakDeploymentSpec.LivenessProbe.TimeoutSeconds
	}

	successThreshold := SuccessThresholdCount
	if (cr.Spec.KeycloakDeploymentSpec.LivenessProbe.SuccessThreshold != 0) {
		successThreshold = cr.Spec.KeycloakDeploymentSpec.LivenessProbe.SuccessThreshold
	}

	failureThreshold := FailureThresholdCount
	if (cr.Spec.KeycloakDeploymentSpec.LivenessProbe.FailureThreshold != 0) {
		failureThreshold = cr.Spec.KeycloakDeploymentSpec.LivenessProbe.FailureThreshold
	}

	handler := v1.Handler{
		Exec: &v1.ExecAction{
			Command: []string{
				"/bin/sh",
				"-c",
				"/probes/" + LivenessProbeProperty,
			},
		},
	}

	if (cr.Spec.KeycloakDeploymentSpec.LivenessProbe.Handler.Exec != nil ||
		cr.Spec.KeycloakDeploymentSpec.LivenessProbe.Handler.HTTPGet != nil ||
        cr.Spec.KeycloakDeploymentSpec.LivenessProbe.Handler.TCPSocket != nil) {
		handler = cr.Spec.KeycloakDeploymentSpec.LivenessProbe.Handler
	}

	return &v1.Probe{
		Handler:		handler,
		InitialDelaySeconds:	initialDelay,
		TimeoutSeconds:     timeout,
		PeriodSeconds:      period,
		SuccessThreshold:	successThreshold,
		FailureThreshold:	failureThreshold,
	}
}

func readinessProbe(cr *v1alpha1.Keycloak) *v1.Probe {
	initialDelay := ReadinessProbeInitialDelay
	if (cr.Spec.KeycloakDeploymentSpec.ReadinessProbe.InitialDelaySeconds != 0) {
		initialDelay = cr.Spec.KeycloakDeploymentSpec.ReadinessProbe.InitialDelaySeconds
	}

	period := ProbeTimeBetweenRunsSeconds
	if (cr.Spec.KeycloakDeploymentSpec.ReadinessProbe.PeriodSeconds != 0) {
		period = cr.Spec.KeycloakDeploymentSpec.ReadinessProbe.PeriodSeconds
	}

	timeout := ProbeTimeoutSeconds
	if (cr.Spec.KeycloakDeploymentSpec.ReadinessProbe.TimeoutSeconds != 0) {
		timeout = cr.Spec.KeycloakDeploymentSpec.ReadinessProbe.TimeoutSeconds
	}

	successThreshold := SuccessThresholdCount
	if (cr.Spec.KeycloakDeploymentSpec.ReadinessProbe.SuccessThreshold != 0) {
		successThreshold = cr.Spec.KeycloakDeploymentSpec.ReadinessProbe.SuccessThreshold
	}

	failureThreshold := FailureThresholdCount
	if (cr.Spec.KeycloakDeploymentSpec.ReadinessProbe.FailureThreshold != 0) {
		failureThreshold = cr.Spec.KeycloakDeploymentSpec.ReadinessProbe.FailureThreshold
	}

	handler := v1.Handler{
		Exec: &v1.ExecAction{
			Command: []string{
				"/bin/sh",
				"-c",
				"/probes/" + ReadinessProbeProperty,
			},
		},
	}

	if (cr.Spec.KeycloakDeploymentSpec.ReadinessProbe.Handler.Exec != nil ||
		cr.Spec.KeycloakDeploymentSpec.ReadinessProbe.Handler.HTTPGet != nil ||
        cr.Spec.KeycloakDeploymentSpec.ReadinessProbe.Handler.TCPSocket != nil) {
		handler = cr.Spec.KeycloakDeploymentSpec.ReadinessProbe.Handler
	}

	return &v1.Probe{
		Handler:		handler,
		InitialDelaySeconds:	initialDelay,
		TimeoutSeconds:		timeout,
		PeriodSeconds:		period,
		SuccessThreshold:	successThreshold,
		FailureThreshold:	failureThreshold,
	}
}
