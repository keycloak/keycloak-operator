package common

import (
	"testing"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	routev1 "github.com/openshift/api/route/v1"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func TestIsResourcesReady_Test_Route_Disabled_(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			ExternalAccess: v1alpha1.KeycloakExternalAccess{
				Enabled: false,
			},
		},
	}
	clusterState := prepareTestingClusterState()

	// when
	GetStateManager().SetState(RouteKind, true)
	ready, err := clusterState.IsResourcesReady(cr)

	// then
	assert.Nil(t, err)
	assert.Equal(t, true, ready)
}

func TestIsResourcesReady_Test_Route_Enabled_(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			ExternalAccess: v1alpha1.KeycloakExternalAccess{
				Enabled: true,
			},
		},
	}
	clusterState := prepareTestingClusterState()

	// when
	GetStateManager().SetState(RouteKind, true)
	ready, err := clusterState.IsResourcesReady(cr)

	// then
	assert.Nil(t, err)
	assert.Equal(t, false, ready)
}
func TestIsResourcesReady_Test_Route_Enabled_Exists_(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			ExternalAccess: v1alpha1.KeycloakExternalAccess{
				Enabled: true,
			},
		},
	}
	clusterState := prepareTestingClusterState()
	clusterState.KeycloakRoute = &routev1.Route{
		Status: routev1.RouteStatus{
			Ingress: []routev1.RouteIngress{
				{
					Conditions: []routev1.RouteIngressCondition{
						{
							Type:   routev1.RouteAdmitted,
							Status: corev1.ConditionTrue,
						},
					},
				},
			},
		},
	}

	// when
	GetStateManager().SetState(RouteKind, true)
	ready, err := clusterState.IsResourcesReady(cr)

	// then
	assert.Nil(t, err)
	assert.Equal(t, true, ready)
}

func prepareTestingClusterState() *ClusterState {
	clusterState := NewClusterState()

	var replicas int32 = 3
	clusterState.KeycloakDeployment = &appsv1.StatefulSet{
		Spec: appsv1.StatefulSetSpec{
			Replicas: &replicas,
		},
		Status: appsv1.StatefulSetStatus{
			Replicas:        replicas,
			ReadyReplicas:   replicas,
			CurrentRevision: "1",
			UpdateRevision:  "1",
		},
	}
	clusterState.PostgresqlDeployment = &appsv1.Deployment{
		Status: appsv1.DeploymentStatus{
			Conditions: []appsv1.DeploymentCondition{},
		},
	}
	return clusterState
}
