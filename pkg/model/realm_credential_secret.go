package model

import (
	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func RealmCredentialSecret(cr *v1alpha1.KeycloakRealm, user *v1alpha1.KeycloakAPIUser, keycloak *v1alpha1.Keycloak) *v1.Secret {
	outputSecretName := GetRealmUserSecretName(keycloak.Namespace, cr.Spec.Realm.Realm, user.UserName)

	outputSecret := &v1.Secret{}
	outputSecret.ObjectMeta = v12.ObjectMeta{
		Namespace: cr.Namespace,
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

func RealmCredentialSecretSelector(cr *v1alpha1.KeycloakRealm, user *v1alpha1.KeycloakAPIUser, keycloak *v1alpha1.Keycloak) client.ObjectKey {
	outputSecretName := GetRealmUserSecretName(keycloak.Namespace, cr.Spec.Realm.Realm, user.UserName)

	return client.ObjectKey{
		Name:      outputSecretName,
		Namespace: cr.Namespace,
	}
}
