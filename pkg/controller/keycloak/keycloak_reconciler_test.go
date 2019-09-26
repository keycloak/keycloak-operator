package keycloak

import (
	"testing"

	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	integreatlyv1alpha1 "github.com/integr8ly/grafana-operator/pkg/apis/integreatly/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/common"
	"github.com/keycloak/keycloak-operator/pkg/model/keycloak"
	"github.com/stretchr/testify/assert"
)

func TestKeycloakReconciler_Test_Creating_All(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{}
	currentState := common.NewClusterState()

	//Set monitoring resources exist to true
	stateManager := common.GetStateManager()
	stateManager.SetState(monitoringv1.PrometheusRuleKind, true)
	stateManager.SetState(monitoringv1.ServiceMonitorsKind, true)
	stateManager.SetState(integreatlyv1alpha1.GrafanaDashboardKind, true)

	// when
	reconciler := NewKeycloakReconciler()
	desiredState, error := reconciler.Reconcile(currentState, cr)

	// then
	// Expectation:
	//    1) Keycloak Service
	//    2) Service Monitor
	//    3) Prometheus Rule
	//    4) Grafana Dashboard
	assert.Nil(t, error)
	assert.IsType(t, common.GenericCreateAction{}, desiredState[0])
	assert.IsType(t, common.GenericCreateAction{}, desiredState[1])
	assert.IsType(t, common.GenericCreateAction{}, desiredState[2])
	assert.IsType(t, common.GenericCreateAction{}, desiredState[3])
	assert.IsType(t, keycloak.Service(cr), desiredState[0].(common.GenericCreateAction).Ref)
	assert.IsType(t, keycloak.ServiceMonitor(cr), desiredState[1].(common.GenericCreateAction).Ref)
	assert.IsType(t, keycloak.PrometheusRule(cr), desiredState[2].(common.GenericCreateAction).Ref)
	assert.IsType(t, keycloak.GrafanaDashboard(cr), desiredState[3].(common.GenericCreateAction).Ref)
}

func TestKeycloakReconciler_Test_Updating_All(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{}
	currentState := &common.ClusterState{
		KeycloakService:          keycloak.Service(cr),
		KeycloakServiceMonitor:   keycloak.ServiceMonitor(cr),
		KeycloakPrometheusRule:   keycloak.PrometheusRule(cr),
		KeycloakGrafanaDashboard: keycloak.GrafanaDashboard(cr),
	}

	//Set monitoring resources exist to true
	stateManager := common.GetStateManager()
	stateManager.SetState(monitoringv1.PrometheusRuleKind, true)
	stateManager.SetState(monitoringv1.ServiceMonitorsKind, true)
	stateManager.SetState(integreatlyv1alpha1.GrafanaDashboardKind, true)

	// when
	reconciler := NewKeycloakReconciler()
	desiredState, error := reconciler.Reconcile(currentState, cr)

	// then
	// Expectation:
	//    1) Keycloak Service
	//    2) Service Monitor
	//    3) Prometheus Rule
	//    4) Grafana Dashboard
	assert.Nil(t, error)
	assert.IsType(t, common.GenericUpdateAction{}, desiredState[0])
	assert.IsType(t, common.GenericUpdateAction{}, desiredState[1])
	assert.IsType(t, common.GenericUpdateAction{}, desiredState[2])
	assert.IsType(t, common.GenericUpdateAction{}, desiredState[3])
	assert.IsType(t, keycloak.Service(cr), desiredState[0].(common.GenericUpdateAction).Ref)
	assert.IsType(t, keycloak.ServiceMonitor(cr), desiredState[1].(common.GenericUpdateAction).Ref)
	assert.IsType(t, keycloak.PrometheusRule(cr), desiredState[2].(common.GenericUpdateAction).Ref)
	assert.IsType(t, keycloak.GrafanaDashboard(cr), desiredState[3].(common.GenericUpdateAction).Ref)
}

func TestKeycloakReconciler_Test_No_Action_When_Monitoring_Resources_Dont_Exist(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{}
	currentState := common.NewClusterState()

	//Set monitoring resources exist to true
	stateManager := common.GetStateManager()
	stateManager.SetState(monitoringv1.PrometheusRuleKind, false)
	stateManager.SetState(monitoringv1.ServiceMonitorsKind, false)
	stateManager.SetState(integreatlyv1alpha1.GrafanaDashboardKind, false)

	// when
	reconciler := NewKeycloakReconciler()
	prometheusRuleAction := reconciler.GetKeycloakPrometheusRuleDesiredState(currentState, cr)
	serviceMonitorAction := reconciler.GetKeycloakServiceMonitorDesiredState(currentState, cr)
	grafanaDashboardAction := reconciler.GetKeycloakGrafanaDashboardDesiredState(currentState, cr)

	// then
	// Expectation:
	//    nil returned from all functions
	assert.Nil(t, prometheusRuleAction)
	assert.Nil(t, serviceMonitorAction)
	assert.Nil(t, grafanaDashboardAction)
}
