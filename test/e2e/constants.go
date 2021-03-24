package e2e

import "time"

const (
	testKeycloakCRName               = "keycloak-test"
	testKeycloakUnmanagedCRName      = "unmanaged-keycloak-test"
	testKeycloakExternalCRName       = "external-keycloak-test"
	operatorCRName                   = "keycloak-operator"
	testKeycloakRealmCRName          = "keycloak-realm-test"
	testKeycloakClientCRName         = "keycloak-client-test"
	testAuthZKeycloakClientCRName    = "authz-keycloak-client-test"
	testSecondKeycloakClientCRName   = "second-keycloak-client-test"
	testExternalKeycloakClientCRName = "external-keycloak-client-test"
	testKeycloakUserCRName           = "keycloak-user-test"
	cleanupRetryInterval             = time.Second * 5
	cleanupTimeout                   = time.Minute * 2
	pollRetryInterval                = time.Second * 10
	pollTimeout                      = time.Minute * 10
)
