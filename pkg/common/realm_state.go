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
	Realm             *kc.KeycloakRealm
	RealmUserSecrets  map[string]*v1.Secret
	Context           context.Context
	Keycloak          *kc.Keycloak
	KeycloakCLIClient *kc.KeycloakClient
}

func NewRealmState(context context.Context, keycloak kc.Keycloak) *RealmState {
	return &RealmState{
		Context:  context,
		Keycloak: &keycloak,
	}
}

func (i *RealmState) ReadRealmCurrentState(cr *kc.KeycloakRealm, realmClient KeycloakInterface) (*kc.KeycloakRealm, error) {
	realm, err := realmClient.GetRealm(cr.Spec.Realm.Realm)
	if err != nil {
		return nil, err
	}
	return realm, nil
}

func (i *RealmState) Read(cr *kc.KeycloakRealm, realmClient KeycloakInterface, controllerClient client.Client) error {
	var err error
	i.Realm, err = i.ReadRealmCurrentState(cr, realmClient)
	if err != nil {
		return err
	}

	if i.Realm == nil {
		return nil
	}

	err = i.readKeycloakOpenShiftCLIClientCurrentState(cr, controllerClient)
	if err != nil {
		return err
	}

	if len(cr.Spec.Realm.Users) == 0 {
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

func (i *RealmState) readKeycloakOpenShiftCLIClientCurrentState(cr *kc.KeycloakRealm, controllerClient client.Client) error {
	keycloakCLIClient := model.KeycloakOperatorCLIClient(cr)
	keycloakCLIClientSelector := model.KeycloakOperatorCLIClientSelector(cr)

	err := controllerClient.Get(i.Context, keycloakCLIClientSelector, keycloakCLIClient)
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	} else {
		i.KeycloakCLIClient = keycloakCLIClient.DeepCopy()
		cr.UpdateStatusSecondaryResources(i.KeycloakCLIClient.Kind, i.KeycloakCLIClient.Name)
	}
	return nil
}
