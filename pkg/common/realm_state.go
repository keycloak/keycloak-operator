package common

import (
	"context"

	kc "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/model"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type RealmState struct {
	Realm            *kc.KeycloakRealm
	RealmUserSecrets map[string]*v1.Secret
	Context          context.Context
	Keycloak         *kc.Keycloak
}

func NewRealmState(context context.Context, keycloak kc.Keycloak) *RealmState {
	return &RealmState{
		Context:  context,
		Keycloak: &keycloak,
	}
}

func (i *RealmState) Read(cr *kc.KeycloakRealm, realmClient KeycloakInterface, controllerClient client.Client) error {
	realm, err := realmClient.GetRealm(cr.Spec.Realm.Realm)
	if err != nil {
		i.Realm = nil
		return err
	}

	i.Realm = realm
	if realm == nil || len(cr.Spec.Realm.Users) == 0 {
		return nil
	}

	// Get the state of the realm users
	i.RealmUserSecrets = make(map[string]*v1.Secret)
	for _, user := range cr.Spec.Realm.Users {
		secret, err := i.readRealmUserSecret(cr, user, controllerClient)
		if err != nil {
			return err
		}
		i.RealmUserSecrets[user.UserName] = secret

		cr.UpdateStatusSecondaryResources(SecretKind, model.GetRealmUserSecretName(i.Keycloak.Namespace, cr.Spec.Realm.Realm, user.UserName))
	}

	return nil
}

func (i *RealmState) readRealmUserSecret(realm *kc.KeycloakRealm, user *kc.KeycloakAPIUser, controllerClient client.Client) (*v1.Secret, error) {
	key := model.RealmCredentialSecretSelector(realm, user, i.Keycloak)
	secret := &v1.Secret{}

	// Try to find the user credential secret
	err := controllerClient.Get(i.Context, key, secret)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	return secret, err
}
