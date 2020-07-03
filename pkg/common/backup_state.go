package common

import (
	"context"

	kc "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/model"
	v12 "k8s.io/api/batch/v1"
	"k8s.io/api/batch/v1beta1"
	v1 "k8s.io/api/core/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type BackupState struct {
	LocalPersistentVolumeJob   *v12.Job
	LocalPersistentVolumeClaim *v1.PersistentVolumeClaim
	AwsJob                     *v12.Job
	AwsPeriodicJob             *v1beta1.CronJob
	Keycloak                   *kc.Keycloak
}

func NewBackupState(keycloak kc.Keycloak) *BackupState {
	return &BackupState{
		Keycloak: &keycloak,
	}
}

func NewClusterState() *ClusterState {
	return &ClusterState{}
}

func (i *BackupState) Read(context context.Context, cr *kc.KeycloakBackup, controllerClient client.Client) error {
	err := i.readLocalBackupJob(context, cr, controllerClient)
	if err != nil {
		return err
	}

	err = i.readLocalBackupPersistentVolumeClaim(context, cr, controllerClient)
	if err != nil {
		return err
	}

	err = i.readAwsBackupJob(context, cr, controllerClient)
	if err != nil {
		return err
	}

	err = i.readAwsPeriodicBackupJob(context, cr, controllerClient)
	if err != nil {
		return err
	}

	return err
}

func (i *BackupState) readLocalBackupJob(context context.Context, cr *kc.KeycloakBackup, controllerClient client.Client) error {
	// decide Job type first
	if cr.Spec.AWS.CredentialsSecretName != "" {
		return nil
	}

	localBackupJob := model.PostgresqlBackup(cr)
	localBackupJobSelector := model.PostgresqlBackupSelector(cr)

	err := controllerClient.Get(context, localBackupJobSelector, localBackupJob)
	if err != nil {
		if !apiErrors.IsNotFound(err) {
			return err
		}
	} else {
		i.LocalPersistentVolumeJob = localBackupJob
		cr.UpdateStatusSecondaryResources(i.LocalPersistentVolumeJob.Kind, i.LocalPersistentVolumeJob.Name)
	}
	return nil
}

func (i *BackupState) readLocalBackupPersistentVolumeClaim(context context.Context, cr *kc.KeycloakBackup, controllerClient client.Client) error {
	localBackupPersistentVolumeClaim := model.PostgresqlBackupPersistentVolumeClaim(cr)
	localBackupPersistentVolumeClaimSelector := model.PostgresqlBackupPersistentVolumeClaimSelector(cr)

	err := controllerClient.Get(context, localBackupPersistentVolumeClaimSelector, localBackupPersistentVolumeClaim)
	if err != nil {
		if !apiErrors.IsNotFound(err) {
			return err
		}
	} else {
		i.LocalPersistentVolumeClaim = localBackupPersistentVolumeClaim
		cr.UpdateStatusSecondaryResources(i.LocalPersistentVolumeClaim.Kind, i.LocalPersistentVolumeClaim.Name)
	}
	return nil
}

func (i *BackupState) readAwsBackupJob(context context.Context, cr *kc.KeycloakBackup, controllerClient client.Client) error {
	if cr.Spec.AWS.CredentialsSecretName == "" {
		return nil
	}

	awsBackupJob := model.PostgresqlAWSBackup(cr)
	awsBackupJobSelector := model.PostgresqlAWSBackupSelector(cr)

	err := controllerClient.Get(context, awsBackupJobSelector, awsBackupJob)
	if err != nil {
		if !apiErrors.IsNotFound(err) {
			return err
		}
	} else {
		i.AwsJob = awsBackupJob
		cr.UpdateStatusSecondaryResources(i.AwsJob.Kind, i.AwsJob.Name)
	}
	return nil
}

func (i *BackupState) readAwsPeriodicBackupJob(context context.Context, cr *kc.KeycloakBackup, controllerClient client.Client) error {
	awsPeriodicBackupJob := model.PostgresqlAWSPeriodicBackup(cr)
	awsPeriodicBackupJobSelector := model.PostgresqlAWSPeriodicBackupSelector(cr)

	err := controllerClient.Get(context, awsPeriodicBackupJobSelector, awsPeriodicBackupJob)
	if err != nil {
		if !apiErrors.IsNotFound(err) {
			return err
		}
	} else {
		i.AwsPeriodicJob = awsPeriodicBackupJob
		cr.UpdateStatusSecondaryResources(i.AwsPeriodicJob.Kind, i.AwsPeriodicJob.Name)
	}
	return nil
}

func (i *BackupState) IsResourcesReady() (bool, error) {
	switch {
	case i.AwsJob != nil:
		return IsJobReady(i.AwsJob)
	case i.LocalPersistentVolumeJob != nil:
		return IsJobReady(i.LocalPersistentVolumeJob)
	case i.AwsPeriodicJob != nil:
		// We don't manage readiness check for CronJobs
		return true, nil
	default:
		return false, nil
	}
}
