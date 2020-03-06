package common

import (
	"context"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/model"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type UserState struct {
	User                 *v1alpha1.KeycloakAPIUser
	ClientRoles          map[string][]*v1alpha1.KeycloakUserRole
	RealmRoles           []*v1alpha1.KeycloakUserRole
	AvailableClientRoles map[string][]*v1alpha1.KeycloakUserRole
	AvailableRealmRoles  []*v1alpha1.KeycloakUserRole
	Clients              []*v1alpha1.KeycloakAPIClient
	Secret               *v1.Secret
	Keycloak             v1alpha1.Keycloak
	Context              context.Context
}

func NewUserState(keycloak v1alpha1.Keycloak) *UserState {
	return &UserState{
		ClientRoles:          map[string][]*v1alpha1.KeycloakUserRole{},
		AvailableClientRoles: map[string][]*v1alpha1.KeycloakUserRole{},
		Keycloak:             keycloak,
	}
}

func (i *UserState) Read(keycloakClient KeycloakInterface, userClient client.Client, user *v1alpha1.KeycloakUser, realm v1alpha1.KeycloakRealm) error {
	err := i.readUser(keycloakClient, user, realm.Spec.Realm.Realm)
	if err != nil {
		// If there was an error reading the user then don't attempt
		// to read the roles. This user might not yet exist
		return nil
	}

	// Don't continue if the user could not be found
	if i.User == nil {
		return nil
	}

	err = i.readRealmRoles(keycloakClient, user, realm.Spec.Realm.Realm)
	if err != nil {
		return err
	}

	err = i.readClientRoles(keycloakClient, user, realm.Spec.Realm.Realm)
	if err != nil {
		return err
	}

	return i.readSecretState(userClient, user, &realm)
}

func (i *UserState) readUser(client KeycloakInterface, user *v1alpha1.KeycloakUser, realm string) error {
	if user.Spec.User.ID != "" {
		keycloakUser, err := client.GetUser(user.Spec.User.ID, realm)
		if err != nil {
			return err
		}
		i.User = keycloakUser
	}
	return nil
}

func (i *UserState) readRealmRoles(client KeycloakInterface, user *v1alpha1.KeycloakUser, realm string) error {
	// Get all the realm roles of this user
	roles, err := client.ListUserRealmRoles(realm, i.User.ID)
	if err != nil {
		return err
	}
	i.RealmRoles = roles

	// Get the roles that are still available to this user
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
		// Get all client roles of this user
		roles, err := client.ListUserClientRoles(realm, c.ID, i.User.ID)
		if err != nil {
			return err
		}
		i.ClientRoles[c.ClientID] = roles

		// Get the roles that are still available to this user
		availableRoles, err := client.ListAvailableUserClientRoles(realm, c.ID, i.User.ID)
		if err != nil {
			return err
		}
		i.AvailableClientRoles[c.ClientID] = availableRoles
	}
	return nil
}

func (i *UserState) readSecretState(userClient client.Client, user *v1alpha1.KeycloakUser, realm *v1alpha1.KeycloakRealm) error {
	key := model.RealmCredentialSecretSelector(realm, &user.Spec.User, &i.Keycloak)
	secret := &v1.Secret{}

	// Try to find the user credential secret
	err := userClient.Get(i.Context, key, secret)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	i.Secret = secret
	return nil
}

// Check if a realm role is part of the available roles for this user
// Don't allow to assign unavailable roles
func (i *UserState) GetAvailableRealmRole(name string) *v1alpha1.KeycloakUserRole {
	for _, role := range i.AvailableRealmRoles {
		if role.Name == name {
			return role
		}
	}
	return nil
}

// Check if a client role is part of the available roles for this user
// Don't allow to assign unavailable roles
func (i *UserState) GetAvailableClientRole(name, clientID string) *v1alpha1.KeycloakUserRole {
	for _, role := range i.AvailableClientRoles[clientID] {
		if role.Name == name {
			return role
		}
	}
	return nil
}

// Keycloak clients have `ID` and `ClientID` properties and depending on the action we
// need one or the other. This function translates between the two
func (i *UserState) GetClientByID(clientID string) *v1alpha1.KeycloakAPIClient {
	for _, client := range i.Clients {
		if client.ClientID == clientID {
			return client
		}
	}
	return nil
}
