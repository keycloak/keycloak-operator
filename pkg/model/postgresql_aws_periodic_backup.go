package model

import (
	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	v13 "k8s.io/api/batch/v1"
	"k8s.io/api/batch/v1beta1"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func PostgresqlAWSPeriodicBackup(cr *v1alpha1.KeycloakBackup) *v1beta1.CronJob {
	return &v1beta1.CronJob{
		ObjectMeta: v12.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app":       ApplicationName,
				"component": PostgresqlBackupComponent,
			},
		},
		Spec: v1beta1.CronJobSpec{
			Schedule: cr.Spec.AWS.Schedule,
			JobTemplate: v1beta1.JobTemplateSpec{
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
			},
		},
	}
}

func PostgresqlAWSPeriodicBackupSelector(cr *v1alpha1.KeycloakBackup) client.ObjectKey {
	return client.ObjectKey{
		Name:      cr.Name,
		Namespace: cr.Namespace,
	}
}

func PostgresqlAWSPeriodicBackupReconciled(cr *v1alpha1.KeycloakBackup, currentState *v1beta1.CronJob) *v1beta1.CronJob {
	reconciled := currentState.DeepCopy()
	reconciled.Spec.Schedule = cr.Spec.AWS.Schedule
	reconciled.Spec.JobTemplate.Spec.Template.Spec.Containers = postgresqlAwsBackupCommonContainers(cr)
	reconciled.Spec.JobTemplate.Spec.Template.Spec.RestartPolicy = v1.RestartPolicyNever
	reconciled.Spec.JobTemplate.Spec.Template.Spec.ServiceAccountName = PostgresqlBackupServiceAccountName
	return reconciled
}
