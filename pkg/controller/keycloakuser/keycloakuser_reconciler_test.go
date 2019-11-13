package keycloakuser

import (
	"testing"

	v12 "k8s.io/api/core/v1"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/common"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getDummyState(keycloak v1alpha1.Keycloak) *common.UserState {
	return common.NewUserState(keycloak)
}

func getDummyUser() *v1alpha1.KeycloakUser {
	return &v1alpha1.KeycloakUser{
		ObjectMeta: v1.ObjectMeta{
			Name:      "dummy",
			Namespace: "dummy",
		},
		Spec: v1alpha1.KeycloakUserSpec{
			RealmSelector: &v1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "sso",
				},
			},
			User: v1alpha1.KeycloakAPIUser{
				UserName:    "",
				ClientRoles: nil,
				RealmRoles:  []string{"dummy_role"},
				Credentials: []v1alpha1.KeycloakCredential{
					{
						Type:  "password",
						Value: "12345",
					},
				},
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
				Users: []*v1alpha1.KeycloakAPIUser{
					{
						ID:         "dummy",
						UserName:   "dummy",
						FirstName:  "dummy",
						LastName:   "dummy",
						Enabled:    true,
						RealmRoles: []string{"dummy_role"},
						Credentials: []v1alpha1.KeycloakCredential{
							{
								Type:      "password",
								Value:     "password",
								Temporary: false,
							},
						},
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
	reconciler := NewKeycloakuserReconciler(keycloak, realm)
	state := getDummyState(keycloak)
	user := getDummyUser()

	// when
	desiredState := reconciler.Reconcile(state, user)

	// then
	// 0 - check keycloak available
	// 1 - create user
	// 2 - create user secret
	assert.IsType(t, &common.PingAction{}, desiredState[0])
	assert.IsType(t, &common.CreateUserAction{}, desiredState[1])
	assert.IsType(t, &common.GenericCreateAction{}, desiredState[2])

	state.User = &user.Spec.User
	state.Secret = &v12.Secret{}
	state.AvailableRealmRoles = []*v1alpha1.KeycloakUserRole{
		{
			ID:         "dummy_role",
			Name:       "dummy_role",
			ClientRole: false,
		},
	}

	desiredState = reconciler.Reconcile(state, user)
	assert.IsType(t, &common.PingAction{}, desiredState[0])
	assert.IsType(t, &common.UpdateUserAction{}, desiredState[1])
	assert.IsType(t, &common.AssignRealmRoleAction{}, desiredState[2])
}
