package common

import (
	"context"

	kc "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/model/keycloak"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// The desired cluster state is defined by a list of actions that have to be run to
// get from the current state to the desired state
type DesiredClusterState []ClusterAction

type ClusterState struct {
	KeycloakService *v1.Service
}

func NewClusterState() *ClusterState {
	return &ClusterState{
		KeycloakService: nil,
	}
}

func (i *ClusterState) Read(cr *kc.Keycloak, controllerClient client.Client) error {
	return i.readKeycloakServiceCurrentState(cr, controllerClient)
}

// Keycloak service
func (i *ClusterState) readKeycloakServiceCurrentState(cr *kc.Keycloak, controllerClient client.Client) error {
	keycloakService := keycloak.Service(cr)

	selector := client.ObjectKey{
		Name:      keycloakService.Name,
		Namespace: keycloakService.Namespace,
	}
	err := controllerClient.Get(context.TODO(), selector, keycloakService)

	if err != nil {
		if errors.IsNotFound(err) {
			i.KeycloakService = nil
		} else {
			return err
		}
	} else {
		i.KeycloakService = keycloakService.DeepCopy()
	}
	return nil
}
