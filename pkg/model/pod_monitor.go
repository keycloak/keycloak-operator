package model

import (
	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func PodMonitor(cr *v1alpha1.Keycloak) *monitoringv1.PodMonitor {
	return &monitoringv1.PodMonitor{
		ObjectMeta: v12.ObjectMeta{
			Name:      PodMonitorName,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"monitoring-key": MonitoringKey,
			},
		},
		Spec: monitoringv1.PodMonitorSpec{
			PodMetricsEndpoints: []monitoringv1.PodMetricsEndpoint{{
				Path: "/auth/realms/master/metrics",
				TargetPort: &intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: KeycloakServicePort,
				},
			}},
			Selector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"component": KeycloakDeploymentComponent,
				},
			},
		},
	}
}

func PodMonitorSelector(cr *v1alpha1.Keycloak) client.ObjectKey {
	return client.ObjectKey{
		Name:      PodMonitorName,
		Namespace: cr.Namespace,
	}
}
