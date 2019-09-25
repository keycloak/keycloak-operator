package keycloak

import (
	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ServiceMonitor(cr *v1alpha1.Keycloak) *monitoringv1.ServiceMonitor {
	return &monitoringv1.ServiceMonitor{
		ObjectMeta: v12.ObjectMeta{
			Name:      "keycloak-monitoring",
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"monitoring-key": "middleware",
			},
		},
		Spec: monitoringv1.ServiceMonitorSpec{
			Endpoints: []monitoringv1.Endpoint{{
				Path:   "/auth/realms/master/metrics",
				Port:   "sso",
				Scheme: "https",
				TLSConfig: &monitoringv1.TLSConfig{
					InsecureSkipVerify: true,
				},
			}},
			Selector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"application": "sso",
				},
			},
		},
	}
}
