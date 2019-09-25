package keycloak

import (
	"reflect"
	"testing"

	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	integreatlyv1alpha1 "github.com/integr8ly/grafana-operator/pkg/apis/integreatly/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/common"
	"github.com/keycloak/keycloak-operator/pkg/model"
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
	//    0) Prometheus Rule
	//    1) Service Monitor
	//    2) Grafana Dashboard
	//    3) Postgresql Persistent Volume Claim
	//    4) Postgresql Deployment
	//    5) Postgresql Service
	//    6) Keycloak Service
	assert.Nil(t, error)
	assert.IsType(t, common.GenericCreateAction{}, desiredState[0])
	assert.IsType(t, common.GenericCreateAction{}, desiredState[1])
	assert.IsType(t, common.GenericCreateAction{}, desiredState[2])
	assert.IsType(t, common.GenericCreateAction{}, desiredState[3])
	assert.IsType(t, common.GenericCreateAction{}, desiredState[4])
	assert.IsType(t, common.GenericCreateAction{}, desiredState[5])
	assert.IsType(t, common.GenericCreateAction{}, desiredState[6])
	assert.IsType(t, model.PrometheusRule(cr), desiredState[0].(common.GenericCreateAction).Ref)
	assert.IsType(t, model.ServiceMonitor(cr), desiredState[1].(common.GenericCreateAction).Ref)
	assert.IsType(t, model.GrafanaDashboard(cr), desiredState[2].(common.GenericCreateAction).Ref)
	assert.IsType(t, model.PostgresqlPersistentVolumeClaim(cr), desiredState[3].(common.GenericCreateAction).Ref)
	assert.IsType(t, model.PostgresqlDeployment(cr), desiredState[4].(common.GenericCreateAction).Ref)
	assert.IsType(t, model.PostgresqlService(cr), desiredState[5].(common.GenericCreateAction).Ref)
	assert.IsType(t, model.KeycloakService(cr), desiredState[6].(common.GenericCreateAction).Ref)
}

func TestKeycloakReconciler_Test_Updating_All(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{}
	currentState := &common.ClusterState{
		KeycloakService:                 model.KeycloakService(cr),
		KeycloakServiceMonitor:          model.ServiceMonitor(cr),
		KeycloakPrometheusRule:          model.PrometheusRule(cr),
		KeycloakGrafanaDashboard:        model.GrafanaDashboard(cr),
		PostgresqlPersistentVolumeClaim: model.PostgresqlPersistentVolumeClaim(cr),
		PostgresqlService:               model.KeycloakService(cr),
		PostgresqlDeployment:            model.PostgresqlDeployment(cr),
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
	//    0) Prometheus Rule
	//    1) Service Monitor
	//    2) Grafana Dashboard
	//    3) Postgresql Persistent Volume Claim
	//    4) Postgresql Deployment
	//    5) Postgresql Service
	//    6) Keycloak Service
	assert.Nil(t, error)
	assert.IsType(t, common.GenericUpdateAction{}, desiredState[0])
	assert.IsType(t, common.GenericUpdateAction{}, desiredState[1])
	assert.IsType(t, common.GenericUpdateAction{}, desiredState[2])
	assert.IsType(t, common.GenericUpdateAction{}, desiredState[3])
	assert.IsType(t, common.GenericUpdateAction{}, desiredState[4])
	assert.IsType(t, common.GenericUpdateAction{}, desiredState[5])
	assert.IsType(t, common.GenericUpdateAction{}, desiredState[6])
	assert.IsType(t, model.PrometheusRule(cr), desiredState[0].(common.GenericUpdateAction).Ref)
	assert.IsType(t, model.ServiceMonitor(cr), desiredState[1].(common.GenericUpdateAction).Ref)
	assert.IsType(t, model.GrafanaDashboard(cr), desiredState[2].(common.GenericUpdateAction).Ref)
	assert.IsType(t, model.PostgresqlPersistentVolumeClaim(cr), desiredState[3].(common.GenericUpdateAction).Ref)
	assert.IsType(t, model.PostgresqlDeployment(cr), desiredState[4].(common.GenericUpdateAction).Ref)
	assert.IsType(t, model.PostgresqlService(cr), desiredState[5].(common.GenericUpdateAction).Ref)
	assert.IsType(t, model.KeycloakService(cr), desiredState[6].(common.GenericUpdateAction).Ref)
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
	desiredState, error := reconciler.Reconcile(currentState, cr)

	// then
	assert.Nil(t, error)
	for _, element := range desiredState {
		assert.IsType(t, common.GenericCreateAction{}, element)
		assert.NotEqual(t, reflect.TypeOf(model.PrometheusRule(cr)), reflect.TypeOf(element.(common.GenericCreateAction).Ref))
		assert.NotEqual(t, reflect.TypeOf(model.GrafanaDashboard(cr)), reflect.TypeOf(element.(common.GenericCreateAction).Ref))
		assert.NotEqual(t, reflect.TypeOf(model.ServiceMonitor(cr)), reflect.TypeOf(element.(common.GenericCreateAction).Ref))
	}
}
