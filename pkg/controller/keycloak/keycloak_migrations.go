package keycloak

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/common"
	"github.com/keycloak/keycloak-operator/pkg/model"
	v13 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var errBackup = errors.New("migrate backup fails")
var errNoMigrator = errors.New("migrator not found")

type Migrator interface {
	Migrate(cr *v1alpha1.Keycloak, currentState *common.ClusterState, desiredState common.DesiredClusterState) (common.DesiredClusterState, error)
}

type RecreateMigrator struct {
}

type RollingMigrator struct {
}

func GetMigrator(cr *v1alpha1.Keycloak) (Migrator, error) {
	switch cr.Spec.Migration.MigrationStrategy {
	case v1alpha1.NoStrategy, v1alpha1.StrategyRecreate:
		return &RecreateMigrator{}, nil
	case v1alpha1.StrategyRolling:
		return &RollingMigrator{}, nil
	default:
		return nil, errNoMigrator
	}
}

func (i *RecreateMigrator) Migrate(cr *v1alpha1.Keycloak, currentState *common.ClusterState, desiredState common.DesiredClusterState) (common.DesiredClusterState, error) {
	if needsMigration(cr, currentState) {
		desiredImage := model.Profiles.GetKeycloakOrRHSSOImage(cr)
		log.Info(fmt.Sprintf("Performing migration from '%s' to '%s'", currentState.KeycloakDeployment.Spec.Template.Spec.Containers[0].Image, desiredImage))
		deployment := findDeployment(&desiredState)

		// The backup should be made when Keycloak container is down.
		// This way, we minimize the chance of skipping important updated
		// that administrators/users could have made just before the backup.
		//
		// The current replicas status is checked here in desired state instead
		// of current state which is actually correct. KC controller creates
		// the desired state from current state, i.e. clones current replicas
		// status into the desired state.
		if deployment != nil && deployment.Status.Replicas > 0 {
			log.Info("Number of replicas decreased to 0")
			deployment.Spec.Replicas = &[]int32{0}[0]
			deployment.Spec.Template.Spec.Containers[0].Image = currentState.KeycloakDeployment.Spec.Template.Spec.Containers[0].Image
			return desiredState, nil
		}

		if cr.Spec.Migration.Backups.Enabled {
			var err error
			desiredState, err = oneTimeLocalBackup(cr, currentState, desiredState)
			if err != nil {
				return desiredState, err
			}
		}
	}

	return desiredState, nil
}

func (i *RollingMigrator) Migrate(cr *v1alpha1.Keycloak, currentState *common.ClusterState, desiredState common.DesiredClusterState) (common.DesiredClusterState, error) {
	return desiredState, nil
}

func needsMigration(cr *v1alpha1.Keycloak, currentState *common.ClusterState) bool {
	if currentState.KeycloakDeployment == nil {
		return false
	}
	deployedImage := currentState.KeycloakDeployment.Spec.Template.Spec.Containers[0].Image
	currentImage := model.Profiles.GetKeycloakOrRHSSOImage(cr)
	return deployedImage != currentImage
}

func findDeployment(desiredState *common.DesiredClusterState) *v13.StatefulSet {
	for _, v := range *desiredState {
		if (reflect.TypeOf(v) == reflect.TypeOf(common.GenericUpdateAction{})) {
			updateAction := v.(common.GenericUpdateAction)
			if (reflect.TypeOf(updateAction.Ref) == reflect.TypeOf(&v13.StatefulSet{})) {
				statefulSet := updateAction.Ref.(*v13.StatefulSet)
				if statefulSet.ObjectMeta.Name == model.KeycloakDeploymentName {
					return statefulSet
				}
			}
		}
	}
	return nil
}

func oneTimeLocalBackup(cr *v1alpha1.Keycloak, currentState *common.ClusterState, desiredState common.DesiredClusterState) (common.DesiredClusterState, error) {
	keycloakBackup := currentState.KeycloakBackup
	switch {
	case keycloakBackup == nil:
		backupCr := &v1alpha1.KeycloakBackup{}
		backupCr.Namespace = cr.Namespace
		backupCr.Name = model.MigrateBackupName + "-" + common.BackupTime
		labelSelect := metav1.LabelSelector{
			MatchLabels: cr.Labels,
		}
		backupCr.Spec.InstanceSelector = &labelSelect
		backupCr.Spec.StorageClassName = cr.Spec.StorageClassName

		migrationBackupCR := common.GenericCreateAction{
			Ref: model.KeycloakMigrationOneTimeBackup(backupCr),
			Msg: "Create Local Backup CR",
		}

		backupDesiredState := common.DesiredClusterState{}
		backupDesiredState = backupDesiredState.AddAction(migrationBackupCR)
		return backupDesiredState, nil
	case keycloakBackup.Status.Phase == v1alpha1.BackupPhaseCreated:
		log.Info("migrate backup succeeds")
		return desiredState, nil
	case keycloakBackup.Status.Phase == v1alpha1.BackupPhaseFailing:
		return nil, errBackup
	default:
		emptyDesiredState := common.DesiredClusterState{}
		log.Info("wait for migrate backup's creating")
		return emptyDesiredState, nil
	}
}
