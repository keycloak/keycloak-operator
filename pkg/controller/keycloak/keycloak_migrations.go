package keycloak

import (
	"fmt"
	"reflect"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/common"
	"github.com/keycloak/keycloak-operator/pkg/model"
	v13 "k8s.io/api/apps/v1"
)

type Migrator interface {
	Migrate(cr *v1alpha1.Keycloak, currentState *common.ClusterState, desiredState common.DesiredClusterState) (common.DesiredClusterState, error)
}

type DefaultMigrator struct {
}

func NewDefaultMigrator() *DefaultMigrator {
	return &DefaultMigrator{}
}

func (i *DefaultMigrator) Migrate(cr *v1alpha1.Keycloak, currentState *common.ClusterState, desiredState common.DesiredClusterState) (common.DesiredClusterState, error) {
	if needsMigration(cr, currentState) {
		desiredImage := model.KeycloakImage
		if cr.Spec.Profile == common.RHSSOProfile {
			desiredImage = model.RHSSOImage
		}
		log.Info(fmt.Sprintf("Performing migration from '%s' to '%s'", currentState.KeycloakDeployment.Spec.Template.Spec.Containers[0].Image, desiredImage))
		deployment := findDeployment(desiredState)
		if deployment != nil {
			log.Info("Number of replicas decreased to 1")
			deployment.Spec.Replicas = &[]int32{1}[0]
		}
	}

	return desiredState, nil
}

func needsMigration(cr *v1alpha1.Keycloak, currentState *common.ClusterState) bool {
	if currentState.KeycloakDeployment == nil {
		return false
	}
	deployedImage := currentState.KeycloakDeployment.Spec.Template.Spec.Containers[0].Image
	currentImage := model.KeycloakImage
	if cr.Spec.Profile == common.RHSSOProfile {
		currentImage = model.RHSSOImage
	}
	return deployedImage != currentImage
}

func findDeployment(desiredState common.DesiredClusterState) *v13.StatefulSet {
	for _, v := range desiredState {
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
