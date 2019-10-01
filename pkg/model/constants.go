package model

// Constants for a community Keycloak installation
const (
	ApplicationName                     = "keycloak"
	MonitoringKey                       = "middleware"
	DatabaseSecretName                  = ApplicationName + "-db-secret"
	PostgresqlPersistentVolumeName      = ApplicationName + "-postgresql-claim"
	PostgresqlDeploymentName            = ApplicationName + "-postgresql"
	PostgresqlDeploymentComponent       = "database"
	PostgresqlServiceName               = ApplicationName + "-postgresql"
	PostgresqlImage                     = "postgres:12"
	KeycloakImage                       = "jboss/keycloak"
	KeycloakDiscoveryServiceName        = ApplicationName + "-discovery"
	KeycloakDeploymentName              = ApplicationName
	KeycloakDeploymentComponent         = "keycloak"
	KeycloakSecretName                  = ApplicationName + "-secret"
	PostgresqlDatabase                  = "root"
	PostgresqlPersistentVolumeCapacity  = "1Gi"
	DatabaseSecretUsernameProperty      = "user"
	DatabaseSecretPasswordProperty      = "password"
	KeycloakSecretAdminUsernameProperty = "KEYCLOAK_USER"
	KeycloakSecretAdminPasswordProperty = "KEYCLOAK_PASSWORD"
)
