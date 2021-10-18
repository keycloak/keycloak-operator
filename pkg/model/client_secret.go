package model

import (
	"unicode"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

func MakeK8sCompatibleName(text string) string {
	// we only want letters and numbers
	reg := []rune(SanitizeResourceName(text))

	// we guarantee first and last char is an alpha and not a dot or a dash
	if !unicode.IsLetter(reg[0]) && !unicode.IsNumber(reg[0]) {
		reg = append([]rune{'a'}, reg...)
	}
	if !unicode.IsLetter(reg[len(reg)-1]) && !unicode.IsNumber(reg[len(reg)-1]) {
		reg = append(reg, 'a')
	}

	return string(reg)
}

func ClientSecret(cr *v1alpha1.KeycloakClient) *v1.Secret {
	escapedClientIDName := MakeK8sCompatibleName(cr.Spec.Client.ClientID)
	return &v1.Secret{
		ObjectMeta: v12.ObjectMeta{
			Name:      ClientSecretName + "-" + escapedClientIDName,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app": ApplicationName,
			},
		},
		Data: map[string][]byte{
			ClientSecretClientIDProperty:     []byte(cr.Spec.Client.ClientID),
			ClientSecretClientSecretProperty: []byte(cr.Spec.Client.Secret),
		},
	}
}

func ClientSecretSelector(cr *v1alpha1.KeycloakClient) client.ObjectKey {
	escapedClientIDName := SanitizeResourceName(cr.Spec.Client.ClientID)
	return client.ObjectKey{
		Name:      ClientSecretName + "-" + escapedClientIDName,
		Namespace: cr.Namespace,
	}
}

func ClientSecretReconciled(cr *v1alpha1.KeycloakClient, currentState *v1.Secret) *v1.Secret {
	reconciled := currentState.DeepCopy()
	// Since the client is synced upon update, we always override what's there...
	reconciled.Data = map[string][]byte{
		ClientSecretClientIDProperty:     []byte(cr.Spec.Client.ClientID),
		ClientSecretClientSecretProperty: []byte(cr.Spec.Client.Secret),
	}
	return reconciled
}
