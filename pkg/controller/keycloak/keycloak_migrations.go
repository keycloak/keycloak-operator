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
var errSelectorCantBeMigrated = errors.New("statefulSet Selector mismatch; please use Recreate migration strategy")

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
	deployment, deploymentIndex := findDeployment(&desiredState)

	// We can't modify existing selector on StatefulSet.
	// The selector might be wrongly set by e.g. RH-SSO 7.5.2.
	// In such case, we need to recreate the StatefulSet
	if needsStatefulSetRecreation(currentState, deployment) {
		log.Info("Detected StatefulSet has mismatching Selector")

		// First, gracefully scale down the cluster to make sure all the pods are gone (hence not deleting the SS right away).
		// This is to not mess with image upgrade that might be happening at the same time.
		// We need to make sure only one KC version is accessing the DB at a time.
		if deployment.Status.Replicas > 0 {
			scaleDownAndDontUpgradeImage(deployment, currentState)
			// Keep selector and labels intact for the moment
			deployment.Spec.Selector = currentState.KeycloakDeployment.Spec.Selector
			deployment.Spec.Template.Labels = currentState.KeycloakDeployment.Spec.Template.Labels
			return desiredState, nil
		}

		// Now, if scaled down, let's delete and recreate the SS.
		log.Info("Recreating the StatefulSet")
		deployment.ObjectMeta.ResourceVersion = "" // version can't be set when creating a resource
		deleteAction := common.GenericDeleteAction{
			Ref: deployment,
			Msg: "Delete server Deployment (StatefulSet)",
		}
		createAction := common.GenericCreateAction{
			Ref: deployment,
			Msg: "Recreate server Deployment (StatefulSet)",
		}
		newDesiredState := append(desiredState[:deploymentIndex], deleteAction, createAction)
		newDesiredState = append(newDesiredState, desiredState[deploymentIndex+1:]...)
		desiredState = newDesiredState // replace update with delete and create
		// no need to return now, we can let the DB backup to proceed
	}

	if needsImageMigration(cr, currentState) {
		desiredImage := model.Profiles.GetKeycloakOrRHSSOImage(cr)
		log.Info(fmt.Sprintf("Performing migration from '%s' to '%s'", currentState.KeycloakDeployment.Spec.Template.Spec.Containers[0].Image, desiredImage))

		// The backup should be made when Keycloak container is down.
		// This way, we minimize the chance of skipping important updated
		// that administrators/users could have made just before the backup.
		//
		// The current replicas status is checked here in desired state instead
		// of current state which is actually correct. KC controller creates
		// the desired state from current state, i.e. clones current replicas
		// status into the desired state.
		if deployment != nil && deployment.Status.Replicas > 0 {
			scaleDownAndDontUpgradeImage(deployment, currentState)
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
	deployment, _ := findDeployment(&desiredState)
	if needsStatefulSetRecreation(currentState, deployment) {
		return nil, errSelectorCantBeMigrated
	}
	return desiredState, nil
}

func needsImageMigration(cr *v1alpha1.Keycloak, currentState *common.ClusterState) bool {
	if currentState.KeycloakDeployment == nil {
		return false
	}
	deployedImage := currentState.KeycloakDeployment.Spec.Template.Spec.Containers[0].Image
	currentImage := model.Profiles.GetKeycloakOrRHSSOImage(cr)
	return deployedImage != currentImage
}

func needsStatefulSetRecreation(currentState *common.ClusterState, desiredDeployment *v13.StatefulSet) bool {
	if currentState.KeycloakDeployment == nil || desiredDeployment == nil {
		return false
	}
	// selectors can't be modified
	return !reflect.DeepEqual(currentState.KeycloakDeployment.Spec.Selector.MatchLabels, desiredDeployment.Spec.Selector.MatchLabels)
}

func findDeployment(desiredState *common.DesiredClusterState) (*v13.StatefulSet, int) {
	for i, v := range *desiredState {
		if (reflect.TypeOf(v) == reflect.TypeOf(common.GenericUpdateAction{})) {
			updateAction := v.(common.GenericUpdateAction)
			if (reflect.TypeOf(updateAction.Ref) == reflect.TypeOf(&v13.StatefulSet{})) {
				statefulSet := updateAction.Ref.(*v13.StatefulSet)
				if statefulSet.ObjectMeta.Name == model.KeycloakDeploymentName {
					return statefulSet, i
				}
			}
		}
	}
	return nil, 0
}

func scaleDownAndDontUpgradeImage(deployment *v13.StatefulSet, currentState *common.ClusterState) {
	log.Info("Number of replicas decreased to 0")
	deployment.Spec.Replicas = &[]int32{0}[0]
	deployment.Spec.Template.Spec.Containers[0].Image = currentState.KeycloakDeployment.Spec.Template.Spec.Containers[0].Image
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
