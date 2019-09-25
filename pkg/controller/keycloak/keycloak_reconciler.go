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

func (i *KeycloakReconciler) Reconcile(clusterState *common.ClusterState, cr *kc.Keycloak) (common.DesiredClusterState, error) {
	desired := common.DesiredClusterState{}
	desired = desired.AddAction(i.GetKeycloakPrometheusRuleDesiredState(clusterState, cr))
	desired = desired.AddAction(i.GetKeycloakServiceMonitorDesiredState(clusterState, cr))
	desired = desired.AddAction(i.GetKeycloakGrafanaDashboardDesiredState(clusterState, cr))
	desired = desired.AddAction(i.getPostgresqlPersistentVolumeClaimDesiredState(clusterState, cr))
	desired = desired.AddAction(i.getPostgresqlDeploymentDesiredState(clusterState, cr))
	desired = desired.AddAction(i.getPostgresqlServiceDesiredState(clusterState, cr))
	desired = desired.AddAction(i.getKeycloakServiceDesiredState(clusterState, cr))

	return desired, nil
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

func (i *KeycloakReconciler) GetKeycloakPrometheusRuleDesiredState(clusterState *common.ClusterState, cr *kc.Keycloak) common.ClusterAction {
	stateManager := common.GetStateManager()
	resourceExists, keyExists := stateManager.GetState(monitoringv1.PrometheusRuleKind).(bool)
	// Only add or update the monitoring resources if the resource type exists on the cluster. These booleans are set in the common/autodetect logic
	if !keyExists || !resourceExists {
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
	resourceExists, keyExists := stateManager.GetState(monitoringv1.ServiceMonitorsKind).(bool)
	// Only add or update the monitoring resources if the resource type exists on the cluster. These booleans are set in the common/autodetect logic
	if !keyExists || !resourceExists {
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
	resourceExists, keyExists := stateManager.GetState(integreatlyv1alpha1.GrafanaDashboardKind).(bool)
	// Only add or update the monitoring resources if the resource type exists on the cluster. These booleans are set in the common/autodetect logic
	if !keyExists || !resourceExists {
		return nil
	}

	grafanadashboard := model.GrafanaDashboard(cr)

	if clusterState.KeycloakGrafanaDashboard == nil {
		return common.GenericCreateAction{
			Ref: grafanadashboard,
			Msg: "create keycloak grafana dashboard",
		}
	}

	grafanadashboard.ResourceVersion = clusterState.KeycloakGrafanaDashboard.ResourceVersion
	return common.GenericUpdateAction{
		Ref: grafanadashboard,
		Msg: "update keycloak grafana dashboard",
	}
}
