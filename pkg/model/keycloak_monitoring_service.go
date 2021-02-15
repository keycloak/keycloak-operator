package model

import (
	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func KeycloakMonitoringService(cr *v1alpha1.Keycloak) *v1.Service {
	return &v1.Service{
		ObjectMeta: v12.ObjectMeta{
			Name:      KeycloakMonitoringServiceName,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app": ApplicationName,
			},
			Annotations: map[string]string{
				"description": "The monitoring service for Prometheus",
			},
		},
		Spec: v1.ServiceSpec{
			Selector: map[string]string{
				"app":       ApplicationName,
				"component": KeycloakDeploymentComponent,
			},
			Ports: []v1.ServicePort{
				{
					Port:       9990,
					TargetPort: intstr.FromInt(9990),
					Name:       KeycloakMonitoringServiceName,
					Protocol:   "TCP",
				},
			},
		},
	}
}

func KeycloakMonitoringServiceSelector(cr *v1alpha1.Keycloak) client.ObjectKey {
	return client.ObjectKey{
		Name:      KeycloakMonitoringServiceName,
		Namespace: cr.Namespace,
	}
}

func KeycloakMonitoringServiceReconciled(cr *v1alpha1.Keycloak, currentState *v1.Service) *v1.Service {
	reconciled := currentState.DeepCopy()
	reconciled.Spec.Ports = []v1.ServicePort{
		{
			Port:       9990,
			TargetPort: intstr.FromInt(9990),
			Name:       KeycloakMonitoringServiceName,
			Protocol:   "TCP",
		},
	}
	return reconciled
}
