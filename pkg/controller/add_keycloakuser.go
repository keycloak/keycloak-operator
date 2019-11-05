package controller

import (
	"github.com/keycloak/keycloak-operator/pkg/controller/keycloakuser"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, keycloakuser.Add)
}
