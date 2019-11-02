package keycloak

import (
	"reflect"
	"testing"

	"k8s.io/api/extensions/v1beta1"

	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	integreatlyv1alpha1 "github.com/integr8ly/grafana-operator/pkg/apis/integreatly/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/common"
	"github.com/keycloak/keycloak-operator/pkg/model"
	"github.com/stretchr/testify/assert"
	v13 "k8s.io/api/apps/v1"
)

func TestKeycloakReconciler_Test_Creating_All(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{}
	cr.Spec.ExternalAccess = v1alpha1.KeycloakExternalAccess{
		Enabled: true,
	}

	currentState := common.NewClusterState()

	//Set monitoring resources exist to true
	stateManager := common.GetStateManager()
	defer stateManager.Clear()
	stateManager.SetState(common.GetStateFieldName(ControllerName, monitoringv1.PrometheusRuleKind), true)
	stateManager.SetState(common.GetStateFieldName(ControllerName, monitoringv1.ServiceMonitorsKind), true)
	stateManager.SetState(common.GetStateFieldName(ControllerName, integreatlyv1alpha1.GrafanaDashboardKind), true)
	stateManager.SetState(common.RouteKind, true)

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
	//   10) Keycloak Route
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
	assert.IsType(t, common.GenericCreateAction{}, desiredState[11])
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
	assert.IsType(t, model.KeycloakRoute(cr), desiredState[11].(common.GenericCreateAction).Ref)
}

func TestKeycloakReconciler_Test_Creating_RHSSO(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			ExternalAccess: v1alpha1.KeycloakExternalAccess{
				Enabled: true,
			},
			Profile: common.RHSSOProfile,
		},
	}
	currentState := common.NewClusterState()

	// when
	reconciler := NewKeycloakReconciler()
	desiredState := reconciler.Reconcile(currentState, cr)

	// then
	var allCreateActions = true
	var deployment *v13.StatefulSet
	var ingress *v1beta1.Ingress
	for _, v := range desiredState {
		if reflect.TypeOf(v) != reflect.TypeOf(common.GenericCreateAction{}) {
			allCreateActions = false
		}
		if reflect.TypeOf(v.(common.GenericCreateAction).Ref) == reflect.TypeOf(model.RHSSODeployment(cr)) {
			deployment = v.(common.GenericCreateAction).Ref.(*v13.StatefulSet)
		}
		if reflect.TypeOf(v.(common.GenericCreateAction).Ref) == reflect.TypeOf(model.KeycloakIngress(cr)) {
			ingress = v.(common.GenericCreateAction).Ref.(*v1beta1.Ingress)
		}
	}
	assert.True(t, allCreateActions)
	assert.NotNil(t, deployment)
	assert.NotNil(t, ingress)
	assert.Equal(t, model.RHSSODeployment(cr), deployment)
}

func TestKeycloakReconciler_Test_Updating_RHSSO(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			ExternalAccess: v1alpha1.KeycloakExternalAccess{
				Enabled: true,
			},
			Profile: common.RHSSOProfile,
		},
	}
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
		KeycloakDeployment:              model.RHSSODeployment(cr),
		KeycloakAdminSecret:             model.KeycloakAdminSecret(cr),
		KeycloakIngress:                 model.KeycloakIngress(cr),
	}

	// when
	reconciler := NewKeycloakReconciler()
	desiredState := reconciler.Reconcile(currentState, cr)

	// then
	var allUpdateActions = true
	var deployment *v13.StatefulSet
	for _, v := range desiredState {
		if reflect.TypeOf(v) != reflect.TypeOf(common.GenericUpdateAction{}) {
			allUpdateActions = false
		}
		if reflect.TypeOf(v.(common.GenericUpdateAction).Ref) == reflect.TypeOf(model.RHSSODeployment(cr)) {
			deployment = v.(common.GenericUpdateAction).Ref.(*v13.StatefulSet)
		}
	}
	assert.True(t, allUpdateActions)
	assert.NotNil(t, deployment)
	assert.Equal(t, model.RHSSODeployment(cr), deployment)
}

func TestKeycloakReconciler_Test_Updating_All(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{}
	cr.Spec.ExternalAccess = v1alpha1.KeycloakExternalAccess{
		Enabled: true,
	}

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
		KeycloakRoute:                   model.KeycloakRoute(cr),
	}

	//Set monitoring resources exist to true
	stateManager := common.GetStateManager()
	stateManager.SetState(common.GetStateFieldName(ControllerName, monitoringv1.PrometheusRuleKind), true)
	stateManager.SetState(common.GetStateFieldName(ControllerName, monitoringv1.ServiceMonitorsKind), true)
	stateManager.SetState(common.GetStateFieldName(ControllerName, integreatlyv1alpha1.GrafanaDashboardKind), true)
	stateManager.SetState(common.RouteKind, true)
	defer stateManager.Clear()

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
	//   10) Keycloak Route
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
	assert.IsType(t, common.GenericUpdateAction{}, desiredState[11])
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
	assert.IsType(t, model.KeycloakRoute(cr), desiredState[11].(common.GenericUpdateAction).Ref)
}

func TestKeycloakReconciler_Test_No_Action_When_Monitoring_Resources_Dont_Exist(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{}
	currentState := common.NewClusterState()

	//Set monitoring resources exist to true
	stateManager := common.GetStateManager()
	stateManager.SetState(common.GetStateFieldName(ControllerName, monitoringv1.PrometheusRuleKind), false)
	stateManager.SetState(common.GetStateFieldName(ControllerName, monitoringv1.ServiceMonitorsKind), false)
	stateManager.SetState(common.GetStateFieldName(ControllerName, integreatlyv1alpha1.GrafanaDashboardKind), false)
	defer stateManager.Clear()

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

func TestKeycloakReconciler_Test_Recreate_Credentials_When_Missig(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{}
	secret := model.KeycloakAdminSecret(cr)

	// when
	secret.Data[model.AdminUsernameProperty] = nil
	secret.Data[model.AdminPasswordProperty] = nil
	secret = model.KeycloakAdminSecretReconciled(cr, secret)

	// then
	assert.NotEmpty(t, secret.Data[model.AdminUsernameProperty])
	assert.NotEmpty(t, secret.Data[model.AdminPasswordProperty])
}

func TestKeycloakReconciler_Test_Recreate_Does_Not_Update_Existing_Credentials(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{}
	secret := model.KeycloakAdminSecret(cr)

	// when
	username := secret.Data[model.AdminUsernameProperty]
	password := secret.Data[model.AdminPasswordProperty]
	secret = model.KeycloakAdminSecretReconciled(cr, secret)

	// then
	assert.Equal(t, username, secret.Data[model.AdminUsernameProperty])
	assert.Equal(t, password, secret.Data[model.AdminPasswordProperty])
}
