package common

import (
	"context"

	v13 "github.com/openshift/api/route/v1"
	"k8s.io/api/extensions/v1beta1"

	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	integreatlyv1alpha1 "github.com/integr8ly/grafana-operator/pkg/apis/integreatly/v1alpha1"
	kc "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/model"
	v12 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// The desired cluster state is defined by a list of actions that have to be run to
// get from the current state to the desired state
type DesiredClusterState []ClusterAction

func (d *DesiredClusterState) AddAction(action ClusterAction) DesiredClusterState {
	if action != nil {
		*d = append(*d, action)
	}
	return *d
}

func (d *DesiredClusterState) AddActions(actions []ClusterAction) DesiredClusterState {
	for _, action := range actions {
		*d = d.AddAction(action)
	}
	return *d
}

type ClusterState struct {
	KeycloakServiceMonitor          *monitoringv1.ServiceMonitor
	KeycloakPrometheusRule          *monitoringv1.PrometheusRule
	KeycloakGrafanaDashboard        *integreatlyv1alpha1.GrafanaDashboard
	DatabaseSecret                  *v1.Secret
	PostgresqlPersistentVolumeClaim *v1.PersistentVolumeClaim
	PostgresqlService               *v1.Service
	PostgresqlDeployment            *v12.Deployment
	KeycloakService                 *v1.Service
	KeycloakDiscoveryService        *v1.Service
	KeycloakDeployment              *v12.StatefulSet
	KeycloakAdminSecret             *v1.Secret
	KeycloakIngress                 *v1beta1.Ingress
	KeycloakRoute                   *v13.Route
	PostgresqlServiceEndpoints      *v1.Endpoints
}

func (i *ClusterState) Read(context context.Context, cr *kc.Keycloak, controllerClient client.Client) error { //nolint
	stateManager := GetStateManager()
	routeKindExists, keyExists := stateManager.GetState(RouteKind).(bool)

	err := i.readKeycloakAdminSecretCurrentState(context, cr, controllerClient)
	if err != nil {
		return err
	}

	err = i.readKeycloakServiceMonitorCurrentState(context, cr, controllerClient)
	if err != nil {
		return err
	}

	err = i.readKeycloakPrometheusRuleCurrentState(context, cr, controllerClient)
	if err != nil {
		return err
	}

	err = i.readKeycloakGrafanaDashboardCurrentState(context, cr, controllerClient)
	if err != nil {
		return err
	}

	err = i.readDatabaseSecretCurrentState(context, cr, controllerClient)
	if err != nil {
		return err
	}

	err = i.readPostgresqlPersistentVolumeClaimCurrentState(context, cr, controllerClient)
	if err != nil {
		return err
	}

	err = i.readPostgresqlDeploymentCurrentState(context, cr, controllerClient)
	if err != nil {
		return err
	}

	err = i.readPostgresqlServiceCurrentState(context, cr, controllerClient)
	if err != nil {
		return err
	}

	err = i.readPostgresqlServiceEndpointsCurrentState(context, cr, controllerClient)
	if err != nil {
		return err
	}

	err = i.readKeycloakServiceCurrentState(context, cr, controllerClient)
	if err != nil {
		return err
	}

	err = i.readKeycloakDiscoveryServiceCurrentState(context, cr, controllerClient)
	if err != nil {
		return err
	}

	err = i.readKeycloakOrRHSSODeploymentCurrentState(context, cr, controllerClient)
	if err != nil {
		return err
	}

	if keyExists && routeKindExists {
		err = i.readKeycloakRouteCurrentState(context, cr, controllerClient)
		if err != nil {
			return err
		}
	} else {
		err = i.readKeycloakIngressCurrentState(context, cr, controllerClient)
		if err != nil {
			return err
		}
	}

	// Read other things
	return nil
}

func (i *ClusterState) readKeycloakAdminSecretCurrentState(context context.Context, cr *kc.Keycloak, controllerClient client.Client) error {
	keycloakAdminSecret := model.KeycloakAdminSecret(cr)
	keycloakAdminSecretSelector := model.KeycloakAdminSecretSelector(cr)

	err := controllerClient.Get(context, keycloakAdminSecretSelector, keycloakAdminSecret)

	if err != nil {
		// If the resource type doesn't exist on the cluster or does exist but is not found
		if meta.IsNoMatchError(err) || apiErrors.IsNotFound(err) {
			i.KeycloakAdminSecret = nil
		} else {
			return err
		}
	} else {
		i.KeycloakAdminSecret = keycloakAdminSecret.DeepCopy()
		cr.UpdateStatusSecondaryResources(i.KeycloakAdminSecret.Kind, i.KeycloakAdminSecret.Name)
	}
	return nil
}

func (i *ClusterState) readPostgresqlPersistentVolumeClaimCurrentState(context context.Context, cr *kc.Keycloak, controllerClient client.Client) error {
	postgresqlPersistentVolumeClaim := model.PostgresqlPersistentVolumeClaim(cr)
	postgresqlPersistentVolumeClaimSelector := model.PostgresqlPersistentVolumeClaimSelector(cr)

	err := controllerClient.Get(context, postgresqlPersistentVolumeClaimSelector, postgresqlPersistentVolumeClaim)
	if err != nil {
		if !apiErrors.IsNotFound(err) {
			return err
		}
	} else {
		i.PostgresqlPersistentVolumeClaim = postgresqlPersistentVolumeClaim.DeepCopy()
		cr.UpdateStatusSecondaryResources(i.PostgresqlPersistentVolumeClaim.Kind, i.PostgresqlPersistentVolumeClaim.Name)
	}
	return nil
}

func (i *ClusterState) readPostgresqlServiceCurrentState(context context.Context, cr *kc.Keycloak, controllerClient client.Client) error {
	postgresqlService := model.PostgresqlService(cr)
	postgresqlServiceSelector := model.PostgresqlServiceSelector(cr)

	err := controllerClient.Get(context, postgresqlServiceSelector, postgresqlService)
	if err != nil {
		if !apiErrors.IsNotFound(err) {
			return err
		}
	} else {
		i.PostgresqlService = postgresqlService.DeepCopy()
		cr.UpdateStatusSecondaryResources(i.PostgresqlService.Kind, i.PostgresqlService.Name)
	}
	return nil
}

func (i *ClusterState) readPostgresqlServiceEndpointsCurrentState(context context.Context, cr *kc.Keycloak, controllerClient client.Client) error {
	postgresqlServiceEndpoints := model.PostgresqlServiceEndpoints(cr)
	postgresqlServiceEndpointsSelector := model.PostgresqlServiceEndpointsSelector(cr)

	err := controllerClient.Get(context, postgresqlServiceEndpointsSelector, postgresqlServiceEndpoints)
	if err != nil {
		if !apiErrors.IsNotFound(err) {
			return err
		}
	} else {
		i.PostgresqlServiceEndpoints = postgresqlServiceEndpoints.DeepCopy()
		if cr.Spec.ExternalDatabase.Enabled {
			cr.UpdateStatusSecondaryResources(i.PostgresqlService.Kind, i.PostgresqlService.Name)
		}
	}
	return nil
}

func (i *ClusterState) readPostgresqlDeploymentCurrentState(context context.Context, cr *kc.Keycloak, controllerClient client.Client) error {
	postgresqlDeployment := model.PostgresqlDeployment(cr)
	postgresqlDeploymentSelector := model.PostgresqlDeploymentSelector(cr)

	err := controllerClient.Get(context, postgresqlDeploymentSelector, postgresqlDeployment)
	if err != nil {
		if !apiErrors.IsNotFound(err) {
			return err
		}
	} else {
		i.PostgresqlDeployment = postgresqlDeployment.DeepCopy()
		cr.UpdateStatusSecondaryResources(i.PostgresqlDeployment.Kind, i.PostgresqlDeployment.Name)
	}
	return nil
}

func (i *ClusterState) readKeycloakServiceCurrentState(context context.Context, cr *kc.Keycloak, controllerClient client.Client) error {
	keycloakService := model.KeycloakService(cr)
	keycloakServiceSelector := model.KeycloakServiceSelector(cr)

	err := controllerClient.Get(context, keycloakServiceSelector, keycloakService)
	if err != nil {
		if !apiErrors.IsNotFound(err) {
			return err
		}
	} else {
		i.KeycloakService = keycloakService.DeepCopy()
		cr.UpdateStatusSecondaryResources(i.KeycloakService.Kind, i.KeycloakService.Name)
	}
	return nil
}

/*
 *
 * Monitoring Resources
 *
 */

// Keycloak Service Monitor. Resource type provided by Prometheus operator
func (i *ClusterState) readKeycloakServiceMonitorCurrentState(context context.Context, cr *kc.Keycloak, controllerClient client.Client) error {
	keycloakServiceMonitor := model.ServiceMonitor(cr)
	keycloakServiceMonitorSelector := model.ServiceMonitorSelector(cr)

	err := controllerClient.Get(context, keycloakServiceMonitorSelector, keycloakServiceMonitor)

	if err != nil {
		// If the resource type doesn't exist on the cluster or does exist but is not found
		if meta.IsNoMatchError(err) || apiErrors.IsNotFound(err) {
			i.KeycloakServiceMonitor = nil
		} else {
			return err
		}
	} else {
		i.KeycloakServiceMonitor = keycloakServiceMonitor.DeepCopy()
		cr.UpdateStatusSecondaryResources(i.KeycloakServiceMonitor.Kind, i.KeycloakServiceMonitor.Name)
	}
	return nil
}

// Keycloak Prometheus Rule. Resource type provided by Prometheus operator
func (i *ClusterState) readKeycloakPrometheusRuleCurrentState(context context.Context, cr *kc.Keycloak, controllerClient client.Client) error {
	keycloakPrometheusRule := model.PrometheusRule(cr)
	keycloakPrometheusRuleSelector := model.PrometheusRuleSelector(cr)

	err := controllerClient.Get(context, keycloakPrometheusRuleSelector, keycloakPrometheusRule)

	if err != nil {
		// If the resource type doesn't exist on the cluster or does exist but is not found
		if meta.IsNoMatchError(err) || apiErrors.IsNotFound(err) {
			i.KeycloakPrometheusRule = nil
		} else {
			return err
		}
	} else {
		i.KeycloakPrometheusRule = keycloakPrometheusRule.DeepCopy()
		cr.UpdateStatusSecondaryResources(i.KeycloakPrometheusRule.Kind, i.KeycloakPrometheusRule.Name)
	}
	return nil
}

// Keycloak Grafana Dashboard. Resource type provided by Grafana operator
func (i *ClusterState) readKeycloakGrafanaDashboardCurrentState(context context.Context, cr *kc.Keycloak, controllerClient client.Client) error {
	keycloakGrafanaDashboard := model.GrafanaDashboard(cr)
	keycloakGrafanaDashboardSelector := model.GrafanaDashboardSelector(cr)

	err := controllerClient.Get(context, keycloakGrafanaDashboardSelector, keycloakGrafanaDashboard)

	if err != nil {
		// If the resource type doesn't exist on the cluster or does exist but is not found
		if meta.IsNoMatchError(err) || apiErrors.IsNotFound(err) {
			i.KeycloakGrafanaDashboard = nil
		} else {
			return err
		}
	} else {
		i.KeycloakGrafanaDashboard = keycloakGrafanaDashboard.DeepCopy()
		cr.UpdateStatusSecondaryResources(i.KeycloakGrafanaDashboard.Kind, i.KeycloakGrafanaDashboard.Name)
	}
	return nil
}

func (i *ClusterState) readDatabaseSecretCurrentState(context context.Context, cr *kc.Keycloak, controllerClient client.Client) error {
	databaseSecret := model.DatabaseSecret(cr)
	databaseSecretSelector := model.DatabaseSecretSelector(cr)

	err := controllerClient.Get(context, databaseSecretSelector, databaseSecret)

	if err != nil {
		// If the resource type doesn't exist on the cluster or does exist but is not found
		if meta.IsNoMatchError(err) || apiErrors.IsNotFound(err) {
			i.DatabaseSecret = nil
		} else {
			return err
		}
	} else {
		i.DatabaseSecret = databaseSecret.DeepCopy()
		cr.UpdateStatusSecondaryResources(i.DatabaseSecret.Kind, i.DatabaseSecret.Name)
	}
	return nil
}

func (i *ClusterState) readKeycloakOrRHSSODeploymentCurrentState(context context.Context, cr *kc.Keycloak, controllerClient client.Client) error {
	isRHSSO := cr.Spec.Profile == RHSSOProfile

	deployment := model.KeycloakDeployment(cr)
	selector := model.KeycloakDeploymentSelector(cr)
	if isRHSSO {
		deployment = model.RHSSODeployment(cr)
		selector = model.RHSSODeploymentSelector(cr)
	}

	err := controllerClient.Get(context, selector, deployment)
	if err != nil {
		if !apiErrors.IsNotFound(err) {
			return err
		}
	} else {
		i.KeycloakDeployment = deployment.DeepCopy()
		cr.UpdateStatusSecondaryResources(i.KeycloakDeployment.Kind, i.KeycloakDeployment.Name)
	}
	return nil
}

func (i *ClusterState) readKeycloakDiscoveryServiceCurrentState(context context.Context, cr *kc.Keycloak, controllerClient client.Client) error {
	keycloakDiscoveryService := model.KeycloakDiscoveryService(cr)
	keycloakDiscoveryServiceSelector := model.KeycloakDiscoveryServiceSelector(cr)

	err := controllerClient.Get(context, keycloakDiscoveryServiceSelector, keycloakDiscoveryService)
	if err != nil {
		if !apiErrors.IsNotFound(err) {
			return err
		}
	} else {
		i.KeycloakDiscoveryService = keycloakDiscoveryService.DeepCopy()
		cr.UpdateStatusSecondaryResources(i.KeycloakDiscoveryService.Kind, i.KeycloakDiscoveryService.Name)
	}
	return nil
}

func (i *ClusterState) readKeycloakRouteCurrentState(context context.Context, cr *kc.Keycloak, controllerClient client.Client) error {
	keycloakRoute := model.KeycloakRoute(cr)
	keycloakRouteSelector := model.KeycloakRouteSelector(cr)

	err := controllerClient.Get(context, keycloakRouteSelector, keycloakRoute)
	if err != nil {
		if !apiErrors.IsNotFound(err) {
			return err
		}
	} else {
		i.KeycloakRoute = keycloakRoute.DeepCopy()
		cr.UpdateStatusSecondaryResources(i.KeycloakRoute.Kind, i.KeycloakRoute.Name)
	}
	return nil
}

func (i *ClusterState) readKeycloakIngressCurrentState(context context.Context, cr *kc.Keycloak, controllerClient client.Client) error {
	keycloakIngress := model.KeycloakIngress(cr)
	keycloakIngressSelector := model.KeycloakIngressSelector(cr)

	err := controllerClient.Get(context, keycloakIngressSelector, keycloakIngress)
	if err != nil {
		if !apiErrors.IsNotFound(err) {
			return err
		}
	} else {
		i.KeycloakIngress = keycloakIngress.DeepCopy()
		cr.UpdateStatusSecondaryResources(i.KeycloakIngress.Kind, i.KeycloakIngress.Name)
	}
	return nil
}

func (i *ClusterState) IsResourcesReady() (bool, error) {
	// Check keycloak statefulset is ready
	keycloakDeploymentReady, _ := IsStatefulSetReady(i.KeycloakDeployment)
	// Default Route ready to true in case we are running on native Kubernetes
	keycloakRouteReady := true
	// Check keycloak postgres deployment is ready
	postgresqlDeploymentReady, err := IsDeploymentReady(i.PostgresqlDeployment)
	if err != nil {
		return false, err
	}

	// If running on OpenShift, check the Route is ready
	stateManager := GetStateManager()
	openshift, keyExists := stateManager.GetState(RouteKind).(bool)
	if keyExists && openshift {
		keycloakRouteReady = IsRouteReady(i.KeycloakRoute)
	}

	return keycloakDeploymentReady && postgresqlDeploymentReady && keycloakRouteReady, nil
}
