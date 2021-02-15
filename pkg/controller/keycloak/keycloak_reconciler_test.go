package keycloak

import (
	"reflect"
	"strconv"
	"testing"

	"k8s.io/apimachinery/pkg/api/resource"

	v1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	grafanav1alpha1 "github.com/integr8ly/grafana-operator/v3/pkg/apis/integreatly/v1alpha1"
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
	stateManager.SetState(common.GetStateFieldName(ControllerName, grafanav1alpha1.GrafanaDashboardKind), true)
	stateManager.SetState(common.RouteKind, true)
	stateManager.SetState(common.OpenShiftAPIServerKind, true)

	// when
	reconciler := NewKeycloakReconciler()
	desiredState := reconciler.Reconcile(currentState, cr)

	// then
	// Expectation:
	//    0) Keycloak Admin Secret
	//    1) Prometheus Rule
	//    2) Service Monitor
	//    3) Grafana Dashboard
	//    4) Database secret
	//    5) Postgresql Persistent Volume Claim
	//    6) Postgresql Deployment
	//    7) Postgresql Service
	//    8) Keycloak Service
	//    9) Keycloak Discovery Service
	//    10) Keycloak Monitoring Service
	//    11) Keycloak Probe ConfigMap
	//    12) Keycloak StatefulSets
	//        Note, that's no MetricsRoute as it needs an established hostname from the root route
	//    13) Keycloak Route
	assert.Equal(t, len(desiredState), 14)
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
	assert.IsType(t, common.GenericCreateAction{}, desiredState[12])
	assert.IsType(t, common.GenericCreateAction{}, desiredState[13])
	assert.IsType(t, model.KeycloakAdminSecret(cr), desiredState[0].(common.GenericCreateAction).Ref)
	assert.IsType(t, model.PrometheusRule(cr), desiredState[1].(common.GenericCreateAction).Ref)
	assert.IsType(t, model.ServiceMonitor(cr), desiredState[2].(common.GenericCreateAction).Ref)
	assert.IsType(t, model.GrafanaDashboard(cr), desiredState[3].(common.GenericCreateAction).Ref)
	assert.IsType(t, model.DatabaseSecret(cr), desiredState[4].(common.GenericCreateAction).Ref)
	assert.IsType(t, model.PostgresqlPersistentVolumeClaim(cr), desiredState[5].(common.GenericCreateAction).Ref)
	assert.IsType(t, model.PostgresqlDeployment(cr, true), desiredState[6].(common.GenericCreateAction).Ref)
	assert.IsType(t, model.PostgresqlService(cr, model.DatabaseSecret(cr), false), desiredState[7].(common.GenericCreateAction).Ref)
	assert.IsType(t, model.KeycloakService(cr), desiredState[8].(common.GenericCreateAction).Ref)
	assert.IsType(t, model.KeycloakDiscoveryService(cr), desiredState[9].(common.GenericCreateAction).Ref)
	assert.IsType(t, model.KeycloakMonitoringService(cr), desiredState[10].(common.GenericCreateAction).Ref)
	assert.IsType(t, model.KeycloakProbes(cr), desiredState[11].(common.GenericCreateAction).Ref)
	assert.IsType(t, model.KeycloakDeployment(cr, model.DatabaseSecret(cr)), desiredState[12].(common.GenericCreateAction).Ref)
	assert.IsType(t, model.KeycloakRoute(cr), desiredState[13].(common.GenericCreateAction).Ref)
}

func TestKeycloakReconciler_Test_Creating_RHSSO(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			ExternalAccess: v1alpha1.KeycloakExternalAccess{
				Enabled: true,
			},
			Profile: model.RHSSOProfile,
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
		if reflect.TypeOf(v.(common.GenericCreateAction).Ref) == reflect.TypeOf(model.RHSSODeployment(cr, model.DatabaseSecret(cr))) {
			deployment = v.(common.GenericCreateAction).Ref.(*v13.StatefulSet)
		}
		if reflect.TypeOf(v.(common.GenericCreateAction).Ref) == reflect.TypeOf(model.KeycloakIngress(cr)) {
			ingress = v.(common.GenericCreateAction).Ref.(*v1beta1.Ingress)
		}
	}
	assert.True(t, allCreateActions)
	assert.NotNil(t, deployment)
	assert.NotNil(t, ingress)
	assert.Equal(t, model.RHSSODeployment(cr, nil), deployment)
}

