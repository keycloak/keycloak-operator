package e2e

import "time"

const (
	testKeycloakCRName      = "keycloak-test"
	operatorCRName          = "keycloak-operator"
	testKeycloakRealmCRName = "keycloak-realm-test"
	cleanupRetryInterval    = time.Second * 5
	cleanupTimeout          = time.Minute * 2
	pollRetryInterval       = time.Second * 10
	pollTimeout             = time.Minute * 9
)
