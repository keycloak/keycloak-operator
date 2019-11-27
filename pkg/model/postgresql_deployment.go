package model

import (
	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	v13 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func PostgresqlDeployment(cr *v1alpha1.Keycloak) *v13.Deployment {
	return &v13.Deployment{
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
							Image: PostgresqlImage,
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
											"psql -h 127.0.0.1 -U $POSTGRES_USER -q -d $POSTGRES_DB -c 'SELECT 1'",
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
									Name: "POSTGRES_USER",
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
									Name: "PGUSER",
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
									Name: "POSTGRES_PASSWORD",
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
									Name:  "POSTGRES_DB",
									Value: PostgresqlDatabase,
								},
								{
									// Due to permissions issue, we need to create a subdirectory in the PVC
									Name:  "PGDATA",
									Value: "/var/lib/postgresql/data/pgdata",
								},
							},
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      PostgresqlPersistentVolumeName,
									MountPath: "/var/lib/postgresql/data",
								},
							},
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
			Image: PostgresqlImage,
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
							"psql -h 127.0.0.1 -U $POSTGRES_USER -q -d $POSTGRES_DB -c 'SELECT 1'",
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
					Name: "POSTGRES_USER",
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
					Name: "PGUSER",
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
					Name: "POSTGRES_PASSWORD",
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
					Name:  "POSTGRES_DB",
					Value: PostgresqlDatabase,
				},
				{
					// Due to permissions issue, we need to create a subdirectory in the PVC
					Name:  "PGDATA",
					Value: "/var/lib/postgresql/data/pgdata",
				},
			},
			VolumeMounts: []v1.VolumeMount{
				{
					Name:      PostgresqlPersistentVolumeName,
					MountPath: "/var/lib/postgresql/data",
				},
			},
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
