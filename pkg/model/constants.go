package model

// Constants for a community Keycloak installation
const (
	ApplicationName                = "keycloak"
	MonitoringKey                  = "middleware"
	PostgresqlPersistentVolumeName = ApplicationName + "-postgresql-claim"
	PostgresqlDeploymentName       = ApplicationName + "-postgresql"
	PostgresqlServiceName          = ApplicationName + "-postgresql"
	PostgresqlImage                = "postgres:12"
	//FIXME: This probably needs to be externalized to a secret
	PostgresqlUsername = "admin"
	//FIXME: This probably needs to be externalized to a secret
	PostgresqlPassword                 = "admin"
	PostgresqlDatabase                 = "root"
	PostgresqlPersistentVolumeCapacity = "1Gi"
)
