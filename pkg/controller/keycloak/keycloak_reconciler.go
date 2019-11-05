package keycloak

import (
	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	integreatlyv1alpha1 "github.com/integr8ly/grafana-operator/pkg/apis/integreatly/v1alpha1"
	kc "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/common"
	"github.com/keycloak/keycloak-operator/pkg/model"
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
	} else {
		desired = desired.AddAction(i.getPostgresqlServiceEndpointsDesiredState(clusterState, cr))
	}
	desired = desired.AddAction(i.getPostgresqlServiceDesiredState(clusterState, cr))

	desired = desired.AddAction(i.getKeycloakServiceDesiredState(clusterState, cr))
	desired = desired.AddAction(i.getKeycloakDiscoveryServiceDesiredState(clusterState, cr))
	desired = desired.AddAction(i.getKeycloakDeploymentOrRHSSODesiredState(clusterState, cr))
	i.reconcileExternalAccess(&desired, clusterState, cr)

	return desired
}

func (i *KeycloakReconciler) reconcileExternalAccess(desired *common.DesiredClusterState, clusterState *common.ClusterState, cr *kc.Keycloak) {
	if !cr.Spec.ExternalAccess.Enabled {
		return
	}

	// Find out if we're on OpenShift or Kubernetes and create either a Route or
	// an Ingress
	stateManager := common.GetStateManager()
	openshift, keyExists := stateManager.GetState(common.RouteKind).(bool)

	if keyExists && openshift {
		desired.AddAction(i.getKeycloakRouteDesiredState(clusterState, cr))
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

func (i *KeycloakReconciler) getPostgresqlServiceDesiredState(clusterState *common.ClusterState, cr *kc.Keycloak) common.ClusterAction {
	postgresqlService := model.PostgresqlService(cr)
	if clusterState.PostgresqlService == nil {
		return common.GenericCreateAction{
			Ref: postgresqlService,
			Msg: "Create Postgresql KeycloakService",
		}
	}
	return common.GenericUpdateAction{
		Ref: model.PostgresqlServiceReconciled(cr, clusterState.PostgresqlService),
		Msg: "Update Postgresql KeycloakService",
	}
}

func (i *KeycloakReconciler) getPostgresqlDeploymentDesiredState(clusterState *common.ClusterState, cr *kc.Keycloak) common.ClusterAction {
	postgresqlDeployment := model.PostgresqlDeployment(cr)
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
	resourceWatchExists, keyExists := stateManager.GetState(common.GetStateFieldName(ControllerName, integreatlyv1alpha1.GrafanaDashboardKind)).(bool)
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
	isRHSSO := cr.Spec.Profile == common.RHSSOProfile

	deployment := model.KeycloakDeployment(cr)
	deploymentName := "Keycloak"

	if isRHSSO {
		deployment = model.RHSSODeployment(cr)
		deploymentName = common.RHSSOProfile
	}

	if clusterState.KeycloakDeployment == nil {
		return common.GenericCreateAction{
			Ref: deployment,
			Msg: "Create " + deploymentName + " Deployment (StatefulSet)",
		}
	}

	deploymentReconciled := model.KeycloakDeploymentReconciled(cr, clusterState.KeycloakDeployment)
	if isRHSSO {
		deploymentReconciled = model.RHSSODeploymentReconciled(cr, clusterState.KeycloakDeployment)
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
