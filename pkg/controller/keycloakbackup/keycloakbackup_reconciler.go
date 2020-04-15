package keycloakbackup

import (
	kc "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/common"
	"github.com/keycloak/keycloak-operator/pkg/model"
)

type Reconciler interface {
	Reconcile(cr *kc.KeycloakBackup) (common.DesiredClusterState, error)
}

type KeycloakBackupReconciler struct { // nolint
	Keycloak kc.Keycloak
}

func NewKeycloakBackupReconciler(keycloak kc.Keycloak) *KeycloakBackupReconciler {
	return &KeycloakBackupReconciler{
		Keycloak: keycloak,
	}
}

func (i *KeycloakBackupReconciler) Reconcile(currentState *common.BackupState, cr *kc.KeycloakBackup) common.DesiredClusterState {
	desired := common.DesiredClusterState{}

	if cr.Spec.AWS != (kc.KeycloakAWSSpec{}) {
		if cr.Spec.AWS.Schedule == "" {
			desired = desired.AddAction(i.GetAwsBackupDesiredState(currentState, cr))
		} else {
			desired = desired.AddAction(i.GetAwsPeriodicBackupDesiredState(currentState, cr))
		}
	} else {
		desired = desired.AddAction(i.GetLocalBackupPersistentVolumeDesiredState(currentState, cr))
		desired = desired.AddAction(i.GetLocalBackupDesiredState(currentState, cr))
	}

	return desired
}

func (i *KeycloakBackupReconciler) GetAwsPeriodicBackupDesiredState(currentState *common.BackupState, cr *kc.KeycloakBackup) common.ClusterAction {
	if currentState.AwsPeriodicJob == nil {
		return common.GenericCreateAction{
			Ref: model.PostgresqlAWSPeriodicBackup(cr),
			Msg: "Create AWS Periodic Backup job",
		}
	}

	return common.GenericUpdateAction{
		Ref: model.PostgresqlAWSPeriodicBackupReconciled(cr, currentState.AwsPeriodicJob),
		Msg: "Update AWS Periodic Backup job",
	}
}

func (i *KeycloakBackupReconciler) GetAwsBackupDesiredState(currentState *common.BackupState, cr *kc.KeycloakBackup) common.ClusterAction {
	if currentState.AwsJob == nil {
		return common.GenericCreateAction{
			Ref: model.PostgresqlAWSBackup(cr),
			Msg: "Create AWS Backup job",
		}
	}

	return common.GenericUpdateAction{
		Ref: model.PostgresqlAWSBackupReconciled(cr, currentState.AwsJob),
		Msg: "Update AWS Backup job",
	}
}

func (i *KeycloakBackupReconciler) GetLocalBackupDesiredState(currentState *common.BackupState, cr *kc.KeycloakBackup) common.ClusterAction {
	if currentState.LocalPersistentVolumeJob == nil {
		return common.GenericCreateAction{
			Ref: model.PostgresqlBackup(cr),
			Msg: "Create Local Backup job",
		}
	}

	return common.GenericUpdateAction{
		Ref: model.PostgresqlBackupReconciled(cr, currentState.LocalPersistentVolumeJob),
		Msg: "Update Local Backup job",
	}
}

func (i *KeycloakBackupReconciler) GetLocalBackupPersistentVolumeDesiredState(currentState *common.BackupState, cr *kc.KeycloakBackup) common.ClusterAction {
	if currentState.LocalPersistentVolumeJob == nil {
		return common.GenericCreateAction{
			Ref: model.PostgresqlBackupPersistentVolumeClaim(cr),
			Msg: "Create Local Backup Persistent Volume Claim",
		}
	}

	return common.GenericUpdateAction{
		Ref: model.PostgresqlBackupPersistentVolumeClaimReconciled(cr, currentState.LocalPersistentVolumeClaim),
		Msg: "Update Local Backup Persistent Volume Claim",
	}
}
