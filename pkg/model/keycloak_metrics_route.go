package model

import (
	kc "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	v1 "github.com/openshift/api/route/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func KeycloakMetricsRoute(cr *kc.Keycloak, keycloakMainRoute *v1.Route) *v1.Route {
	keycloakMainRouteCopy := keycloakMainRoute.DeepCopy()
	return &v1.Route{
		ObjectMeta: v12.ObjectMeta{
			Name:      KeycloakMetricsRouteName,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app": ApplicationName,
			},
			Annotations: map[string]string{
				"haproxy.router.openshift.io/balance":        RouteLoadBalancingStrategy,
				"haproxy.router.openshift.io/rewrite-target": KeycloakMetricsRouteRewritePath,
			},
		},
		Spec: v1.RouteSpec{
			Host: keycloakMainRouteCopy.Spec.Host,
			Path: KeycloakMetricsRoutePath,
			Port: keycloakMainRouteCopy.Spec.Port,
			TLS:  keycloakMainRouteCopy.Spec.TLS,
			To:   keycloakMainRouteCopy.Spec.To,
		},
	}
}

func KeycloakMetricsRouteReconciled(cr *kc.Keycloak, currentState *v1.Route, keycloakMainRoute *v1.Route) *v1.Route {
	reconciled := currentState.DeepCopy()
	reconciled.Annotations = map[string]string{
		"haproxy.router.openshift.io/balance":        RouteLoadBalancingStrategy,
		"haproxy.router.openshift.io/rewrite-target": KeycloakMetricsRouteRewritePath,
	}
	reconciled.Spec = v1.RouteSpec{
		Host: keycloakMainRoute.Spec.Host,
		Path: KeycloakMetricsRoutePath,
		Port: keycloakMainRoute.Spec.Port,
		TLS:  keycloakMainRoute.Spec.TLS.DeepCopy(),
		To:   keycloakMainRoute.Spec.To,
	}

	return reconciled
}

func KeycloakMetricsRouteSelector(cr *kc.Keycloak) client.ObjectKey {
	return client.ObjectKey{
		Name:      KeycloakMetricsRouteName,
		Namespace: cr.Namespace,
	}
}
