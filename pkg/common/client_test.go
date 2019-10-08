package common

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/stretchr/testify/assert"
)

const (
	RealmsCreatePath = "/auth/admin/realms"
	RealmsDeletePath = "/auth/admin/realms/%s"
	UserCreatePath   = "/auth/admin/realms/%s/users"
	UserDeletePath   = "/auth/admin/realms/%s/users/%s"
)

func getDummyRealm() *v1alpha1.KeycloakRealm {
	return &v1alpha1.KeycloakRealm{
		Spec: v1alpha1.KeycloakRealmSpec{
			KeycloakApiRealm: &v1alpha1.KeycloakApiRealm{
				ID:          "dummy",
				Realm:       "dummy",
				Enabled:     false,
				DisplayName: "dummy",
			},
		},
	}
}

func getDummyUser() *v1alpha1.KeycloakApiUser {
	return &v1alpha1.KeycloakApiUser{
		ID:            "dummy",
		UserName:      "dummy",
		FirstName:     "dummy",
		LastName:      "dummy",
		EmailVerified: false,
		Enabled:       false,
	}
}

func TestClient_CreateRealm(t *testing.T) {
	// given
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, RealmsCreatePath, req.URL.Path)
		w.WriteHeader(201)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	client := Client{
		requester: server.Client(),
		URL:       server.URL,
		token:     "dummy",
	}

	realm := getDummyRealm()

	// when
	err := client.CreateRealm(realm)

	// then
	// no error expected
	// correct path expected on httptest server
	assert.NoError(t, err)
}

func TestClient_DeleteRealmRealm(t *testing.T) {
	// given
	realm := getDummyRealm()

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, fmt.Sprintf(RealmsDeletePath, realm.Spec.Realm), req.URL.Path)
		w.WriteHeader(204)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	client := Client{
		requester: server.Client(),
		URL:       server.URL,
		token:     "dummy",
	}

	// when
	err := client.DeleteRealm(realm.Spec.Realm)

	// then
	// correct path expected on httptest server
	assert.NoError(t, err)
}

func TestClient_CreateUser(t *testing.T) {
	// given
	user := getDummyUser()
	realm := getDummyRealm()

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, fmt.Sprintf(UserCreatePath, realm.Spec.Realm), req.URL.Path)
		w.WriteHeader(201)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	client := Client{
		requester: server.Client(),
		URL:       server.URL,
		token:     "dummy",
	}

	// when
	err := client.CreateUser(user, realm.Spec.Realm)

	// then
	// correct path expected on httptest server
	assert.NoError(t, err)
}

func TestClient_DeleteUser(t *testing.T) {
	// given
	user := getDummyUser()
	realm := getDummyRealm()

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, fmt.Sprintf(UserDeletePath, realm.Spec.Realm, user.ID), req.URL.Path)
		w.WriteHeader(204)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	client := Client{
		requester: server.Client(),
		URL:       server.URL,
		token:     "dummy",
	}

	// when
	err := client.DeleteUser(user.ID, realm.Spec.Realm)

	// then
	// correct path expected on httptest server
	assert.NoError(t, err)
}
