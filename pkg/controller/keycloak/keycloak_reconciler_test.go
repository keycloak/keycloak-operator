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
	stateManager.SetState(getStateFieldName(monitoringv1.PrometheusRuleKind), true)
	stateManager.SetState(getStateFieldName(monitoringv1.ServiceMonitorsKind), true)
	stateManager.SetState(getStateFieldName(integreatlyv1alpha1.GrafanaDashboardKind), true)

	// when
	reconciler := NewKeycloakReconciler()
	desiredState := reconciler.Reconcile(currentState, cr)

	// then
	// Expectation:
	//    0) Prometheus Rule
	//    1) Service Monitor
	//    2) Grafana Dashboard
	//    3) Database secret
	//    4) Postgresql Persistent Volume Claim
	//    5) Postgresql Deployment
	//    6) Postgresql Service
	//    7) Keycloak Service
	//    8) Keycloak Discovery Service
	//    9) Keycloak StatefulSets
	assert.IsType(t, common.GenericCreateAction{}, desiredState[0])
	assert.IsType(t, common.GenericCreateAction{}, desiredState[1])
	assert.IsType(t, common.GenericCreateAction{}, desiredState[2])
	assert.IsType(t, common.GenericCreateAction{}, desiredState[3])
	assert.IsType(t, common.GenericCreateAction{}, desiredState[4])
	assert.IsType(t, common.GenericCreateAction{}, desiredState[5])
	assert.IsType(t, common.GenericCreateAction{}, desiredState[6])
	assert.IsType(t, common.GenericCreateAction{}, desiredState[7])
	assert.IsType(t, common.GenericCreateAction{}, desiredState[8])
	assert.IsType(t, common.GenericCreateAction{}, desiredState[9])
	assert.IsType(t, common.GenericCreateAction{}, desiredState[10])
	assert.IsType(t, model.KeycloakAdminSecret(cr), desiredState[0].(common.GenericCreateAction).Ref)
	assert.IsType(t, model.PrometheusRule(cr), desiredState[1].(common.GenericCreateAction).Ref)
	assert.IsType(t, model.ServiceMonitor(cr), desiredState[2].(common.GenericCreateAction).Ref)
	assert.IsType(t, model.GrafanaDashboard(cr), desiredState[3].(common.GenericCreateAction).Ref)
	assert.IsType(t, model.DatabaseSecret(cr), desiredState[4].(common.GenericCreateAction).Ref)
	assert.IsType(t, model.PostgresqlPersistentVolumeClaim(cr), desiredState[5].(common.GenericCreateAction).Ref)
	assert.IsType(t, model.PostgresqlDeployment(cr), desiredState[6].(common.GenericCreateAction).Ref)
	assert.IsType(t, model.PostgresqlService(cr), desiredState[7].(common.GenericCreateAction).Ref)
	assert.IsType(t, model.KeycloakService(cr), desiredState[8].(common.GenericCreateAction).Ref)
	assert.IsType(t, model.KeycloakDiscoveryService(cr), desiredState[9].(common.GenericCreateAction).Ref)
	assert.IsType(t, model.KeycloakDeployment(cr), desiredState[10].(common.GenericCreateAction).Ref)
}

