package model

// Constants for a community Keycloak installation
const (
	ApplicationName                      = "keycloak"
	MonitoringKey                        = "middleware"
	DatabaseSecretName                   = ApplicationName + "-db-secret"
	PostgresqlPersistentVolumeName       = ApplicationName + "-postgresql-claim"
	PostgresqlBackupPersistentVolumeName = ApplicationName + "-backup"
	PostgresqlDeploymentName             = ApplicationName + "-postgresql"
	KeycloakProbesName                   = ApplicationName + "-probes"
	PostgresqlDeploymentComponent        = "database"
	PostgresqlServiceName                = ApplicationName + "-postgresql"
	PostgresqlImage                      = "postgres:11.5"
	KeycloakImage                        = "quay.io/keycloak/keycloak:9.0.2"
	KeycloakInitContainerImage           = "quay.io/keycloak/keycloak-init-container:master"
	RHSSOImage                           = "registry.redhat.io/rh-sso-7/sso74-openshift-rhel8:7.4-1"
	BackupImage                          = "quay.io/integreatly/backup-container:1.0.10"
	KeycloakDiscoveryServiceName         = ApplicationName + "-discovery"
	KeycloakDeploymentName               = ApplicationName
	KeycloakDeploymentComponent          = "keycloak"
	PostgresqlBackupComponent            = "database-backup"
	PostgresqlDatabase                   = "root"
	PostgresqlUsername                   = ApplicationName
	PostgresqlPasswordLength             = 32
	PostgresqlPersistentVolumeCapacity   = "1Gi"
	DatabaseSecretUsernameProperty       = "POSTGRES_USERNAME" // nolint
	DatabaseSecretPasswordProperty       = "POSTGRES_PASSWORD" // nolint
	// Required by the Integreately Backup Image
	DatabaseSecretHostProperty = "POSTGRES_HOST" // nolint
	// Required by the Integreately Backup Image
	DatabaseSecretDatabaseProperty = "POSTGRES_DATABASE" // nolint
	// Required by the Integreately Backup Image
	DatabaseSecretSuperuserProperty       = "POSTGRES_SUPERUSER"        // nolint
	DatabaseSecretExternalAddressProperty = "POSTGRES_EXTERNAL_ADDRESS" // nolint
	DatabaseSecretExternalPortProperty    = "POSTGRES_EXTERNAL_PORT"    // nolint
	KeycloakServicePort                   = 8443
	PostgresDefaultPort                   = 5432
	AdminUsernameProperty                 = "ADMIN_USERNAME"        // nolint
	AdminPasswordProperty                 = "ADMIN_PASSWORD"        // nolint
	ServingCertSecretName                 = "sso-x509-https-secret" // nolint
	LivenessProbeProperty                 = "liveness_probe.sh"
	ReadinessProbeProperty                = "readiness_probe.sh"
	RouteLoadBalancingStrategy            = "source"
	PostgresqlBackupServiceAccountName    = "keycloak-operator"
	KeycloakExtensionEnvVar               = "KEYCLOAK_EXTENSIONS"
	KeycloakExtensionPath                 = "/opt/jboss/keycloak/providers"
	KeycloakExtensionsInitContainerPath   = "/opt/extensions"
	RhssoExtensionPath                    = "/opt/eap/providers"
	ClientSecretName                      = ApplicationName + "-client-secret"
	ClientSecretClientIDProperty          = "CLIENT_ID"
	ClientSecretClientSecretProperty      = "CLIENT_SECRET"
	MaxUnavailableNumberOfPods            = 1
	PodMonitorName                        = ApplicationName + "-pod-monitor"
	ServiceMonitorName                    = ApplicationName + "-service-monitor"
)
