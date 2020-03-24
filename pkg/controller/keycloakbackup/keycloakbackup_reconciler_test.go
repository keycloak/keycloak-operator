package keycloakbackup

import (
	"testing"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/common"
	"github.com/keycloak/keycloak-operator/pkg/model"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/batch/v1"
	"k8s.io/api/batch/v1beta1"
	v12 "k8s.io/api/core/v1"
)

func TestKeycloakBackupReconciler_Test_Creating_Local_Backup_Job(t *testing.T) {
	// given
	cr := &v1alpha1.KeycloakBackup{}
	keycloak := v1alpha1.Keycloak{}

	currentState := common.NewBackupState(keycloak)

	// when
	reconciler := NewKeycloakBackupReconciler(keycloak)
	desiredState := reconciler.Reconcile(currentState, cr)

	// then
	assert.IsType(t, common.GenericCreateAction{}, desiredState[0])
	assert.IsType(t, common.GenericCreateAction{}, desiredState[1])
	assert.IsType(t, model.PostgresqlBackupPersistentVolumeClaim(cr), desiredState[0].(common.GenericCreateAction).Ref)
	assert.IsType(t, model.PostgresqlBackup(cr), desiredState[1].(common.GenericCreateAction).Ref)
}

func TestKeycloakBackupReconciler_Test_Updating_Local_Backup_Job(t *testing.T) {
	// given
	cr := &v1alpha1.KeycloakBackup{}
	keycloak := v1alpha1.Keycloak{}

	currentState := &common.BackupState{
		LocalPersistentVolumeJob:   &v1.Job{},
		LocalPersistentVolumeClaim: &v12.PersistentVolumeClaim{},
	}

	// when
	reconciler := NewKeycloakBackupReconciler(keycloak)
	desiredState := reconciler.Reconcile(currentState, cr)

	// then
	assert.IsType(t, common.GenericUpdateAction{}, desiredState[0])
	assert.IsType(t, common.GenericUpdateAction{}, desiredState[1])
	assert.IsType(t, model.PostgresqlBackupPersistentVolumeClaim(cr), desiredState[0].(common.GenericUpdateAction).Ref)
	assert.IsType(t, model.PostgresqlBackup(cr), desiredState[1].(common.GenericUpdateAction).Ref)
}

func TestKeycloakBackupReconciler_Test_Creating_AWS_Job(t *testing.T) {
	// given
	cr := &v1alpha1.KeycloakBackup{
		Spec: v1alpha1.KeycloakBackupSpec{
			AWS: v1alpha1.KeycloakAWSSpec{
				CredentialsSecretName: "aws-secret",
			},
		},
	}
	keycloak := v1alpha1.Keycloak{}

	currentState := &common.BackupState{
		AwsJob: &v1.Job{},
	}

	// when
	reconciler := NewKeycloakBackupReconciler(keycloak)
	desiredState := reconciler.Reconcile(currentState, cr)

	// then
	assert.IsType(t, common.GenericUpdateAction{}, desiredState[0])
	assert.IsType(t, model.PostgresqlAWSBackup(cr), desiredState[0].(common.GenericUpdateAction).Ref)
}

func TestKeycloakBackupReconciler_Test_Updating_AWS_Job(t *testing.T) {
	// given
	cr := &v1alpha1.KeycloakBackup{
		Spec: v1alpha1.KeycloakBackupSpec{
			AWS: v1alpha1.KeycloakAWSSpec{
				CredentialsSecretName: "aws-secret",
			},
		},
	}
	keycloak := v1alpha1.Keycloak{}

	currentState := common.NewBackupState(keycloak)

	// when
	reconciler := NewKeycloakBackupReconciler(keycloak)
	desiredState := reconciler.Reconcile(currentState, cr)

	// then
	assert.IsType(t, common.GenericCreateAction{}, desiredState[0])
	assert.IsType(t, model.PostgresqlAWSBackup(cr), desiredState[0].(common.GenericCreateAction).Ref)
}

func TestKeycloakBackupReconciler_Test_Creating_AWS_Periodic_Job(t *testing.T) {
	// given
	cr := &v1alpha1.KeycloakBackup{
		Spec: v1alpha1.KeycloakBackupSpec{
			AWS: v1alpha1.KeycloakAWSSpec{
				Schedule:              "*/2 * * * *",
				CredentialsSecretName: "aws-secret",
			},
		},
	}
	keycloak := v1alpha1.Keycloak{}

	currentState := common.NewBackupState(keycloak)

	// when
	reconciler := NewKeycloakBackupReconciler(keycloak)
	desiredState := reconciler.Reconcile(currentState, cr)

	// then
	assert.IsType(t, common.GenericCreateAction{}, desiredState[0])
	assert.IsType(t, model.PostgresqlAWSPeriodicBackup(cr), desiredState[0].(common.GenericCreateAction).Ref)
}

func TestKeycloakBackupReconciler_Test_Updating_AWS_Periodic_Job(t *testing.T) {
	// given
	cr := &v1alpha1.KeycloakBackup{
		Spec: v1alpha1.KeycloakBackupSpec{
			AWS: v1alpha1.KeycloakAWSSpec{
				Schedule:              "*/2 * * * *",
				CredentialsSecretName: "aws-secret",
			},
		},
	}
	keycloak := v1alpha1.Keycloak{}

	currentState := &common.BackupState{
		AwsPeriodicJob: &v1beta1.CronJob{},
	}

	// when
	reconciler := NewKeycloakBackupReconciler(keycloak)
	desiredState := reconciler.Reconcile(currentState, cr)

	// then
	assert.IsType(t, common.GenericUpdateAction{}, desiredState[0])
	assert.IsType(t, model.PostgresqlAWSPeriodicBackup(cr), desiredState[0].(common.GenericUpdateAction).Ref)
}
