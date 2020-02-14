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
	Client       *kc.KeycloakAPIClient
	ClientSecret *v1.Secret
	Context      context.Context
	Realm        *kc.KeycloakRealm
}

func NewClientState(context context.Context, realm *kc.KeycloakRealm) *ClientState {
	return &ClientState{
		Context: context,
		Realm:   realm,
	}
}

func (i *ClientState) Read(context context.Context, cr *kc.KeycloakClient, realmClient KeycloakInterface, controllerClient client.Client) error {
	client, err := realmClient.GetClient(cr.Spec.Client.ID, i.Realm.Spec.Realm.Realm)
	i.Client = client
	if err != nil {
		return err
	}

	if client == nil {
		return nil
	}

	clientSecret, err := realmClient.GetClientSecret(cr.Spec.Client.ID, i.Realm.Spec.Realm.Realm)
	if err != nil {
		return err
	}
	cr.Spec.Client.Secret = clientSecret

	err = i.readClientSecret(context, cr, i.Client, controllerClient)
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
