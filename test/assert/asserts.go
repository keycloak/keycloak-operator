package assert

import (
	"testing"

	"github.com/keycloak/keycloak-operator/pkg/common"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/apps/v1"
)

func ReplicasCount(t *testing.T, state common.DesiredClusterState, expectedCount int32) {
	assert.Equal(t, &[]int32{expectedCount}[0], state[0].(common.GenericUpdateAction).Ref.(*v1.StatefulSet).Spec.Replicas)
}

func ReplicasCountRecreate(t *testing.T, state common.DesiredClusterState, expectedCount int32) {
	assert.Equal(t, &[]int32{expectedCount}[0], state[1].(common.GenericCreateAction).Ref.(*v1.StatefulSet).Spec.Replicas)
}

func KcDeploymentRecreated(t *testing.T, state common.DesiredClusterState, kcDeployment *v1.StatefulSet, expectedRecreation bool) {
	found := false

	for i, v := range state {
		deleteAction, ok := v.(common.GenericDeleteAction)
		if !ok {
			continue
		}

		statefulSet, ok := deleteAction.Ref.(*v1.StatefulSet)
		if !ok || kcDeployment.Name != statefulSet.Name {
			continue
		}

		// if the previous action was deletion the next action must be creation
		createAction, ok := state[i+1].(common.GenericCreateAction)
		assert.True(t, ok)
		statefulSet, ok = createAction.Ref.(*v1.StatefulSet)
		assert.True(t, ok)
		assert.Equal(t, kcDeployment.Name, statefulSet.Name)

		found = true
		break
	}
	assert.Equal(t, expectedRecreation, found)
}

func KcDeploymentUpdated(t *testing.T, state common.DesiredClusterState, kcDeployment *v1.StatefulSet, expectedUpdate bool) {
	found := false

	for _, v := range state {
		updateAction, ok := v.(common.GenericUpdateAction)
		if !ok {
			continue
		}

		statefulSet, ok := updateAction.Ref.(*v1.StatefulSet)
		found = ok && kcDeployment.Name == statefulSet.Name
		if found {
			break
		}
	}
	assert.Equal(t, expectedUpdate, found)
}
