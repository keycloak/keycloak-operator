package common

import (
	"context"

	kc "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/model"
	v1 "k8s.io/api/core/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ClientState struct {
	Client                *kc.KeycloakAPIClient
	ClientSecret          *v1.Secret
	Context               context.Context
	Realm                 *kc.KeycloakRealm
	Roles                 []kc.RoleRepresentation
	ScopeMappings         *kc.MappingsRepresentation
	AvailableClientScopes []kc.KeycloakClientScope
	DefaultClientScopes   []kc.KeycloakClientScope
	OptionalClientScopes  []kc.KeycloakClientScope
}

func NewClientState(context context.Context, realm *kc.KeycloakRealm) *ClientState {
	return &ClientState{
		Context: context,
		Realm:   realm,
	}
}

func (i *ClientState) Read(context context.Context, cr *kc.KeycloakClient, realmClient KeycloakInterface, controllerClient client.Client) error {
	if cr.Spec.Client.ID == "" {
		return nil
	}

	client, err := realmClient.GetClient(cr.Spec.Client.ID, i.Realm.Spec.Realm.Realm)

	if err != nil {
		return err
	}

	i.Client = client

	// CR could have updated with new secret, so set saved secret to Spec only when empty
	// Otherwise let reconcile loop to update secret with desired secret in CR
	if cr.Spec.Client.Secret == "" {
		clientSecret, err := realmClient.GetClientSecret(cr.Spec.Client.ID, i.Realm.Spec.Realm.Realm)
		if err != nil {
			return err
		}

		cr.Spec.Client.Secret = clientSecret
	}

	err = i.readClientSecret(context, cr, i.Client, controllerClient)
	if err != nil {
		return err
	}

	if i.Client != nil {
		i.Roles, err = realmClient.ListClientRoles(cr.Spec.Client.ID, i.Realm.Spec.Realm.Realm)
		if err != nil {
			return err
		}

		i.ScopeMappings, err = realmClient.ListScopeMappings(cr.Spec.Client.ID, i.Realm.Spec.Realm.Realm)
		if err != nil {
			return err
		}

		err := i.readClientScopes(cr, realmClient)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *ClientState) readClientScopes(cr *kc.KeycloakClient, realmClient KeycloakInterface) (err error) {
	// It is not strictly a property of the client but rather of the realm.
	// However could not figure out a better way to convey it to populate default and optional
	// client scopes which requires client scope IDs.
	i.AvailableClientScopes, err = realmClient.ListAvailableClientScopes(i.Realm.Spec.Realm.Realm)
	if err != nil {
		return err
	}

	i.DefaultClientScopes, err = realmClient.ListDefaultClientScopes(cr.Spec.Client.ID, i.Realm.Spec.Realm.Realm)
	if err != nil {
		return err
	}

	i.OptionalClientScopes, err = realmClient.ListOptionalClientScopes(cr.Spec.Client.ID, i.Realm.Spec.Realm.Realm)
	if err != nil {
		return err
	}
	return nil
}

func (i *ClientState) readClientSecret(context context.Context, cr *kc.KeycloakClient, clientSpec *kc.KeycloakAPIClient, controllerClient client.Client) error {
	key := model.ClientSecretSelector(cr)
	secret := model.ClientSecret(cr)

	err := controllerClient.Get(context, key, secret)
	if err != nil {
		if !apiErrors.IsNotFound(err) {
			return err
		}
	} else {
		i.ClientSecret = secret.DeepCopy()
		cr.UpdateStatusSecondaryResources(i.ClientSecret.Kind, i.ClientSecret.Name)
	}
	return nil
}
