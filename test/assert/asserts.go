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