func TestKeycloakReconciler_Test_Updating_RHSSO(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			ExternalAccess: v1alpha1.KeycloakExternalAccess{
				Enabled: true,
			},
			Profile:   model.RHSSOProfile,
			Instances: 1,
		},
	}
	currentState := &common.ClusterState{
		KeycloakServiceMonitor:          model.ServiceMonitor(cr),
		KeycloakPrometheusRule:          model.PrometheusRule(cr),
		KeycloakGrafanaDashboard:        model.GrafanaDashboard(cr),
		DatabaseSecret:                  model.DatabaseSecret(cr),
		PostgresqlPersistentVolumeClaim: model.PostgresqlPersistentVolumeClaim(cr),
		PostgresqlService:               model.PostgresqlService(cr, model.DatabaseSecret(cr), false),
		PostgresqlDeployment:            model.PostgresqlDeployment(cr, true),
		KeycloakService:                 model.KeycloakService(cr),
		KeycloakDiscoveryService:        model.KeycloakDiscoveryService(cr),
		KeycloakDeployment:              model.RHSSODeployment(cr, model.DatabaseSecret(cr)),
		KeycloakAdminSecret:             model.KeycloakAdminSecret(cr),
		KeycloakIngress:                 model.KeycloakIngress(cr),
		KeycloakProbes:                  model.KeycloakProbes(cr),
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
		if reflect.TypeOf(v.(common.GenericUpdateAction).Ref) == reflect.TypeOf(model.RHSSODeployment(cr, model.DatabaseSecret(cr))) {
			deployment = v.(common.GenericUpdateAction).Ref.(*v13.StatefulSet)
		}
	}
	assert.True(t, allUpdateActions)
	assert.NotNil(t, deployment)
	assert.Equal(t, model.RHSSODeployment(cr, model.DatabaseSecret(cr)), deployment)
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
		PostgresqlService:               model.PostgresqlService(cr, model.DatabaseSecret(cr), false),
		PostgresqlDeployment:            model.PostgresqlDeployment(cr, true),
		KeycloakService:                 model.KeycloakService(cr),
		KeycloakDiscoveryService:        model.KeycloakDiscoveryService(cr),
		KeycloakMonitoringService:       model.KeycloakMonitoringService(cr),
		KeycloakDeployment:              model.KeycloakDeployment(cr, model.DatabaseSecret(cr)),
		KeycloakAdminSecret:             model.KeycloakAdminSecret(cr),
		KeycloakRoute:                   model.KeycloakRoute(cr),
		KeycloakMetricsRoute:            model.KeycloakMetricsRoute(cr, model.KeycloakRoute(cr)),
		KeycloakProbes:                  model.KeycloakProbes(cr),
	}

	//Set monitoring resources exist to true
	stateManager := common.GetStateManager()
	stateManager.SetState(common.GetStateFieldName(ControllerName, monitoringv1.PrometheusRuleKind), true)
	stateManager.SetState(common.GetStateFieldName(ControllerName, monitoringv1.ServiceMonitorsKind), true)
	stateManager.SetState(common.GetStateFieldName(ControllerName, grafanav1alpha1.GrafanaDashboardKind), true)
	stateManager.SetState(common.RouteKind, true)
	stateManager.SetState(common.OpenShiftAPIServerKind, true)
	defer stateManager.Clear()

	// when
	reconciler := NewKeycloakReconciler()
	desiredState := reconciler.Reconcile(currentState, cr)

	// then
	// Expectation:
	//    0) Keycloak Admin Secret
	//    1) Prometheus Rule
	//    2) Service Monitor
	//    3) Grafana Dashboard
	//    4) Database secret
	//    5) Postgresql Persistent Volume Claim
	//    6) Postgresql Deployment
	//    7) Postgresql Service
	//    8) Keycloak Service
	//    9) Keycloak Discovery Service
	//    10) Keycloak Monitoring Service
	//    11) Keycloak StatefulSets
	//    12) Keycloak Route
	//    13) Keycloak Metrics Route
	assert.Equal(t, 14, len(desiredState))
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
	assert.IsType(t, common.GenericUpdateAction{}, desiredState[12])
	assert.IsType(t, model.KeycloakAdminSecret(cr), desiredState[0].(common.GenericUpdateAction).Ref)
	assert.IsType(t, model.PrometheusRule(cr), desiredState[1].(common.GenericUpdateAction).Ref)
	assert.IsType(t, model.ServiceMonitor(cr), desiredState[2].(common.GenericUpdateAction).Ref)
	assert.IsType(t, model.GrafanaDashboard(cr), desiredState[3].(common.GenericUpdateAction).Ref)
	assert.IsType(t, model.DatabaseSecret(cr), desiredState[4].(common.GenericUpdateAction).Ref)
	assert.IsType(t, model.PostgresqlPersistentVolumeClaim(cr), desiredState[5].(common.GenericUpdateAction).Ref)
	assert.IsType(t, model.PostgresqlDeployment(cr, true), desiredState[6].(common.GenericUpdateAction).Ref)
	assert.IsType(t, model.PostgresqlService(cr, model.DatabaseSecret(cr), false), desiredState[7].(common.GenericUpdateAction).Ref)
	assert.IsType(t, model.KeycloakService(cr), desiredState[8].(common.GenericUpdateAction).Ref)
	assert.IsType(t, model.KeycloakDiscoveryService(cr), desiredState[9].(common.GenericUpdateAction).Ref)
	assert.IsType(t, model.KeycloakMonitoringService(cr), desiredState[10].(common.GenericUpdateAction).Ref)
	assert.IsType(t, model.KeycloakDeployment(cr, model.DatabaseSecret(cr)), desiredState[11].(common.GenericUpdateAction).Ref)
	assert.IsType(t, model.KeycloakMetricsRoute(cr, model.KeycloakRoute(cr)), desiredState[12].(common.GenericUpdateAction).Ref)
}

