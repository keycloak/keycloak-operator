package model

import (
	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	v13 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func PostgresqlAWSBackup(cr *v1alpha1.KeycloakBackup) *v13.Job {
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
					Containers:         postgresqlAwsBackupCommonContainers(cr),
					RestartPolicy:      v1.RestartPolicyNever,
					ServiceAccountName: PostgresqlBackupServiceAccountName,
				},
			},
		},
	}
}

func PostgresqlAWSBackupSelector(cr *v1alpha1.KeycloakBackup) client.ObjectKey {
	return client.ObjectKey{
		Name:      cr.Name,
		Namespace: cr.Namespace,
	}
}

func PostgresqlAWSBackupReconciled(cr *v1alpha1.KeycloakBackup, currentState *v13.Job) *v13.Job {
	reconciled := currentState.DeepCopy()
	reconciled.Spec.Template.Spec.Containers = postgresqlAwsBackupCommonContainers(cr)
	reconciled.Spec.Template.Spec.RestartPolicy = v1.RestartPolicyNever
	reconciled.Spec.Template.Spec.ServiceAccountName = PostgresqlBackupServiceAccountName
	return reconciled
}
