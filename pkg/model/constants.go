package model

// Constants for a community Keycloak installation
const (
	ApplicationName                    = "keycloak"
	MonitoringKey                      = "middleware"
	DatabaseSecretName                 = ApplicationName + "-db-secret"
	PostgresqlPersistentVolumeName     = ApplicationName + "-postgresql-claim"
	PostgresqlDeploymentName           = ApplicationName + "-postgresql"
	PostgresqlDeploymentComponent      = "database"
	PostgresqlServiceName              = ApplicationName + "-postgresql"
	PostgresqlImage                    = "postgres:9.5"
	KeycloakImage                      = "quay.io/keycloak/keycloak:7.0.1"
	RHSSOImage                         = "registry.access.redhat.com/redhat-sso-7/sso73-openshift:1.0"
	KeycloakDiscoveryServiceName       = ApplicationName + "-discovery"
	KeycloakDeploymentName             = ApplicationName
	KeycloakDeploymentComponent        = "keycloak"
	PostgresqlDatabase                 = "root"
	PostgresqlPersistentVolumeCapacity = "1Gi"
	DatabaseSecretUsernameProperty     = "user"
	DatabaseSecretPasswordProperty     = "password"
	KeycloakServicePort                = 8443
	AdminUsernameProperty              = "ADMIN_USERNAME"        // nolint
	AdminPasswordProperty              = "ADMIN_PASSWORD"        // nolint
	ServingCertSecretName              = "sso-x509-https-secret" // nolint
	RouteLoadBalancingStrategy         = "source"
)
