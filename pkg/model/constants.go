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
	PostgresqlImage                    = "postgres:12"
	KeycloakImage                      = "jboss/keycloak"
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
)
