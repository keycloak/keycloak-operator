package model

import (
	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	v13 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func getPostgresResources(cr *v1alpha1.Keycloak) v1.ResourceRequirements {
	requirements := v1.ResourceRequirements{}
	requirements.Limits = v1.ResourceList{}
	requirements.Requests = v1.ResourceList{}

	cpu, err := resource.ParseQuantity(cr.Spec.PostgresDeploymentSpec.Resources.Requests.Cpu().String())
	if err == nil && cpu.String() != "0" {
		requirements.Requests[v1.ResourceCPU] = cpu
	}
	memory, err := resource.ParseQuantity(cr.Spec.PostgresDeploymentSpec.Resources.Requests.Memory().String())
	if err == nil && memory.String() != "0" {
		requirements.Requests[v1.ResourceMemory] = memory
	}

	cpu, err = resource.ParseQuantity(cr.Spec.PostgresDeploymentSpec.Resources.Limits.Cpu().String())
	if err == nil && cpu.String() != "0" {
		requirements.Limits[v1.ResourceCPU] = cpu
	}
	memory, err = resource.ParseQuantity(cr.Spec.PostgresDeploymentSpec.Resources.Limits.Memory().String())
	if err == nil && memory.String() != "0" {
		requirements.Limits[v1.ResourceMemory] = memory
	}
	return requirements
}

func PostgresqlDeployment(cr *v1alpha1.Keycloak, isOpenshift bool) *v13.Deployment {
	v13Deployment := &v13.Deployment{
		ObjectMeta: v12.ObjectMeta{
			Name:      PostgresqlDeploymentName,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app":       ApplicationName,
				"component": PostgresqlDeploymentComponent,
			},
		},
		Spec: v13.DeploymentSpec{
			Selector: &v12.LabelSelector{
				MatchLabels: map[string]string{
					"app":       ApplicationName,
					"component": PostgresqlDeploymentComponent,
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: v12.ObjectMeta{
					Name:      PostgresqlDeploymentName,
					Namespace: cr.Namespace,
					Labels: map[string]string{
						"app":       ApplicationName,
						"component": PostgresqlDeploymentComponent,
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  PostgresqlDeploymentName,
							Image: Images.Images[PostgresqlImage],
							Ports: []v1.ContainerPort{
								{
									ContainerPort: 5432,
									Protocol:      "TCP",
								},
							},
							ReadinessProbe: &v1.Probe{
								TimeoutSeconds:      1,
								InitialDelaySeconds: 5,
								Handler: v1.Handler{
									Exec: &v1.ExecAction{
										Command: []string{
											"/bin/sh",
											"-c",
											"psql -h 127.0.0.1 -U $POSTGRESQL_USER -q -d $POSTGRESQL_DATABASE -c 'SELECT 1'",
										},
									},
								},
							},
							LivenessProbe: &v1.Probe{
								InitialDelaySeconds: 30,
								TimeoutSeconds:      1,
								Handler: v1.Handler{
									TCPSocket: &v1.TCPSocketAction{
										Port: intstr.FromInt(5432),
									},
								},
							},
							Env: []v1.EnvVar{
								{
									Name: "POSTGRESQL_USER",
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
									Name: "POSTGRESQL_PASSWORD",
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
									Name:  "POSTGRESQL_DATABASE",
									Value: PostgresqlDatabase,
								},
							},
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      PostgresqlPersistentVolumeName,
									MountPath: PostgresqlPersistentVolumeMountPath,
								},
							},
							Resources: getPostgresResources(cr),
						},
					},
					Volumes: []v1.Volume{
						{
							Name: PostgresqlPersistentVolumeName,
							VolumeSource: v1.VolumeSource{
								PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
									ClaimName: PostgresqlPersistentVolumeName,
								},
							},
						},
					},
				},
			},
			Strategy: v13.DeploymentStrategy{
				Type: v13.RecreateDeploymentStrategyType,
			},
		},
	}

	if !isOpenshift {
		v13Deployment.Spec.Template.Spec.InitContainers = getPostgresqlDeploymentInitContainer(cr)
	}
	return v13Deployment
}

func getPostgresqlDeploymentInitContainer(cr *v1alpha1.Keycloak) []v1.Container {
	return []v1.Container{
		{
			Name:  "init-pvc",
			Image: Images.Images[PostgresqlImage],
			SecurityContext: &v1.SecurityContext{
				RunAsUser: pointer.Int64Ptr(0),
			},
			Command: []string{
				"sh",
				"-c",
				"chown -R postgres:postgres " + PostgresqlPersistentVolumeMountPath,
			},
			VolumeMounts: []v1.VolumeMount{
				{
					Name:      PostgresqlPersistentVolumeName,
					MountPath: PostgresqlPersistentVolumeMountPath,
				},
			},
		},
	}
}

func PostgresqlDeploymentSelector(cr *v1alpha1.Keycloak) client.ObjectKey {
	return client.ObjectKey{
		Name:      PostgresqlDeploymentName,
		Namespace: cr.Namespace,
	}
}

func PostgresqlDeploymentReconciled(cr *v1alpha1.Keycloak, currentState *v13.Deployment) *v13.Deployment {
	reconciled := currentState.DeepCopy()
	reconciled.ResourceVersion = currentState.ResourceVersion
	reconciled.Spec.Strategy = v13.DeploymentStrategy{
		Type: v13.RecreateDeploymentStrategyType,
	}
	reconciled.Spec.Template.Spec.Containers = []v1.Container{
		{
			Name:  PostgresqlDeploymentName,
			Image: Images.Images[PostgresqlImage],
			Ports: []v1.ContainerPort{
				{
					ContainerPort: 5432,
					Protocol:      "TCP",
				},
			},
			ReadinessProbe: &v1.Probe{
				TimeoutSeconds:      1,
				InitialDelaySeconds: 5,
				Handler: v1.Handler{
					Exec: &v1.ExecAction{
						Command: []string{
							"/bin/sh",
							"-c",
							"psql -h 127.0.0.1 -U $POSTGRESQL_USER -q -d $POSTGRESQL_DATABASE -c 'SELECT 1'",
						},
					},
				},
			},
			LivenessProbe: &v1.Probe{
				InitialDelaySeconds: 30,
				TimeoutSeconds:      1,
				Handler: v1.Handler{
					TCPSocket: &v1.TCPSocketAction{
						Port: intstr.FromInt(5432),
					},
				},
			},
			Env: []v1.EnvVar{
				{
					Name: "POSTGRESQL_USER",
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
					Name: "POSTGRESQL_PASSWORD",
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
					Name:  "POSTGRESQL_DATABASE",
					Value: PostgresqlDatabase,
				},
			},
			VolumeMounts: []v1.VolumeMount{
				{
					Name:      PostgresqlPersistentVolumeName,
					MountPath: PostgresqlPersistentVolumeMountPath,
				},
			},
			Resources: getPostgresResources(cr),
		},
	}
	reconciled.Spec.Template.Spec.Volumes = []v1.Volume{
		{
			Name: PostgresqlPersistentVolumeName,
			VolumeSource: v1.VolumeSource{
				PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
					ClaimName: PostgresqlPersistentVolumeName,
				},
			},
		},
	}
	return reconciled
}
