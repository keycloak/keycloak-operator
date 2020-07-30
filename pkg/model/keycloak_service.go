package model

import (
	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetServicePortName() string {
  if (KeycloakServicePort == 80 || KeycloakServicePort == 8080) {
    return "http"
  } else {
    return "https"
  }
}

func GetServicePorts() []v1.ServicePort {
	return []v1.ServicePort{
		{
			Port:       KeycloakServicePort,
			TargetPort: intstr.FromInt(KeycloakServicePort),
			Name:       GetServicePortName(),
			Protocol:   "TCP",
		},
	}
}

func GetServiceAnnotations(cr *v1alpha1.Keycloak) map[string]string {
	annotations := map[string]string{
		"description": "The web server's https port.",
		"service.alpha.openshift.io/serving-cert-secret-name": ServingCertSecretName,
	}
	for key, value := range cr.Spec.ServiceAnnotations {
		annotations[key] = value
	}
	return annotations
}

func KeycloakService(cr *v1alpha1.Keycloak) *v1.Service {
	return &v1.Service{
		ObjectMeta: v12.ObjectMeta{
			Name:      ApplicationName,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app": ApplicationName,
			},
			Annotations: GetServiceAnnotations(cr),
		},
		Spec: v1.ServiceSpec{
			Selector: map[string]string{
				"app":       ApplicationName,
				"component": KeycloakDeploymentComponent,
			},
			Ports: GetServicePorts(),
		},
	}
}

func KeycloakServiceSelector(cr *v1alpha1.Keycloak) client.ObjectKey {
	return client.ObjectKey{
		Name:      ApplicationName,
		Namespace: cr.Namespace,
	}
}

func KeycloakServiceReconciled(cr *v1alpha1.Keycloak, currentState *v1.Service) *v1.Service {
	reconciled := currentState.DeepCopy()
	reconciled.Spec.Ports = GetServicePorts()
	return reconciled
}
