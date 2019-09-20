package keycloak

import (
	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/common"
	"github.com/keycloak/keycloak-operator/pkg/test"
	"testing"
)

var actionRunner = test.NewMockActionRunner()
var mockCr = v1alpha1.Keycloak{}

func TestKeycloakReconciler_Reconcile(t *testing.T) {
	currentState := &common.ClusterState{
		KeycloakService: nil,
	}

	reconciler := NewKeycloakReconciler(currentState, actionRunner)
	reconciler.Reconcile(&mockCr)

	runner := reconciler.runner.(*test.MockActionRunner)
	if runner.ResourcesCreated != 1 {
		t.Error("invalid number of resources created")
	}

	if runner.ResourcesUpdated != 0 {
		t.Error("invalid number of resources updated")
	}
}
