package model

import (
	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func PostgresqlService(cr *v1alpha1.Keycloak) *v1.Service {
	return &v1.Service{
		ObjectMeta: v12.ObjectMeta{
			Name:      PostgresqlServiceName,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app": ApplicationName,
			},
		},
		Spec: v1.ServiceSpec{
			Selector: map[string]string{
				"app":       ApplicationName,
				"component": PostgresqlDeploymentComponent,
			},
			Ports: []v1.ServicePort{
				{
					Port:       5432,
					TargetPort: intstr.Parse("5432"),
				},
			},
		},
	}
}

func PostgresqlServiceSelector(cr *v1alpha1.Keycloak) client.ObjectKey {
	return client.ObjectKey{
		Name:      PostgresqlServiceName,
		Namespace: cr.Namespace,
	}
}

func PostgresqlServiceReconciled(cr *v1alpha1.Keycloak, currentState *v1.Service) *v1.Service {
	reconciled := currentState.DeepCopy()
	reconciled.Spec.Ports = []v1.ServicePort{
		{
			Port:       5432,
			TargetPort: intstr.Parse("5432"),
		},
	}
	return reconciled
}
