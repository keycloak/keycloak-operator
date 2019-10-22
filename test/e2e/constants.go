package e2e

import "time"

const (
	testKeycloakCRName   = "keycloak-test"
	cleanupRetryInterval = time.Second * 1
	cleanupTimeout       = time.Second * 10
	pollRetryInterval    = time.Second * 15
	pollTimeout          = time.Second * 480
)
