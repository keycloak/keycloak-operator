package keycloak

import (
	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	integreatlyv1alpha1 "github.com/integr8ly/grafana-operator/pkg/apis/integreatly/v1alpha1"
	kc "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/common"
	"github.com/keycloak/keycloak-operator/pkg/model/keycloak"
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
	desired = desired.AddAction(i.getKeycloakServiceDesiredState(clusterState, cr))
	desired = desired.AddAction(i.getKeycloakServiceMonitorDesiredState(clusterState, cr))
	desired = desired.AddAction(i.getKeycloakPrometheusRuleDesiredState(clusterState, cr))
	desired = desired.AddAction(i.getKeycloakGrafanaDashboardDesiredState(clusterState, cr))

	return desired, nil
}

func (i *KeycloakReconciler) getKeycloakServiceDesiredState(clusterState *common.ClusterState, cr *kc.Keycloak) common.ClusterAction {
	stateManager := common.GetStateManager()
	resourceExists, keyExists := stateManager.GetState(monitoringv1.ServiceMonitorsKind).(bool)
	// Only add or update the monitoring resources if the resource type exists on the cluster. These booleans are set in the common/autodetect logic
	if !keyExists || !resourceExists {
		return nil
	}

	service := keycloak.Service(cr)

	if clusterState.KeycloakService == nil {
		return common.GenericCreateAction{
			Ref: service,
			Msg: "create keycloak service",
		}
	}

	// This part may change in the future once we have more resources to reconcile.
	// Perhaps there should be another method, like `keycloak.Service(cr, clusterState)`?
	service.Spec.ClusterIP = clusterState.KeycloakService.Spec.ClusterIP
	service.ResourceVersion = clusterState.KeycloakService.ResourceVersion
	return common.GenericUpdateAction{
		Ref: service,
		Msg: "update keycloak service",
	}
}

func (i *KeycloakReconciler) getKeycloakPrometheusRuleDesiredState(clusterState *common.ClusterState, cr *kc.Keycloak) common.ClusterAction {
	stateManager := common.GetStateManager()
	resourceExists, keyExists := stateManager.GetState(monitoringv1.PrometheusRuleKind).(bool)
	// Only add or update the monitoring resources if the resource type exists on the cluster. These booleans are set in the common/autodetect logic
	if !keyExists || !resourceExists {
		return nil
	}

	prometheusrule := keycloak.PrometheusRule(cr)

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

func (i *KeycloakReconciler) getKeycloakServiceMonitorDesiredState(clusterState *common.ClusterState, cr *kc.Keycloak) common.ClusterAction {
	stateManager := common.GetStateManager()
	resourceExists, keyExists := stateManager.GetState(integreatlyv1alpha1.GrafanaDashboardKind).(bool)
	// Only add or update the monitoring resources if the resource type exists on the cluster. These booleans are set in the common/autodetect logic
	if !keyExists || !resourceExists {
		return nil
	}

	servicemonitor := keycloak.ServiceMonitor(cr)

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

func (i *KeycloakReconciler) getKeycloakGrafanaDashboardDesiredState(clusterState *common.ClusterState, cr *kc.Keycloak) common.ClusterAction {
	grafanadashboard := keycloak.GrafanaDashboard(cr)

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
