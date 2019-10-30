package model

import (
	integreatlyv1alpha1 "github.com/integr8ly/grafana-operator/pkg/apis/integreatly/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GrafanaDashboard(cr *v1alpha1.Keycloak) *integreatlyv1alpha1.GrafanaDashboard {
	return &integreatlyv1alpha1.GrafanaDashboard{
		ObjectMeta: v12.ObjectMeta{
			Name:      ApplicationName,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"monitoring-key": MonitoringKey,
			},
		},
		Spec: integreatlyv1alpha1.GrafanaDashboardSpec{
			Json: GrafanaDashboardJSON,
			Name: "keycloak.json",
			Plugins: []integreatlyv1alpha1.GrafanaPlugin{
				{
					Name:    "grafana-piechart-panel",
					Version: "1.3.9",
				},
			},
		},
	}
}

func GrafanaDashboardReconciled(cr *v1alpha1.Keycloak, currentState *integreatlyv1alpha1.GrafanaDashboard) *integreatlyv1alpha1.GrafanaDashboard {
	reconciled := currentState.DeepCopy()
	reconciled.Spec.Json = GrafanaDashboardJSON
	reconciled.Spec.Name = "keycloak.json"
	reconciled.Spec.Plugins = []integreatlyv1alpha1.GrafanaPlugin{
		{
			Name:    "grafana-piechart-panel",
			Version: "1.3.9",
		},
	}
	return reconciled
}

func GrafanaDashboardSelector(cr *v1alpha1.Keycloak) client.ObjectKey {
	return client.ObjectKey{
		Name:      ApplicationName,
		Namespace: cr.Namespace,
	}
}
