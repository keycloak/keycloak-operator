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
	KeycloakImage                      = "jboss/keycloak:7.0.0"
	RHSSOImage                         = "registry.redhat.io/redhat-sso-7-tech-preview/sso-cd-openshift:6"
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