func TestKeycloakReconciler_Test_No_Action_When_Monitoring_Resources_Dont_Exist(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{}
	currentState := common.NewClusterState()

	//Set monitoring resources exist to true
	stateManager := common.GetStateManager()
	stateManager.SetState(common.GetStateFieldName(ControllerName, monitoringv1.PrometheusRuleKind), false)
	stateManager.SetState(common.GetStateFieldName(ControllerName, monitoringv1.ServiceMonitorsKind), false)
	stateManager.SetState(common.GetStateFieldName(ControllerName, grafanav1alpha1.GrafanaDashboardKind), false)
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

func TestKeycloakReconciler_Test_Creating_All_With_External_Database(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{}
	cr.Spec.ExternalDatabase.Enabled = true

	currentState := common.NewClusterState()
	currentState.DatabaseSecret = model.DatabaseSecret(cr)

	// when
	reconciler := NewKeycloakReconciler()
	desiredState := reconciler.Reconcile(currentState, cr)

	// then
	assert.IsType(t, common.GenericCreateAction{}, desiredState[0])
	assert.IsType(t, common.GenericCreateAction{}, desiredState[1])
	assert.IsType(t, common.GenericCreateAction{}, desiredState[2])
	assert.IsType(t, common.GenericCreateAction{}, desiredState[3])
	assert.IsType(t, common.GenericCreateAction{}, desiredState[4])
	assert.IsType(t, common.GenericCreateAction{}, desiredState[5])
	assert.IsType(t, model.DatabaseSecret(cr), desiredState[0].(common.GenericCreateAction).Ref)
	assert.IsType(t, model.PostgresqlService(cr, model.DatabaseSecret(cr), false), desiredState[1].(common.GenericCreateAction).Ref)
	assert.IsType(t, model.KeycloakService(cr), desiredState[2].(common.GenericCreateAction).Ref)
	assert.IsType(t, model.KeycloakDiscoveryService(cr), desiredState[3].(common.GenericCreateAction).Ref)
	assert.IsType(t, model.KeycloakProbes(cr), desiredState[4].(common.GenericCreateAction).Ref)
	assert.IsType(t, model.KeycloakDeployment(cr, model.DatabaseSecret(cr)), desiredState[5].(common.GenericCreateAction).Ref)
}

func TestKeycloakReconciler_Test_Updating_External_Database(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{}
	cr.Spec.ExternalDatabase.Enabled = true

	currentState := common.NewClusterState()
	currentState.PostgresqlServiceEndpoints = model.PostgresqlServiceEndpoints(cr)
	currentState.DatabaseSecret = model.DatabaseSecret(cr)
	// This conversion is done my K8s. In the tests, we need to fake it.
	currentState.DatabaseSecret.Data = map[string][]byte{
		model.DatabaseSecretExternalAddressProperty: []byte("10.10.10.1"),
		model.DatabaseSecretExternalPortProperty:    []byte("5432"),
	}

	// when
	reconciler := NewKeycloakReconciler()
	desiredState := reconciler.Reconcile(currentState, cr)

	// then
	var endpoints *v1.Endpoints
	for _, v := range desiredState {
		if reflect.TypeOf(v) == reflect.TypeOf(common.GenericUpdateAction{}) {
			if reflect.TypeOf(v.(common.GenericUpdateAction).Ref) == reflect.TypeOf(model.PostgresqlServiceEndpoints(cr)) {
				endpoints = v.(common.GenericUpdateAction).Ref.(*v1.Endpoints)
			}
		}
	}
	assert.NotNil(t, endpoints)
	assert.Equal(t, model.PostgresqlServiceEndpointsReconciled(cr, currentState.PostgresqlServiceEndpoints, currentState.DatabaseSecret), endpoints)
}

func TestKeycloakReconciler_Test_Updating_External_Database_URI(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{}
	cr.Spec.ExternalDatabase.Enabled = true

	currentState := common.NewClusterState()
	currentState.PostgresqlServiceEndpoints = model.PostgresqlServiceEndpoints(cr)
	currentState.DatabaseSecret = model.DatabaseSecret(cr)
	// This conversion is done my K8s. In the tests, we need to fake it.
	currentState.DatabaseSecret.Data = map[string][]byte{
		model.DatabaseSecretExternalAddressProperty: []byte("host.example.database.location"),
		model.DatabaseSecretExternalPortProperty:    []byte("5432"),
	}

	// when
	reconciler := NewKeycloakReconciler()
	desiredState := reconciler.Reconcile(currentState, cr)

	// then
	var service *v1.Service
	for _, v := range desiredState {
		if reflect.TypeOf(v) == reflect.TypeOf(common.GenericCreateAction{}) {
			if reflect.TypeOf(v.(common.GenericCreateAction).Ref) == reflect.TypeOf(model.PostgresqlService(cr, currentState.DatabaseSecret, true)) {
				s := v.(common.GenericCreateAction).Ref.(*v1.Service)
				if s.Name == model.PostgresqlServiceName {
					service = s
				}
			}
		}
	}
	assert.NotNil(t, service)
	assert.Equal(t, service.Spec.Type, v1.ServiceTypeExternalName)
	assert.Equal(t, service.Spec.ExternalName, string(currentState.DatabaseSecret.Data[model.DatabaseSecretExternalAddressProperty]))
}

func TestKeycloakReconciler_Test_Updating_External_Database_URI_From_IP_To_ExternalName(t *testing.T) {
	// given
	const (
		oldIP           = "1.2.3.4"
		oldPort         = 1234
		newExternalname = "host.example.database.location"
		newPort         = 5432
	)
	cr := &v1alpha1.Keycloak{}
	cr.Spec.ExternalDatabase.Enabled = true

	currentState := common.NewClusterState()
	currentState.PostgresqlServiceEndpoints = model.PostgresqlServiceEndpoints(cr)
	currentState.DatabaseSecret = model.DatabaseSecret(cr)
	currentState.DatabaseSecret.Data = map[string][]byte{
		model.DatabaseSecretExternalAddressProperty: []byte(oldIP),
		model.DatabaseSecretExternalPortProperty:    []byte(strconv.Itoa(oldPort)),
	}

	currentState.PostgresqlService = model.PostgresqlService(cr, currentState.DatabaseSecret, false)

	// This conversion is done my K8s. In the tests, we need to fake it.
	currentState.DatabaseSecret.Data = map[string][]byte{
		model.DatabaseSecretExternalAddressProperty: []byte(newExternalname),
		model.DatabaseSecretExternalPortProperty:    []byte(strconv.Itoa(newPort)),
	}

	// when
	reconciler := NewKeycloakReconciler()
	desiredState := reconciler.Reconcile(currentState, cr)

	// then
	var service *v1.Service
	for _, v := range desiredState {
		if reflect.TypeOf(v) == reflect.TypeOf(common.GenericUpdateAction{}) {
			if reflect.TypeOf(v.(common.GenericUpdateAction).Ref) == reflect.TypeOf(model.PostgresqlService(cr, currentState.DatabaseSecret, true)) {
				s := v.(common.GenericUpdateAction).Ref.(*v1.Service)
				if s.Name == model.PostgresqlServiceName {
					service = s
				}
			}
		}
	}
	assert.NotNil(t, service)
	assert.Equal(t, service.Spec.Type, v1.ServiceTypeExternalName)
	assert.Equal(t, service.Spec.ExternalName, string(currentState.DatabaseSecret.Data[model.DatabaseSecretExternalAddressProperty]))
	assert.Equal(t, service.Spec.Ports[0].Port, int32(newPort))
}

func TestKeycloakReconciler_Test_Updating_External_Database_From_ExternalName_To_IP(t *testing.T) {
	// given
	const (
		oldExternalname = "host.example.database.location"
		oldPort         = 1234
		newIP           = "1.2.3.4"
		newPort         = 5432
	)
	cr := &v1alpha1.Keycloak{}
	cr.Spec.ExternalDatabase.Enabled = true

	currentState := common.NewClusterState()
	currentState.DatabaseSecret = model.DatabaseSecret(cr)
	currentState.DatabaseSecret.Data = map[string][]byte{
		model.DatabaseSecretExternalAddressProperty: []byte(oldExternalname),
		model.DatabaseSecretExternalPortProperty:    []byte(strconv.Itoa(oldPort)),
	}

	currentState.PostgresqlService = model.PostgresqlService(cr, currentState.DatabaseSecret, false)

	// This conversion is done my K8s. In the tests, we need to fake it.
	currentState.DatabaseSecret.Data = map[string][]byte{
		model.DatabaseSecretExternalAddressProperty: []byte(newIP),
		model.DatabaseSecretExternalPortProperty:    []byte(strconv.Itoa(newPort)),
	}
	currentState.PostgresqlServiceEndpoints = model.PostgresqlServiceEndpoints(cr)

	// when
	reconciler := NewKeycloakReconciler()
	desiredState := reconciler.Reconcile(currentState, cr)

	// then
	var service *v1.Service
	var endpoints *v1.Endpoints
	for _, v := range desiredState {
		if reflect.TypeOf(v) == reflect.TypeOf(common.GenericUpdateAction{}) {
			if reflect.TypeOf(v.(common.GenericUpdateAction).Ref) == reflect.TypeOf(model.PostgresqlService(cr, currentState.DatabaseSecret, true)) {
				s := v.(common.GenericUpdateAction).Ref.(*v1.Service)
				if s.Name == model.PostgresqlServiceName {
					service = s
				}
			}
		}
		if reflect.TypeOf(v) == reflect.TypeOf(common.GenericUpdateAction{}) {
			if reflect.TypeOf(v.(common.GenericUpdateAction).Ref) == reflect.TypeOf(model.PostgresqlServiceEndpoints(cr)) {
				endpoints = v.(common.GenericUpdateAction).Ref.(*v1.Endpoints)
			}
		}
	}
	assert.NotNil(t, service)
	assert.Equal(t, service.Spec.Type, v1.ServiceTypeClusterIP)
	assert.Equal(t, model.PostgresqlServiceEndpointsReconciled(cr, currentState.PostgresqlServiceEndpoints, currentState.DatabaseSecret), endpoints)
	assert.Equal(t, service.Spec.ExternalName, "")
	assert.Equal(t, endpoints.Subsets[0].Ports[0].Port, int32(newPort))
	assert.Equal(t, endpoints.Subsets[0].Addresses[0].IP, newIP)
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

func TestKeycloakReconciler_Test_Should_Create_PDB(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{}
	cr.Spec.PodDisruptionBudget.Enabled = true

	currentState := common.NewClusterState()

	// when
	reconciler := NewKeycloakReconciler()
	desiredState := reconciler.Reconcile(currentState, cr)

	// then
	assert.Equal(t, len(desiredState), 10)
	assert.IsType(t, common.GenericCreateAction{}, desiredState[9])
	assert.IsType(t, model.PodDisruptionBudget(cr), desiredState[9].(common.GenericCreateAction).Ref)
}

func TestKeycloakReconciler_Test_Should_Update_PDB(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{}
	cr.Spec.PodDisruptionBudget.Enabled = true

	currentState := &common.ClusterState{
		PodDisruptionBudget: model.PodDisruptionBudget(cr),
	}

	// when
	reconciler := NewKeycloakReconciler()
	desiredState := reconciler.Reconcile(currentState, cr)

	// then
	assert.Equal(t, len(desiredState), 10)
	assert.IsType(t, common.GenericUpdateAction{}, desiredState[9])
	assert.IsType(t, model.PodDisruptionBudget(cr), desiredState[9].(common.GenericUpdateAction).Ref)
}

func TestKeycloakReconciler_Test_Setting_Resources(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{}
	cr.Spec.ExternalAccess = v1alpha1.KeycloakExternalAccess{
		Enabled: true,
	}

	resource50m := resource.MustParse("50m")
	resource100Mi := resource.MustParse("100Mi")
	resource1900m := resource.MustParse("1900m")
	resource700Mi := resource.MustParse("700Mi")

	var resourceListKeycloak = make(v1.ResourceList)
	resourceListKeycloak[v1.ResourceCPU] = resource1900m
	resourceListKeycloak[v1.ResourceMemory] = resource700Mi

	var resourceListPostgres = make(v1.ResourceList)
	resourceListPostgres[v1.ResourceCPU] = resource50m
	resourceListPostgres[v1.ResourceMemory] = resource100Mi

	cr.Spec.KeycloakDeploymentSpec = v1alpha1.KeycloakDeploymentSpec{
		DeploymentSpec: v1alpha1.DeploymentSpec{
			Resources: v1.ResourceRequirements{
				Requests: resourceListKeycloak,
				Limits:   resourceListKeycloak,
			},
		},
	}
	cr.Spec.PostgresDeploymentSpec = v1alpha1.PostgresqlDeploymentSpec{
		DeploymentSpec: v1alpha1.DeploymentSpec{
			Resources: v1.ResourceRequirements{
				Requests: resourceListPostgres,
				Limits:   resourceListPostgres,
			},
		},
	}

	currentState := common.NewClusterState()

	//Set monitoring resources exist to true
	stateManager := common.GetStateManager()
	defer stateManager.Clear()
	stateManager.SetState(common.GetStateFieldName(ControllerName, monitoringv1.PrometheusRuleKind), true)
	stateManager.SetState(common.GetStateFieldName(ControllerName, monitoringv1.ServiceMonitorsKind), true)
	stateManager.SetState(common.GetStateFieldName(ControllerName, grafanav1alpha1.GrafanaDashboardKind), true)
	stateManager.SetState(common.RouteKind, true)

	// when
	reconciler := NewKeycloakReconciler()
	desiredState := reconciler.Reconcile(currentState, cr)
	// then
	// Expectation:
	//    6) Postgresql Deployment
	//    12) Keycloak StatefulSets
	assert.Equal(t, 14, len(desiredState))
	assert.IsType(t, model.PostgresqlDeployment(cr, false), desiredState[6].(common.GenericCreateAction).Ref)
	assert.IsType(t, model.KeycloakDeployment(cr, model.DatabaseSecret(cr)), desiredState[12].(common.GenericCreateAction).Ref)
	keycloakContainer := desiredState[12].(common.GenericCreateAction).Ref.(*v13.StatefulSet).Spec.Template.Spec.Containers[0]
	assert.Equal(t, &resource700Mi, keycloakContainer.Resources.Requests.Memory(), "Keycloak Deployment: Memory-Requests should be: "+resource700Mi.String()+" but is "+keycloakContainer.Resources.Requests.Memory().String())
	assert.Equal(t, &resource1900m, keycloakContainer.Resources.Requests.Cpu(), "Keycloak Deployment: Cpu-Requests should be: "+resource1900m.String()+" but is "+keycloakContainer.Resources.Requests.Cpu().String())
	assert.Equal(t, &resource700Mi, keycloakContainer.Resources.Limits.Memory(), "Keycloak Deployment: Memory-Limit should be: "+resource700Mi.String()+" but is "+keycloakContainer.Resources.Limits.Memory().String())
	assert.Equal(t, &resource1900m, keycloakContainer.Resources.Limits.Cpu(), "Keycloak Deployment:  Cpu-Limit should be: "+resource1900m.String()+" but is "+keycloakContainer.Resources.Limits.Cpu().String())
	postgresContainer := desiredState[6].(common.GenericCreateAction).Ref.(*v13.Deployment).Spec.Template.Spec.Containers[0]
	assert.Equal(t, &resource100Mi, postgresContainer.Resources.Requests.Memory(), "Postgres Deployment: Memory-Requests should be: "+resource100Mi.String()+" but is: "+postgresContainer.Resources.Requests.Memory().String())
	assert.Equal(t, &resource50m, postgresContainer.Resources.Requests.Cpu(), "Postgres Deployment: Cpu-Requests should be: ", resource50m.String()+"b ut is "+postgresContainer.Resources.Requests.Cpu().String())
	assert.Equal(t, &resource100Mi, postgresContainer.Resources.Limits.Memory(), "Postgres Deployment: Memory-Limits should be: "+resource100Mi.String()+" but is: "+postgresContainer.Resources.Limits.Memory().String())
	assert.Equal(t, &resource50m, postgresContainer.Resources.Limits.Cpu(), "Postgres Deployment: Cpu-Limits should be: ", resource50m.String()+"b ut is "+postgresContainer.Resources.Limits.Cpu().String())
}

func TestKeycloakReconciler_Test_No_Resources_Specified(t *testing.T) {
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
	stateManager.SetState(common.GetStateFieldName(ControllerName, grafanav1alpha1.GrafanaDashboardKind), true)
	stateManager.SetState(common.RouteKind, true)
	stateManager.SetState(common.OpenShiftAPIServerKind, true)

	// when
	reconciler := NewKeycloakReconciler()
	desiredState := reconciler.Reconcile(currentState, cr)
	// then
	// Expectation:
	//    6) Postgresql Deployment
	//    12) Keycloak StatefulSets
	assert.Equal(t, 14, len(desiredState))
	assert.IsType(t, model.PostgresqlDeployment(cr, true), desiredState[6].(common.GenericCreateAction).Ref)
	assert.IsType(t, model.KeycloakDeployment(cr, model.DatabaseSecret(cr)), desiredState[12].(common.GenericCreateAction).Ref)
	keycloakContainer := desiredState[12].(common.GenericCreateAction).Ref.(*v13.StatefulSet).Spec.Template.Spec.Containers[0]
	assert.Equal(t, 0, len(keycloakContainer.Resources.Requests), "Requests-List should be empty")
	assert.Equal(t, 0, len(keycloakContainer.Resources.Limits), "Limits-List should be empty")
	postgresContainer := desiredState[6].(common.GenericCreateAction).Ref.(*v13.Deployment).Spec.Template.Spec.Containers[0]
	assert.Equal(t, len(postgresContainer.Resources.Requests), 0, "Request-List should be empty")
	assert.Equal(t, len(postgresContainer.Resources.Limits), 0, "Limits-List should be empty")
}

func TestKeycloakReconciler_Test_Proxy_Settings(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{}

	currentState := common.NewClusterState()
	reconciler := NewKeycloakReconciler()

	// when
	desiredState := reconciler.Reconcile(currentState, cr)

	// then
	// Expectation:
	//    11) Keycloak StatefulSets
	envs := desiredState[8].(common.GenericCreateAction).Ref.(*v13.StatefulSet).Spec.Template.Spec.Containers[0].Env
	proxySet := false
	for _, val := range envs {
		if val.Name == "PROXY_ADDRESS_FORWARDING" && val.Value == "true" {
			proxySet = true
		}
	}
	assert.True(t, proxySet)
}

func TestKeycloakReconciler_Test_Should_Create_Backup(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{}
	cr.Spec.Migration.Backups.Enabled = true
	backupCr := &v1alpha1.KeycloakBackup{}
	labelSelect := metav1.LabelSelector{
		MatchLabels: cr.Labels,
	}

	currentState := common.NewClusterState()
	currentState.KeycloakBackup = &v1alpha1.KeycloakBackup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      model.MigrateBackupName + "-" + common.BackupTime,
			Namespace: cr.Namespace,
			Labels:    cr.Labels,
		},
		Spec: v1alpha1.KeycloakBackupSpec{
			InstanceSelector: &labelSelect,
		},
	}

	// when
	reconciler := NewKeycloakReconciler()
	desiredState := reconciler.Reconcile(currentState, cr)

	// then
	assert.Equal(t, len(desiredState), 10)
	assert.IsType(t, common.GenericUpdateAction{}, desiredState[9])
	assert.IsType(t, model.KeycloakMigrationOneTimeBackup(backupCr), desiredState[9].(common.GenericUpdateAction).Ref)
}
