package keycloak

import (
	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	integreatlyv1alpha1 "github.com/integr8ly/grafana-operator/pkg/apis/integreatly/v1alpha1"
	kc "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/common"
	"github.com/keycloak/keycloak-operator/pkg/model/keycloak"
	"github.com/spf13/viper"
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
	desired = append(desired, i.getKeycloakServiceDesiredState(clusterState, cr))

	// Only add or update the monitoring resources if the resource type exists on the cluster. These booleans are set in the common/autodetect logic
	if viper.GetBool(monitoringv1.ServiceMonitorsKind) {
		desired = append(desired, i.getKeycloakServiceMonitorDesiredState(clusterState, cr))
	}
	if viper.GetBool(monitoringv1.PrometheusRuleKind) {
		desired = append(desired, i.getKeycloakPrometheusRuleDesiredState(clusterState, cr))
	}
	if viper.GetBool(integreatlyv1alpha1.GrafanaDashboardKind) {
		desired = append(desired, i.getKeycloakGrafanaDashboardDesiredState(clusterState, cr))
	}

	return desired, nil
}

func (i *KeycloakReconciler) getKeycloakServiceDesiredState(clusterState *common.ClusterState, cr *kc.Keycloak) common.ClusterAction {
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
