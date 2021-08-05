package keycloakgroup

import (
	"testing"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/common"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getDummyState(keycloak v1alpha1.Keycloak) *common.GroupState {
	return common.NewGroupState(keycloak)
}

func getDummyGroup() *v1alpha1.KeycloakGroup {
	return &v1alpha1.KeycloakGroup{
		ObjectMeta: v1.ObjectMeta{
			Name:      "dummy",
			Namespace: "dummy",
		},
		Spec: v1alpha1.KeycloakGroupSpec{
			RealmSelector: &v1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "sso",
				},
			},
			Group: v1alpha1.KeycloakAPIGroup{
				Name:        "",
				ClientRoles: nil,
				RealmRoles:  []string{"dummy_role"},
			},
		},
	}
}

func getDummyRealm() v1alpha1.KeycloakRealm {
	return v1alpha1.KeycloakRealm{
		Spec: v1alpha1.KeycloakRealmSpec{
			InstanceSelector: &v1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "keycloak",
				},
			},
			RealmOverrides: []*v1alpha1.RedirectorIdentityProviderOverride{
				{
					IdentityProvider: "openshift-v4",
					ForFlow:          "browser",
				},
			},
			Realm: &v1alpha1.KeycloakAPIRealm{
				ID:          "dummy",
				Realm:       "dummy",
				Enabled:     true,
				DisplayName: "dummy",
				Groups: []*v1alpha1.KeycloakAPIGroup{
					{
						Name: "dummy",
					},
				},
			},
		},
	}
}

func TestKeycloakRealmReconciler_Reconcile(t *testing.T) {
	// given
	keycloak := v1alpha1.Keycloak{}
	realm := getDummyRealm()
	reconciler := NewKeycloakgroupReconciler(keycloak, realm)
	state := getDummyState(keycloak)
	group := getDummyGroup()

	// when
	desiredState := reconciler.Reconcile(state, group)

	// then
	// 0 - check keycloak available
	// 1 - create group
	assert.IsType(t, &common.PingAction{}, desiredState[0])
	assert.IsType(t, &common.CreateGroupAction{}, desiredState[1])

	state.Group = &group.Spec.Group
	state.AvailableRealmRoles = []*v1alpha1.KeycloakUserRole{
		{
			ID:         "dummy_role",
			Name:       "dummy_role",
			ClientRole: false,
		},
	}

	desiredState = reconciler.Reconcile(state, group)
	assert.IsType(t, &common.PingAction{}, desiredState[0])
	assert.IsType(t, &common.UpdateGroupAction{}, desiredState[1])
	assert.IsType(t, &common.AssignRealmRoleAction{}, desiredState[2])
}
