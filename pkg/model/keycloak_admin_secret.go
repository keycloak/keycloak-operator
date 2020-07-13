package model

import (
	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func KeycloakAdminSecret(cr *v1alpha1.Keycloak) *v1.Secret {
	return &v1.Secret{
		ObjectMeta: v12.ObjectMeta{
			Name:      "credential-" + cr.Name,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app":           ApplicationName,
				ApplicationName: cr.Name,
			},
		},
		Data: map[string][]byte{
			AdminUsernameProperty: []byte("admin"),
			AdminPasswordProperty: []byte(GenerateRandomString(10)),
		},
		Type: "Opaque",
	}
}

func KeycloakAdminSecretSelector(cr *v1alpha1.Keycloak) client.ObjectKey {
	return client.ObjectKey{
		Name:      "credential-" + cr.Name,
		Namespace: cr.Namespace,
	}
}

func KeycloakAdminSecretReconciled(cr *v1alpha1.Keycloak, currentState *v1.Secret) *v1.Secret {
	reconciled := currentState.DeepCopy()
	if val, ok := reconciled.Data[AdminUsernameProperty]; !ok || len(val) == 0 {
		reconciled.Data[AdminUsernameProperty] = []byte("admin")
	}
	if val, ok := reconciled.Data[AdminPasswordProperty]; !ok || len(val) == 0 {
		reconciled.Data[AdminPasswordProperty] = []byte(GenerateRandomString(10))
	}
	return reconciled
}
