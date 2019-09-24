package keycloak

import (
	"testing"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/common"
	"github.com/keycloak/keycloak-operator/pkg/model/keycloak"
	"github.com/stretchr/testify/assert"
)

func TestKeycloakReconciler_Test_Creating_Example_Service(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{}
	currentState := &common.ClusterState{
		KeycloakService: nil,
	}

	// when
	reconciler := NewKeycloakReconciler()
	desiredState, error := reconciler.Reconcile(currentState, cr)

	// then
	assert.Nil(t, error)
	assert.IsType(t, common.GenericCreateAction{}, desiredState[0])
	assert.IsType(t, keycloak.Service(cr), desiredState[0].(common.GenericCreateAction).Ref)
}
