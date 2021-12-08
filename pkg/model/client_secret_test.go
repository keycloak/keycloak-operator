package model

import (
	"testing"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/stretchr/testify/assert"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestClientSecret_testNoSecretTemplate(t *testing.T) {
	// given
	cr := &v1alpha1.KeycloakClient{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      "test",
			Namespace: "testns",
		},
		Spec: v1alpha1.KeycloakClientSpec{
			RealmSelector: &meta_v1.LabelSelector{
				MatchLabels: map[string]string{"application": "sso"},
			},
			Client: &v1alpha1.KeycloakAPIClient{
				ClientID: "client_id",
				Secret:   "client_secret",
			},
		},
	}

	//when
	secret := ClientSecret(cr)

	//then
	assert.Equal(t, "keycloak-client-secret-client-id", secret.ObjectMeta.Name)
	assert.Equal(t, "testns", secret.ObjectMeta.Namespace)
	assert.Equal(t, 1, len(secret.ObjectMeta.Labels))
	assert.Equal(t, 0, len(secret.ObjectMeta.Annotations))
	assert.Equal(t, []byte("client_id"), secret.Data[ClientSecretClientIDProperty])
	assert.Equal(t, []byte("client_secret"), secret.Data[ClientSecretClientSecretProperty])
}

func TestClientSecret_testSecretTemplateLabels(t *testing.T) {
	// given
	cr := &v1alpha1.KeycloakClient{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      "test",
			Namespace: "testns",
		},
		Spec: v1alpha1.KeycloakClientSpec{
			RealmSelector: &meta_v1.LabelSelector{
				MatchLabels: map[string]string{"application": "sso"},
			},
			Client: &v1alpha1.KeycloakAPIClient{
				ClientID: "client_id",
				Secret:   "client_secret",
				SecretTemplate: &v1alpha1.SecretTemplate{
					Metadata: &v1alpha1.SecretTemplateMetadata{
						Labels: map[string]string{
							"foo":  "bar",
							"toto": "titi",
						},
					},
				},
			},
		},
	}

	//when
	secret := ClientSecret(cr)

	//then
	expectedLabels := map[string]string{
		"app":  "keycloak",
		"foo":  "bar",
		"toto": "titi",
	}
	expectedAnnotations := map[string]string(nil)
	assert.Equal(t, "keycloak-client-secret-client-id", secret.ObjectMeta.Name)
	assert.Equal(t, "testns", secret.ObjectMeta.Namespace)
	assert.Equal(t, expectedLabels, secret.GetLabels())
	assert.Equal(t, expectedAnnotations, secret.GetAnnotations())
	assert.Equal(t, 0, len(secret.ObjectMeta.Annotations))
	assert.Equal(t, []byte("client_id"), secret.Data[ClientSecretClientIDProperty])
	assert.Equal(t, []byte("client_secret"), secret.Data[ClientSecretClientSecretProperty])
}
