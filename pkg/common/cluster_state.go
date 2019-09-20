package common

import (
	"context"
	kc "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/model/keycloak"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// The desired cluster state is defined by a list of actions that have to be run to
// get from the current state to the desired state
type DesiredClusterState []ClusterAction

type ClusterState struct {
	KeycloakService *v1.Service
	client          client.Client
}

func NewClusterState(client client.Client) *ClusterState {
	return &ClusterState{
		KeycloakService: nil,
		client:          client,
	}
}

func (i *ClusterState) Read(cr *kc.Keycloak) {
	i.readKeycloakServiceCurrentState(cr)
}

// Keycloak service
func (i *ClusterState) readKeycloakServiceCurrentState(cr *kc.Keycloak) {
	keycloakService := keycloak.Service(cr)

	selector := client.ObjectKey{
		Name:      keycloakService.Name,
		Namespace: keycloakService.Namespace,
	}

	err := i.client.Get(context.TODO(), selector, keycloakService)
	if err != nil {
		i.KeycloakService = nil
	} else {
		i.KeycloakService = keycloakService.DeepCopy()
	}
}
