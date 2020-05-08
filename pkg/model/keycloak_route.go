package model

import (
	kc "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	v1 "github.com/openshift/api/route/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func KeycloakRoute(cr *kc.Keycloak) *v1.Route {
	return &v1.Route{
		ObjectMeta: v12.ObjectMeta{
			Name:      ApplicationName,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app": ApplicationName,
			},
			Annotations: map[string]string{
				"haproxy.router.openshift.io/balance": RouteLoadBalancingStrategy,
			},
		},
		Spec: v1.RouteSpec{
			Port: &v1.RoutePort{
				TargetPort: intstr.FromString(ApplicationName),
			},
			TLS: &v1.TLSConfig{
				Termination: getTLSTerminationType(cr),
			},
			To: v1.RouteTargetReference{
				Kind: "Service",
				Name: ApplicationName,
			},
		},
	}
}

func KeycloakRouteReconciled(cr *kc.Keycloak, currentState *v1.Route) *v1.Route {
	reconciled := currentState.DeepCopy()
	reconciled.Spec = v1.RouteSpec{
		Port: &v1.RoutePort{
			TargetPort: intstr.FromString(ApplicationName),
		},
		TLS: &v1.TLSConfig{
			Termination: getTLSTerminationType(cr),
		},
		To: v1.RouteTargetReference{
			Kind: "Service",
			Name: ApplicationName,
		},
	}

	return reconciled
}

func getTLSTerminationType(cr *kc.Keycloak) v1.TLSTerminationType {
	if cr.Spec.ExternalAccess.TLSTermination == kc.PassthroughTLSTerminationType {
		return "passthrough"
	}
	return "reencrypt"
}

func KeycloakRouteSelector(cr *kc.Keycloak) client.ObjectKey {
	return client.ObjectKey{
		Name:      ApplicationName,
		Namespace: cr.Namespace,
	}
}