func TestKeycloakReconciler_Test_Updating_All(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{}
	currentState := &common.ClusterState{
		KeycloakServiceMonitor:          model.ServiceMonitor(cr),
		KeycloakPrometheusRule:          model.PrometheusRule(cr),
		KeycloakGrafanaDashboard:        model.GrafanaDashboard(cr),
		DatabaseSecret:                  model.DatabaseSecret(cr),
		PostgresqlPersistentVolumeClaim: model.PostgresqlPersistentVolumeClaim(cr),
		PostgresqlService:               model.PostgresqlService(cr),
		PostgresqlDeployment:            model.PostgresqlDeployment(cr),
		KeycloakService:                 model.KeycloakService(cr),
		KeycloakDiscoveryService:        model.KeycloakDiscoveryService(cr),
		KeycloakDeployment:              model.KeycloakDeployment(cr),
		KeycloakAdminSecret:             model.KeycloakAdminSecret(cr),
	}

	//Set monitoring resources exist to true
	stateManager := common.GetStateManager()
	stateManager.SetState(getStateFieldName(monitoringv1.PrometheusRuleKind), true)
	stateManager.SetState(getStateFieldName(monitoringv1.ServiceMonitorsKind), true)
	stateManager.SetState(getStateFieldName(integreatlyv1alpha1.GrafanaDashboardKind), true)

	// when
	reconciler := NewKeycloakReconciler()
	desiredState := reconciler.Reconcile(currentState, cr)

	// then
	// Expectation:
	//    0) Prometheus Rule
	//    1) Service Monitor
	//    2) Grafana Dashboard
	//    3) Database secret
	//    4) Postgresql Persistent Volume Claim
	//    5) Postgresql Deployment
	//    6) Postgresql Service
	//    7) Keycloak Service
	//    8) Keycloak Discovery Service
	//    9) Keycloak StatefulSets
	assert.IsType(t, common.GenericUpdateAction{}, desiredState[0])
	assert.IsType(t, common.GenericUpdateAction{}, desiredState[1])
	assert.IsType(t, common.GenericUpdateAction{}, desiredState[2])
	assert.IsType(t, common.GenericUpdateAction{}, desiredState[3])
	assert.IsType(t, common.GenericUpdateAction{}, desiredState[4])
	assert.IsType(t, common.GenericUpdateAction{}, desiredState[5])
	assert.IsType(t, common.GenericUpdateAction{}, desiredState[6])
	assert.IsType(t, common.GenericUpdateAction{}, desiredState[7])
	assert.IsType(t, common.GenericUpdateAction{}, desiredState[8])
	assert.IsType(t, common.GenericUpdateAction{}, desiredState[9])
	assert.IsType(t, common.GenericUpdateAction{}, desiredState[10])
	assert.IsType(t, model.KeycloakAdminSecret(cr), desiredState[0].(common.GenericUpdateAction).Ref)
	assert.IsType(t, model.PrometheusRule(cr), desiredState[1].(common.GenericUpdateAction).Ref)
	assert.IsType(t, model.ServiceMonitor(cr), desiredState[2].(common.GenericUpdateAction).Ref)
	assert.IsType(t, model.GrafanaDashboard(cr), desiredState[3].(common.GenericUpdateAction).Ref)
	assert.IsType(t, model.DatabaseSecret(cr), desiredState[4].(common.GenericUpdateAction).Ref)
	assert.IsType(t, model.PostgresqlPersistentVolumeClaim(cr), desiredState[5].(common.GenericUpdateAction).Ref)
	assert.IsType(t, model.PostgresqlDeployment(cr), desiredState[6].(common.GenericUpdateAction).Ref)
	assert.IsType(t, model.PostgresqlService(cr), desiredState[7].(common.GenericUpdateAction).Ref)
	assert.IsType(t, model.KeycloakService(cr), desiredState[8].(common.GenericUpdateAction).Ref)
	assert.IsType(t, model.KeycloakDiscoveryService(cr), desiredState[9].(common.GenericUpdateAction).Ref)
	assert.IsType(t, model.KeycloakDeployment(cr), desiredState[10].(common.GenericUpdateAction).Ref)
}

func TestKeycloakReconciler_Test_No_Action_When_Monitoring_Resources_Dont_Exist(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{}
	currentState := common.NewClusterState()

	//Set monitoring resources exist to true
	stateManager := common.GetStateManager()
	stateManager.SetState(getStateFieldName(monitoringv1.PrometheusRuleKind), false)
	stateManager.SetState(getStateFieldName(monitoringv1.ServiceMonitorsKind), false)
	stateManager.SetState(getStateFieldName(integreatlyv1alpha1.GrafanaDashboardKind), false)

	// when
	reconciler := NewKeycloakReconciler()
	desiredState := reconciler.Reconcile(currentState, cr)

	// then
	for _, element := range desiredState {
		assert.IsType(t, common.GenericCreateAction{}, element)
		assert.NotEqual(t, reflect.TypeOf(model.PrometheusRule(cr)), reflect.TypeOf(element.(common.GenericCreateAction).Ref))
		assert.NotEqual(t, reflect.TypeOf(model.GrafanaDashboard(cr)), reflect.TypeOf(element.(common.GenericCreateAction).Ref))
		assert.NotEqual(t, reflect.TypeOf(model.ServiceMonitor(cr)), reflect.TypeOf(element.(common.GenericCreateAction).Ref))
	}
}
