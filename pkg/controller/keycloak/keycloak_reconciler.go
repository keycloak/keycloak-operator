package keycloak

import (
	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	grafanav1alpha1 "github.com/integr8ly/grafana-operator/v3/pkg/apis/integreatly/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	kc "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/common"
	"github.com/keycloak/keycloak-operator/pkg/model"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Reconciler interface {
	Reconcile(clusterState *common.ClusterState, cr *kc.Keycloak) (common.DesiredClusterState, error)
}

type KeycloakReconciler struct { // nolint
}

func NewKeycloakReconciler() *KeycloakReconciler {
	return &KeycloakReconciler{}
}

func (i *KeycloakReconciler) Reconcile(clusterState *common.ClusterState, cr *kc.Keycloak) common.DesiredClusterState {
	desired := common.DesiredClusterState{}

	desired = desired.AddAction(i.GetKeycloakAdminSecretDesiredState(clusterState, cr))
	desired = desired.AddAction(i.GetKeycloakPrometheusRuleDesiredState(clusterState, cr))
	desired = desired.AddAction(i.GetKeycloakServiceMonitorDesiredState(clusterState, cr))
	desired = desired.AddAction(i.GetKeycloakGrafanaDashboardDesiredState(clusterState, cr))

	if !cr.Spec.ExternalDatabase.Enabled {
		desired = desired.AddAction(i.getDatabaseSecretDesiredState(clusterState, cr))
		desired = desired.AddAction(i.getPostgresqlPersistentVolumeClaimDesiredState(clusterState, cr))
		desired = desired.AddAction(i.getPostgresqlDeploymentDesiredState(clusterState, cr))
		desired = desired.AddAction(i.getPostgresqlServiceDesiredState(clusterState, cr, false))
	} else {
		i.reconcileExternalDatabase(&desired, clusterState, cr)
	}

	desired = desired.AddAction(i.getKeycloakServiceDesiredState(clusterState, cr))
	desired = desired.AddAction(i.getKeycloakDiscoveryServiceDesiredState(clusterState, cr))
	desired = desired.AddAction(i.getKeycloakMonitoringServiceDesiredState(clusterState, cr))
	desired = desired.AddAction(i.GetKeycloakProbesDesiredState(clusterState, cr))
	desired = desired.AddAction(i.getKeycloakDeploymentOrRHSSODesiredState(clusterState, cr))
	i.reconcileExternalAccess(&desired, clusterState, cr)
	desired = desired.AddAction(i.getPodDisruptionBudgetDesiredState(clusterState, cr))

	if cr.Spec.Migration.Backups.Enabled {
		desired = desired.AddAction(i.getKeycloakBackupDesiredState(clusterState, cr))
	}
	return desired
}

func (i *KeycloakReconciler) reconcileExternalDatabase(desired *common.DesiredClusterState, clusterState *common.ClusterState, cr *kc.Keycloak) {
	// If the database secret does not exist we can't continue
	if clusterState.DatabaseSecret == nil {
		return
	}
	if model.IsIP(clusterState.DatabaseSecret.Data[model.DatabaseSecretExternalAddressProperty]) {
		// If the address of the external database is an IP address then we have to
		// set up an endpoints object for the service to send traffic. An externalName
		// type service won't work in this case. For more details, see https://cloud.google.com/blog/products/gcp/kubernetes-best-practices-mapping-external-services
		desired.AddAction(i.getPostgresqlServiceEndpointsDesiredState(clusterState, cr))
		desired.AddAction(i.getPostgresqlServiceDesiredState(clusterState, cr, false))
	} else {
		// If we have an URI for the external database then we can use a service of
		// type externalName
		desired.AddAction(i.getPostgresqlServiceDesiredState(clusterState, cr, true))
	}
}

func (i *KeycloakReconciler) reconcileExternalAccess(desired *common.DesiredClusterState, clusterState *common.ClusterState, cr *kc.Keycloak) {
	if !cr.Spec.ExternalAccess.Enabled {
		return
	}

	// Find out if we're on OpenShift or Kubernetes and create either a Route or
	// an Ingress
	stateManager := common.GetStateManager()
	openshift, keyExists := stateManager.GetState(common.OpenShiftAPIServerKind).(bool)

	if keyExists && openshift {
		desired.AddAction(i.getKeycloakRouteDesiredState(clusterState, cr))
		desired.AddAction(i.getKeycloakMetricsRouteDesiredState(clusterState, cr))
	} else {
		desired.AddAction(i.getKeycloakIngressDesiredState(clusterState, cr))
	}
}

func (i *KeycloakReconciler) GetKeycloakAdminSecretDesiredState(clusterState *common.ClusterState, cr *kc.Keycloak) common.ClusterAction {
	keycloakAdminSecret := model.KeycloakAdminSecret(cr)

	if clusterState.KeycloakAdminSecret == nil {
		return common.GenericCreateAction{
			Ref: keycloakAdminSecret,
			Msg: "Create Keycloak admin secret",
		}
	}
	return common.GenericUpdateAction{
		Ref: model.KeycloakAdminSecretReconciled(cr, clusterState.KeycloakAdminSecret),
		Msg: "Update Keycloak admin secret",
	}
}

func (i *KeycloakReconciler) GetKeycloakProbesDesiredState(clusterState *common.ClusterState, cr *kc.Keycloak) common.ClusterAction {
	keycloakProbesConfigMap := model.KeycloakProbes(cr)

	if clusterState.KeycloakProbes == nil {
		return common.GenericCreateAction{
			Ref: keycloakProbesConfigMap,
			Msg: "Create Keycloak probes configmap",
		}
	}
	return nil
}

func (i *KeycloakReconciler) getPostgresqlPersistentVolumeClaimDesiredState(clusterState *common.ClusterState, cr *kc.Keycloak) common.ClusterAction {
	postgresqlPersistentVolume := model.PostgresqlPersistentVolumeClaim(cr)
	if clusterState.PostgresqlPersistentVolumeClaim == nil {
		return common.GenericCreateAction{
			Ref: postgresqlPersistentVolume,
			Msg: "Create Postgresql PersistentVolumeClaim",
		}
	}
	return common.GenericUpdateAction{
		Ref: model.PostgresqlPersistentVolumeClaimReconciled(cr, clusterState.PostgresqlPersistentVolumeClaim),
		Msg: "Update Postgresql PersistentVolumeClaim",
	}
}

func (i *KeycloakReconciler) getPostgresqlServiceDesiredState(clusterState *common.ClusterState, cr *kc.Keycloak, isExternal bool) common.ClusterAction {
	postgresqlService := model.PostgresqlService(cr, clusterState.DatabaseSecret, isExternal)
	if clusterState.PostgresqlService == nil {
		return common.GenericCreateAction{
			Ref: postgresqlService,
			Msg: "Create Postgresql KeycloakService",
		}
	}
	return common.GenericUpdateAction{
		Ref: model.PostgresqlServiceReconciled(clusterState.PostgresqlService, clusterState.DatabaseSecret, isExternal),
		Msg: "Update Postgresql KeycloakService",
	}
}

func (i *KeycloakReconciler) getPostgresqlDeploymentDesiredState(clusterState *common.ClusterState, cr *kc.Keycloak) common.ClusterAction {
	// Find out if we're on OpenShift or Kubernetes
	stateManager := common.GetStateManager()
	isOpenshift, _ := stateManager.GetState(common.OpenShiftAPIServerKind).(bool)

	postgresqlDeployment := model.PostgresqlDeployment(cr, isOpenshift)

	if clusterState.PostgresqlDeployment == nil {
		return common.GenericCreateAction{
			Ref: postgresqlDeployment,
			Msg: "Create Postgresql Deployment",
		}
	}
	return common.GenericUpdateAction{
		Ref: model.PostgresqlDeploymentReconciled(cr, clusterState.PostgresqlDeployment),
		Msg: "Update Postgresql Deployment",
	}
}

func (i *KeycloakReconciler) getKeycloakServiceDesiredState(clusterState *common.ClusterState, cr *kc.Keycloak) common.ClusterAction {
	keycloakService := model.KeycloakService(cr)

	if clusterState.KeycloakService == nil {
		return common.GenericCreateAction{
			Ref: keycloakService,
			Msg: "Create Keycloak Service",
		}
	}
	return common.GenericUpdateAction{
		Ref: model.KeycloakServiceReconciled(cr, clusterState.KeycloakService),
		Msg: "Update keycloak Service",
	}
}

func (i *KeycloakReconciler) getKeycloakDiscoveryServiceDesiredState(clusterState *common.ClusterState, cr *kc.Keycloak) common.ClusterAction {
	keycloakDiscoveryService := model.KeycloakDiscoveryService(cr)

	if clusterState.KeycloakDiscoveryService == nil {
		return common.GenericCreateAction{
			Ref: keycloakDiscoveryService,
			Msg: "Create Keycloak Discovery Service",
		}
	}
	return common.GenericUpdateAction{
		Ref: model.KeycloakDiscoveryServiceReconciled(cr, clusterState.KeycloakDiscoveryService),
		Msg: "Update keycloak Discovery Service",
	}
}

func (i *KeycloakReconciler) getKeycloakMonitoringServiceDesiredState(clusterState *common.ClusterState, cr *kc.Keycloak) common.ClusterAction {
	stateManager := common.GetStateManager()
	resourceWatchExists, keyExists := stateManager.GetState(common.GetStateFieldName(ControllerName, monitoringv1.ServiceMonitorsKind)).(bool)
	// Only add or update the monitoring resources if the resource type exists on the cluster. These booleans are set in the common/autodetect logic
	if !keyExists || !resourceWatchExists {
		return nil
	}

	keycloakMonitoringService := model.KeycloakMonitoringService(cr)

	if clusterState.KeycloakMonitoringService == nil {
		return common.GenericCreateAction{
			Ref: keycloakMonitoringService,
			Msg: "Create Keycloak Monitoring Service",
		}
	}
	return common.GenericUpdateAction{
		Ref: model.KeycloakMonitoringServiceReconciled(cr, clusterState.KeycloakMonitoringService),
		Msg: "Update keycloak Monitoring Service",
	}
}

func (i *KeycloakReconciler) GetKeycloakPrometheusRuleDesiredState(clusterState *common.ClusterState, cr *kc.Keycloak) common.ClusterAction {
	stateManager := common.GetStateManager()
	resourceWatchExists, keyExists := stateManager.GetState(common.GetStateFieldName(ControllerName, monitoringv1.PrometheusRuleKind)).(bool)
	// Only add or update the monitoring resources if the resource type exists on the cluster. These booleans are set in the common/autodetect logic
	if !keyExists || !resourceWatchExists {
		return nil
	}

	prometheusrule := model.PrometheusRule(cr)

	if clusterState.KeycloakPrometheusRule == nil {
		return common.GenericCreateAction{
			Ref: prometheusrule,
			Msg: "create keycloak prometheus rule",
		}
	}

	prometheusrule.ResourceVersion = clusterState.KeycloakPrometheusRule.ResourceVersion
	return common.GenericUpdateAction{
		Ref: prometheusrule,
		Msg: "update keycloak prometheus rule",
	}
}

func (i *KeycloakReconciler) GetKeycloakServiceMonitorDesiredState(clusterState *common.ClusterState, cr *kc.Keycloak) common.ClusterAction {
	stateManager := common.GetStateManager()
	resourceWatchExists, keyExists := stateManager.GetState(common.GetStateFieldName(ControllerName, monitoringv1.ServiceMonitorsKind)).(bool)
	// Only add or update the monitoring resources if the resource type exists on the cluster. These booleans are set in the common/autodetect logic
	if !keyExists || !resourceWatchExists {
		return nil
	}

	servicemonitor := model.ServiceMonitor(cr)

	if clusterState.KeycloakServiceMonitor == nil {
		return common.GenericCreateAction{
			Ref: servicemonitor,
			Msg: "create keycloak service monitor",
		}
	}

	servicemonitor.ResourceVersion = clusterState.KeycloakServiceMonitor.ResourceVersion
	return common.GenericUpdateAction{
		Ref: servicemonitor,
		Msg: "update keycloak service monitor",
	}
}

func (i *KeycloakReconciler) GetKeycloakGrafanaDashboardDesiredState(clusterState *common.ClusterState, cr *kc.Keycloak) common.ClusterAction {
	stateManager := common.GetStateManager()
	resourceWatchExists, keyExists := stateManager.GetState(common.GetStateFieldName(ControllerName, grafanav1alpha1.GrafanaDashboardKind)).(bool)
	// Only add or update the monitoring resources if the resource type exists on the cluster. These booleans are set in the common/autodetect logic
	if !keyExists || !resourceWatchExists {
		return nil
	}

	grafanadashboard := model.GrafanaDashboard(cr)

	if clusterState.KeycloakGrafanaDashboard == nil {
		return common.GenericCreateAction{
			Ref: grafanadashboard,
			Msg: "create keycloak grafana dashboard",
		}
	}

	return common.GenericUpdateAction{
		Ref: model.GrafanaDashboardReconciled(cr, clusterState.KeycloakGrafanaDashboard),
		Msg: "update keycloak grafana dashboard",
	}
}

func (i *KeycloakReconciler) getDatabaseSecretDesiredState(clusterState *common.ClusterState, cr *kc.Keycloak) common.ClusterAction {
	databaseSecret := model.DatabaseSecret(cr)
	if clusterState.DatabaseSecret == nil {
		return common.GenericCreateAction{
			Ref: databaseSecret,
			Msg: "Create Database Secret",
		}
	}
	return common.GenericUpdateAction{
		Ref: model.DatabaseSecretReconciled(cr, clusterState.DatabaseSecret),
		Msg: "Update Database Secret",
	}
}

func (i *KeycloakReconciler) getKeycloakDeploymentOrRHSSODesiredState(clusterState *common.ClusterState, cr *kc.Keycloak) common.ClusterAction {
	isRHSSO := model.Profiles.IsRHSSO(cr)

	deployment := model.KeycloakDeployment(cr, clusterState.DatabaseSecret)
	deploymentName := "Keycloak"

	if isRHSSO {
		deployment = model.RHSSODeployment(cr, clusterState.DatabaseSecret)
		deploymentName = model.RHSSOProfile
	}

	if clusterState.KeycloakDeployment == nil {
		return common.GenericCreateAction{
			Ref: deployment,
			Msg: "Create " + deploymentName + " Deployment (StatefulSet)",
		}
	}

	deploymentReconciled := model.KeycloakDeploymentReconciled(cr, clusterState.KeycloakDeployment, clusterState.DatabaseSecret)
	if isRHSSO {
		deploymentReconciled = model.RHSSODeploymentReconciled(cr, clusterState.KeycloakDeployment, clusterState.DatabaseSecret)
	}

	return common.GenericUpdateAction{
		Ref: deploymentReconciled,
		Msg: "Update " + deploymentName + " Deployment (StatefulSet)",
	}
}

func (i *KeycloakReconciler) getKeycloakRouteDesiredState(clusterState *common.ClusterState, cr *kc.Keycloak) common.ClusterAction {
	if clusterState.KeycloakRoute == nil {
		return common.GenericCreateAction{
			Ref: model.KeycloakRoute(cr),
			Msg: "Create Keycloak Route",
		}
	}

	return common.GenericUpdateAction{
		Ref: model.KeycloakRouteReconciled(cr, clusterState.KeycloakRoute),
		Msg: "Update Keycloak Route",
	}
}

func (i *KeycloakReconciler) getKeycloakMetricsRouteDesiredState(clusterState *common.ClusterState, cr *kc.Keycloak) common.ClusterAction {
	if clusterState.KeycloakRoute == nil {
		return nil
	}

	if clusterState.KeycloakMetricsRoute == nil {
		return common.GenericCreateAction{
			Ref: model.KeycloakMetricsRoute(cr, clusterState.KeycloakRoute),
			Msg: "Create Keycloak Metrics Route",
		}
	}

	return common.GenericUpdateAction{
		Ref: model.KeycloakMetricsRouteReconciled(cr, clusterState.KeycloakMetricsRoute, clusterState.KeycloakRoute),
		Msg: "Update Keycloak Metrics Route",
	}
}

func (i *KeycloakReconciler) getKeycloakIngressDesiredState(clusterState *common.ClusterState, cr *kc.Keycloak) common.ClusterAction {
	if clusterState.KeycloakIngress == nil {
		return common.GenericCreateAction{
			Ref: model.KeycloakIngress(cr),
			Msg: "Create Keycloak Ingress",
		}
	}

	return common.GenericUpdateAction{
		Ref: model.KeycloakIngressReconciled(cr, clusterState.KeycloakIngress),
		Msg: "Update Keycloak Ingress",
	}
}

func (i *KeycloakReconciler) getPostgresqlServiceEndpointsDesiredState(clusterState *common.ClusterState, cr *kc.Keycloak) common.ClusterAction {
	if clusterState.PostgresqlServiceEndpoints == nil {
		// This happens only during initial run
		return nil
	}
	return common.GenericUpdateAction{
		Ref: model.PostgresqlServiceEndpointsReconciled(cr, clusterState.PostgresqlServiceEndpoints, clusterState.DatabaseSecret),
		Msg: "Update External Database Service Endpoints",
	}
}

func (i *KeycloakReconciler) getPodDisruptionBudgetDesiredState(clusterState *common.ClusterState, cr *kc.Keycloak) common.ClusterAction {
	if cr.Spec.PodDisruptionBudget.Enabled {
		if clusterState.PodDisruptionBudget == nil {
			return common.GenericCreateAction{
				Ref: model.PodDisruptionBudget(cr),
				Msg: "Create PodDisruptionBudget",
			}
		}
		return common.GenericUpdateAction{
			Ref: model.PodDisruptionBudgetReconciled(cr, clusterState.PodDisruptionBudget),
			Msg: "Update PodDisruptionBudget",
		}
	}
	return nil
}

func (i *KeycloakReconciler) getKeycloakBackupDesiredState(clusterState *common.ClusterState, cr *kc.Keycloak) common.ClusterAction {
	backupCr := &v1alpha1.KeycloakBackup{}
	backupCr.Namespace = cr.Namespace
	backupCr.Name = model.MigrateBackupName + "-" + common.BackupTime
	labelSelect := metav1.LabelSelector{
		MatchLabels: cr.Labels,
	}
	backupCr.Spec.InstanceSelector = &labelSelect
	backupCr.Spec.StorageClassName = cr.Spec.StorageClassName

	if clusterState.KeycloakBackup == nil {
		// This happens before migration
		return nil
	}

	keycloakbackup := model.KeycloakMigrationOneTimeBackup(backupCr)
	keycloakbackup.ResourceVersion = clusterState.KeycloakBackup.ResourceVersion
	return common.GenericUpdateAction{
		Ref: keycloakbackup,
		Msg: "Update Postgresql Backup for Keycloak Migration",
	}
}
