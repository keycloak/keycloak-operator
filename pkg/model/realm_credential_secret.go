package model

import (
	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func RealmCredentialSecret(cr v1alpha1.KeycloakRealmReference, user *v1alpha1.KeycloakAPIUser, keycloak v1alpha1.KeycloakReference) *v1.Secret {
	outputSecretName := GetRealmUserSecretName(keycloak.GetNamespace(), cr.Realm(), user.UserName)

	outputSecret := &v1.Secret{}
	outputSecret.ObjectMeta = v12.ObjectMeta{
		Namespace: cr.GetNamespace(),
		Name:      outputSecretName,
	}
	outputSecret.Data = map[string][]byte{
		"username": []byte(user.UserName),
	}

	for _, credential := range user.Credentials {
		outputSecret.Data[credential.Type] = []byte(credential.Value)
	}

	return outputSecret
}

func RealmCredentialSecretSelector(cr v1alpha1.KeycloakRealmReference, user *v1alpha1.KeycloakAPIUser, keycloak v1alpha1.KeycloakReference) client.ObjectKey {
	outputSecretName := GetRealmUserSecretName(keycloak.GetNamespace(), cr.Realm(), user.UserName)

	return client.ObjectKey{
		Name:      outputSecretName,
		Namespace: cr.GetNamespace(),
	}
}
