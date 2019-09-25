package common

import (
	"context"

	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	integreatlyv1alpha1 "github.com/integr8ly/grafana-operator/pkg/apis/integreatly/v1alpha1"
	kc "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/model"
	v12 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// The desired cluster state is defined by a list of actions that have to be run to
// get from the current state to the desired state
type DesiredClusterState []ClusterAction

func (d DesiredClusterState) AddAction(action ClusterAction) DesiredClusterState {
	if action != nil {
		d = append(d, action)
	}
	return d
}

type ClusterState struct {
	KeycloakService                 *v1.Service
	KeycloakServiceMonitor          *monitoringv1.ServiceMonitor
	KeycloakPrometheusRule          *monitoringv1.PrometheusRule
	KeycloakGrafanaDashboard        *integreatlyv1alpha1.GrafanaDashboard
	PostgresqlPersistentVolumeClaim *v1.PersistentVolumeClaim
	PostgresqlService               *v1.Service
	PostgresqlDeployment            *v12.Deployment
}

func NewClusterState() *ClusterState {
	return &ClusterState{}
}

func (i *ClusterState) Read(context context.Context, cr *kc.Keycloak, controllerClient client.Client) error {
	err := i.readKeycloakServiceMonitorCurrentState(context, cr, controllerClient)
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

	err = i.readKeycloakServiceCurrentState(context, cr, controllerClient)
	if err != nil {
		return err
	}

	// Read other things
	return nil
}

func (i *ClusterState) readPostgresqlPersistentVolumeClaimCurrentState(context context.Context, cr *kc.Keycloak, controllerClient client.Client) error {
	postgresqlPersistentVolumeClaim := model.PostgresqlPersistentVolumeClaim(cr)
	postgresqlPersistentVolumeClaimSelector := model.PostgresqlPersistentVolumeClaimSelector(cr)

	err := controllerClient.Get(context, postgresqlPersistentVolumeClaimSelector, postgresqlPersistentVolumeClaim)
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	} else {
		i.PostgresqlPersistentVolumeClaim = postgresqlPersistentVolumeClaim.DeepCopy()
	}
	return nil
}

func (i *ClusterState) readPostgresqlServiceCurrentState(context context.Context, cr *kc.Keycloak, controllerClient client.Client) error {
	postgresqlService := model.PostgresqlService(cr)
	postgresqlServiceSelector := model.PostgresqlServiceSelector(cr)

	err := controllerClient.Get(context, postgresqlServiceSelector, postgresqlService)
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	} else {
		i.PostgresqlService = postgresqlService.DeepCopy()
	}
	return nil
}

func (i *ClusterState) readPostgresqlDeploymentCurrentState(context context.Context, cr *kc.Keycloak, controllerClient client.Client) error {
	postgresqlDeployment := model.PostgresqlDeployment(cr)
	postgresqlDeploymentSelector := model.PostgresqlDeploymentSelector(cr)

	err := controllerClient.Get(context, postgresqlDeploymentSelector, postgresqlDeployment)
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	} else {
		i.PostgresqlDeployment = postgresqlDeployment.DeepCopy()
	}
	return nil
}

func (i *ClusterState) readKeycloakServiceCurrentState(context context.Context, cr *kc.Keycloak, controllerClient client.Client) error {
	keycloakService := model.KeycloakService(cr)
	keycloakServiceSelector := model.KeycloakServiceSelector(cr)

	err := controllerClient.Get(context, keycloakServiceSelector, keycloakService)
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	} else {
		i.KeycloakService = keycloakService.DeepCopy()
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
		if meta.IsNoMatchError(err) || errors.IsNotFound(err) {
			i.KeycloakServiceMonitor = nil
		} else {
			return err
		}
	} else {
		i.KeycloakServiceMonitor = keycloakServiceMonitor.DeepCopy()
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
		if meta.IsNoMatchError(err) || errors.IsNotFound(err) {
			i.KeycloakPrometheusRule = nil
		} else {
			return err
		}
	} else {
		i.KeycloakPrometheusRule = keycloakPrometheusRule.DeepCopy()
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
		if meta.IsNoMatchError(err) || errors.IsNotFound(err) {
			i.KeycloakGrafanaDashboard = nil
		} else {
			return err
		}
	} else {
		i.KeycloakGrafanaDashboard = keycloakGrafanaDashboard.DeepCopy()
	}
	return nil
}
