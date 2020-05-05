package model

import (
	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	v13 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func PostgresqlBackup(cr *v1alpha1.KeycloakBackup) *v13.Job {
	return &v13.Job{
		ObjectMeta: v12.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app":       ApplicationName,
				"component": PostgresqlBackupComponent,
			},
		},
		Spec: v13.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Volumes: []v1.Volume{
						{
							Name: PostgresqlBackupPersistentVolumeName + "-" + cr.Name,
							VolumeSource: v1.VolumeSource{
								PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
									ClaimName: PostgresqlBackupPersistentVolumeName + "-" + cr.Name,
								},
							},
						},
					},
					Containers: []v1.Container{
						{
							Name:    cr.Name,
							Image:   Images.Images[PostgresqlImage],
							Command: []string{"/bin/sh", "-c"},
							Args:    []string{"pg_dump $POSTGRES_DB | tee /backup/backup.sql"},
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
									Name: "PGPASSWORD",
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
									Name:  "PGHOST",
									Value: PostgresqlServiceName,
								},
							},
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      PostgresqlBackupPersistentVolumeName + "-" + cr.Name,
									MountPath: "/backup",
								},
							},
						},
					},
					RestartPolicy:      v1.RestartPolicyNever,
					ServiceAccountName: PostgresqlBackupServiceAccountName,
				},
			},
		},
	}
}

func PostgresqlBackupSelector(cr *v1alpha1.KeycloakBackup) client.ObjectKey {
	return client.ObjectKey{
		Name:      cr.Name,
		Namespace: cr.Namespace,
	}
}

func PostgresqlBackupReconciled(cr *v1alpha1.KeycloakBackup, currentState *v13.Job) *v13.Job {
	reconciled := currentState.DeepCopy()
	reconciled.Spec.Template.Spec.Volumes = []v1.Volume{
		{
			Name: PostgresqlBackupPersistentVolumeName + "-" + cr.Name,
			VolumeSource: v1.VolumeSource{
				PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
					ClaimName: PostgresqlBackupPersistentVolumeName + "-" + cr.Name,
				},
			},
		},
	}
	reconciled.Spec.Template.Spec.Containers = []v1.Container{
		{
			Name:    cr.Name,
			Image:   Images.Images[PostgresqlImage],
			Command: []string{"/bin/sh", "-c"},
			Args:    []string{"pg_dump $POSTGRES_DB | tee /backup/backup.sql"},
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
					Name: "PGPASSWORD",
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
					Name:  "PGHOST",
					Value: PostgresqlServiceName,
				},
			},
			VolumeMounts: []v1.VolumeMount{
				{
					Name:      PostgresqlBackupPersistentVolumeName + "-" + cr.Name,
					MountPath: "/backup",
				},
			},
		},
	}
	reconciled.Spec.Template.Spec.RestartPolicy = v1.RestartPolicyNever
	reconciled.Spec.Template.Spec.ServiceAccountName = PostgresqlBackupServiceAccountName
	return reconciled
}
