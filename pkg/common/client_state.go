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
	Client               *kc.KeycloakAPIClient
	ClientSecret         *v1.Secret
	Context              context.Context
	Realm                *kc.KeycloakRealm
	ClientRoles          map[string][]*kc.KeycloakUserRole
	RealmRoles           []*kc.KeycloakUserRole
	AvailableClientRoles map[string][]*kc.KeycloakUserRole
	AvailableRealmRoles  []*kc.KeycloakUserRole
	Clients              []*kc.KeycloakAPIClient
}

func NewClientState(context context.Context, realm *kc.KeycloakRealm) *ClientState {
	return &ClientState{
		Context:              context,
		Realm:                realm,
		ClientRoles:          map[string][]*kc.KeycloakUserRole{},
		AvailableClientRoles: map[string][]*kc.KeycloakUserRole{},
	}
}

func (i *ClientState) Read(context context.Context, cr *kc.KeycloakClient, realmClient KeycloakInterface, controllerClient client.Client) error {
	if cr.Status.RealmClientIds[i.Realm.Spec.Realm.Realm] == "" {
		return nil
	}

	client, err := realmClient.GetClient(cr.Status.RealmClientIds[i.Realm.Spec.Realm.Realm], i.Realm.Spec.Realm.Realm)

	if err != nil {
		return err
	}

	i.Client = client

	if cr.Spec.Client.ServiceAccountsEnabled && cr.Status.RealmServiceAccountIds[i.Realm.Spec.Realm.Realm] != "" {
		err = i.readRealmRoles(realmClient, cr.Status.RealmServiceAccountIds[i.Realm.Spec.Realm.Realm], i.Realm.Spec.Realm.Realm)
		if err != nil {
			return err
		}

		err = i.readClientRoles(realmClient, cr.Status.RealmServiceAccountIds[i.Realm.Spec.Realm.Realm], i.Realm.Spec.Realm.Realm)
		if err != nil {
			return err
		}
	}

	clientSecret, err := realmClient.GetClientSecret(cr.Status.RealmClientIds[i.Realm.Spec.Realm.Realm], i.Realm.Spec.Realm.Realm)
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

func (i *ClientState) readRealmRoles(client KeycloakInterface, userID, realm string) error {
	// Get all the realm roles of this user
	roles, err := client.ListUserRealmRoles(realm, userID)
	if err != nil {
		return err
	}
	i.RealmRoles = roles

	// Get the roles that are still available to this user
	availableRoles, err := client.ListAvailableUserRealmRoles(realm, userID)
	if err != nil {
		return err
	}
	i.AvailableRealmRoles = availableRoles

	return nil
}

func (i *ClientState) readClientRoles(client KeycloakInterface, userID, realm string) error {
	clients, err := client.ListClients(realm)
	if err != nil {
		return err
	}
	i.Clients = clients

	for _, c := range clients {
		// Get all client roles of this user
		roles, err := client.ListUserClientRoles(realm, c.ID, userID)
		if err != nil {
			return err
		}
		i.ClientRoles[c.ClientID] = roles

		// Get the roles that are still available to this user
		availableRoles, err := client.ListAvailableUserClientRoles(realm, c.ID, userID)
		if err != nil {
			return err
		}
		i.AvailableClientRoles[c.ClientID] = availableRoles
	}
	return nil
}

// Check if a realm role is part of the available roles for this user
// Don't allow to assign unavailable roles
func (i *ClientState) GetAvailableRealmRole(name string) *kc.KeycloakUserRole {
	for _, role := range i.AvailableRealmRoles {
		if role.Name == name {
			return role
		}
	}
	return nil
}

// Check if a client role is part of the available roles for this user
// Don't allow to assign unavailable roles
func (i *ClientState) GetAvailableClientRole(name, clientID string) *kc.KeycloakUserRole {
	for _, role := range i.AvailableClientRoles[clientID] {
		if role.Name == name {
			return role
		}
	}
	return nil
}

// Keycloak clients have `ID` and `ClientID` properties and depending on the action we
// need one or the other. This function translates between the two
func (i *ClientState) GetClientByID(clientID string) *kc.KeycloakAPIClient {
	for _, client := range i.Clients {
		if client.ClientID == clientID {
			return client
		}
	}
	return nil
}
