package common

import (
	"context"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type GroupState struct {
	Group                *v1alpha1.KeycloakAPIGroup
	ClientRoles          map[string][]*v1alpha1.KeycloakUserRole
	RealmRoles           []*v1alpha1.KeycloakUserRole
	AvailableClientRoles map[string][]*v1alpha1.KeycloakUserRole
	AvailableRealmRoles  []*v1alpha1.KeycloakUserRole
	Clients              []*v1alpha1.KeycloakAPIClient
	Keycloak             v1alpha1.Keycloak
	Context              context.Context
}

func NewGroupState(keycloak v1alpha1.Keycloak) *GroupState {
	return &GroupState{
		ClientRoles:          map[string][]*v1alpha1.KeycloakUserRole{},
		AvailableClientRoles: map[string][]*v1alpha1.KeycloakUserRole{},
		Keycloak:             keycloak,
	}
}

func (i *GroupState) Read(keycloakClient KeycloakInterface, groupClient client.Client, group *v1alpha1.KeycloakGroup, realm v1alpha1.KeycloakRealm) error {
	err := i.readGroup(keycloakClient, group, realm.Spec.Realm.Realm)
	if err != nil {
		// If there was an error reading the user then don't attempt
		// to read the roles. This user might not yet exist
		return nil
	}

	// Don't continue if the group could not be found
	if i.Group == nil {
		return nil
	}

	err = i.readRealmRoles(keycloakClient, group, realm.Spec.Realm.Realm)
	if err != nil {
		return err
	}

	return i.readClientRoles(keycloakClient, group, realm.Spec.Realm.Realm)
}

func (i *GroupState) readGroup(client KeycloakInterface, group *v1alpha1.KeycloakGroup, realm string) error {
	if group.Spec.Group.ID != "" {
		keycloakGroup, err := client.GetGroup(group.Spec.Group.ID, realm)
		if err != nil {
			return err
		}
		i.Group = keycloakGroup
	}
	return nil
}

func (i *GroupState) readRealmRoles(client KeycloakInterface, user *v1alpha1.KeycloakGroup, realm string) error {
	// Get all the realm roles of this user
	roles, err := client.ListGroupRealmRoles(realm, i.Group.ID)
	if err != nil {
		return err
	}
	i.RealmRoles = roles

	// Get the roles that are still available to this user
	availableRoles, err := client.ListAvailableGroupRealmRoles(realm, i.Group.ID)
	if err != nil {
		return err
	}
	i.AvailableRealmRoles = availableRoles

	return nil
}

func (i *GroupState) readClientRoles(client KeycloakInterface, group *v1alpha1.KeycloakGroup, realm string) error {
	clients, err := client.ListClients(realm)
	if err != nil {
		return err
	}
	i.Clients = clients

	for _, c := range clients {
		// Get all client roles of this user
		roles, err := client.ListGroupClientRoles(realm, c.ID, i.Group.ID)
		if err != nil {
			return err
		}
		i.ClientRoles[c.ClientID] = roles

		// Get the roles that are still available to this user
		availableRoles, err := client.ListAvailableGroupClientRoles(realm, c.ID, i.Group.ID)
		if err != nil {
			return err
		}
		i.AvailableClientRoles[c.ClientID] = availableRoles
	}
	return nil
}

// Check if a realm role is part of the available roles for this user
// Don't allow to assign unavailable roles
func (i *GroupState) GetAvailableRealmRole(name string) *v1alpha1.KeycloakUserRole {
	for _, role := range i.AvailableRealmRoles {
		if role.Name == name {
			return role
		}
	}
	return nil
}

// Check if a client role is part of the available roles for this user
// Don't allow to assign unavailable roles
func (i *GroupState) GetAvailableClientRole(name, clientID string) *v1alpha1.KeycloakUserRole {
	for _, role := range i.AvailableClientRoles[clientID] {
		if role.Name == name {
			return role
		}
	}
	return nil
}

// Keycloak clients have `ID` and `ClientID` properties and depending on the action we
// need one or the other. This function translates between the two
func (i *GroupState) GetClientByID(clientID string) *v1alpha1.KeycloakAPIClient {
	for _, client := range i.Clients {
		if client.ClientID == clientID {
			return client
		}
	}
	return nil
}
