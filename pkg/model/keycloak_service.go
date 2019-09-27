package model

import (
	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func KeycloakService(cr *v1alpha1.Keycloak) *v1.Service {
	return &v1.Service{
		ObjectMeta: v12.ObjectMeta{
			Name:      ApplicationName,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"application": ApplicationName,
			},
		},
		Spec: v1.ServiceSpec{
			Selector: map[string]string{
				"application": ApplicationName,
				"component":   KeycloakDeploymentComponent,
			},
			Ports: []v1.ServicePort{
				{
					Port:       8443,
					TargetPort: intstr.FromInt(8443),
				},
			},
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
	reconciled.Spec.Ports = []v1.ServicePort{
		{
			Port:       8443,
			TargetPort: intstr.FromInt(8443),
		},
	}
	return reconciled
}
