package model

import (
	"fmt"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func getSpec(dbSecret *v1.Secret, serviceTypeExternal bool) v1.ServiceSpec {
	spec := v1.ServiceSpec{}

	if serviceTypeExternal {
		spec.Type = v1.ServiceTypeExternalName
		spec.Selector = nil
		spec.ExternalName = GetExternalDatabaseHost(dbSecret)
	} else {
		spec.Type = v1.ServiceTypeClusterIP
		spec.Selector = map[string]string{
			"app":       ApplicationName,
			"component": PostgresqlDeploymentComponent,
		}
	}

	spec.Ports = []v1.ServicePort{
		{
			Port:       GetExternalDatabasePort(dbSecret),
			TargetPort: intstr.Parse(fmt.Sprintf("%d", GetExternalDatabasePort(dbSecret))),
		},
	}

	return spec
}

func PostgresqlService(cr *v1alpha1.Keycloak, dbSecret *v1.Secret, serviceTypeExternal bool) *v1.Service {
	return &v1.Service{
		ObjectMeta: v12.ObjectMeta{
			Name:      PostgresqlServiceName,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app": ApplicationName,
			},
		},
		Spec: getSpec(dbSecret, serviceTypeExternal),
	}
}

func PostgresqlServiceSelector(cr *v1alpha1.Keycloak) client.ObjectKey {
	return client.ObjectKey{
		Name:      PostgresqlServiceName,
		Namespace: cr.Namespace,
	}
}

func PostgresqlServiceReconciled(currentState *v1.Service, dbSecret *v1.Secret, serviceTypeExternal bool) *v1.Service {
	reconciled := currentState.DeepCopy()
	if !serviceTypeExternal {
		reconciled.Spec.Type = v1.ServiceTypeClusterIP
		reconciled.Spec.Selector = map[string]string{
			"app":       ApplicationName,
			"component": PostgresqlDeploymentComponent,
		}
		reconciled.Spec.Ports = []v1.ServicePort{
			{
				Port:       5432,
				TargetPort: intstr.Parse("5432"),
			},
		}
	} else {
		reconciled.Spec = getSpec(dbSecret, serviceTypeExternal)
	}
	return reconciled
}
