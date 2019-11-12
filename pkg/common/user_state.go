package common

import (
	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
)

type UserState struct {
	User                 *v1alpha1.KeycloakAPIUser
	ClientRoles          map[string][]*v1alpha1.KeycloakUserRole
	RealmRoles           []*v1alpha1.KeycloakUserRole
	AvailableClientRoles map[string][]*v1alpha1.KeycloakUserRole
	AvailableRealmRoles  []*v1alpha1.KeycloakUserRole
	Clients              []*v1alpha1.KeycloakAPIClient
}

func NewUserState() *UserState {
	return &UserState{
		ClientRoles:          map[string][]*v1alpha1.KeycloakUserRole{},
		AvailableClientRoles: map[string][]*v1alpha1.KeycloakUserRole{},
	}
}

func (i *UserState) Read(client KeycloakInterface, user *v1alpha1.KeycloakUser, realm string) error {
	err := i.readUser(client, user, realm)
	if err != nil {
		// If there was an error reading the user then don't attempt
		// to read the roles. This user might not yet exist
		return nil
	}

	err = i.readRealmRoles(client, user, realm)
	if err != nil {
		return err
	}

	return i.readClientRoles(client, user, realm)
}

func (i *UserState) readUser(client KeycloakInterface, user *v1alpha1.KeycloakUser, realm string) error {
	keycloakUser, err := client.FindUserByUsername(user.Spec.User.UserName, realm)
	if err != nil {
		return err
	}

	i.User = keycloakUser
	return nil
}

func (i *UserState) readRealmRoles(client KeycloakInterface, user *v1alpha1.KeycloakUser, realm string) error {
	roles, err := client.ListUserRealmRoles(realm, i.User.ID)
	if err != nil {
		return err
	}
	i.RealmRoles = roles

	availableRoles, err := client.ListAvailableUserRealmRoles(realm, i.User.ID)
	if err != nil {
		return err
	}
	i.AvailableRealmRoles = availableRoles

	return nil
}

func (i *UserState) readClientRoles(client KeycloakInterface, user *v1alpha1.KeycloakUser, realm string) error {
	clients, err := client.ListClients(realm)
	if err != nil {
		return err
	}
	i.Clients = clients

	for _, c := range clients {
		roles, err := client.ListUserClientRoles(realm, c.ID, i.User.ID)
		if err != nil {
			return err
		}
		i.ClientRoles[c.ClientID] = roles

		availableRoles, err := client.ListAvailableUserClientRoles(realm, c.ID, i.User.ID)
		if err != nil {
			return err
		}
		i.AvailableClientRoles[c.ClientID] = availableRoles
	}
	return nil
}

func (i *UserState) GetAvailableRealmRole(name string) *v1alpha1.KeycloakUserRole {
	for _, role := range i.AvailableRealmRoles {
		if role.Name == name {
			return role
		}
	}
	return nil
}

func (i *UserState) GetAvailableClientRole(name, clientId string) *v1alpha1.KeycloakUserRole {
	for _, role := range i.AvailableClientRoles[clientId] {
		if role.Name == name {
			return role
		}
	}
	return nil
}

func (i *UserState) GetClientByID(clientID string) *v1alpha1.KeycloakAPIClient {
	for _, client := range i.Clients {
		if client.ClientID == clientID {
			return client
		}
	}
	return nil
}
