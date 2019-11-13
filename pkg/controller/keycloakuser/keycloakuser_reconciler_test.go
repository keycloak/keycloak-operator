package keycloakuser

import (
	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/common"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getMockState() *common.UserState {
	return common.NewUserState()
}

func getMockUser() *v1alpha1.KeycloakUser {
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
				RealmRoles:  nil,
				ClientRoles: nil,
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
