package model

import (
	kc "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	networkingv1 "k8s.io/api/networking/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func KeycloakIngress(cr *kc.Keycloak) *networkingv1.Ingress {
	ingressHost := cr.Spec.ExternalAccess.Host
	if ingressHost == "" {
		ingressHost = IngressDefaultHost
	}

	pathTypeImplementationSpecific := networkingv1.PathTypeImplementationSpecific // a workaround to get constant's address

	return &networkingv1.Ingress{
		ObjectMeta: v1.ObjectMeta{
			Name:      ApplicationName,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app": ApplicationName,
			},
			Annotations: map[string]string{
				"nginx.ingress.kubernetes.io/backend-protocol": "HTTPS",
				"nginx.ingress.kubernetes.io/server-snippet": `
                      location ~* "^/auth/realms/master/metrics" {
                          return 301 /auth/realms/master;
                        }`,
			},
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					Host: ingressHost,
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &pathTypeImplementationSpecific,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: ApplicationName,
											Port: networkingv1.ServiceBackendPort{
												Number: KeycloakServicePort,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func KeycloakIngressReconciled(cr *kc.Keycloak, currentState *networkingv1.Ingress) *networkingv1.Ingress {
	reconciled := currentState.DeepCopy()
	reconciledHost := currentState.Spec.Rules[0].Host
	reconciledSpecTLS := currentState.Spec.TLS
	pathTypeImplementationSpecific := networkingv1.PathTypeImplementationSpecific // a workaround to get constant's address

	reconciled.Spec = networkingv1.IngressSpec{
		TLS: reconciledSpecTLS,
		Rules: []networkingv1.IngressRule{
			{
				Host: reconciledHost,
				IngressRuleValue: networkingv1.IngressRuleValue{
					HTTP: &networkingv1.HTTPIngressRuleValue{
						Paths: []networkingv1.HTTPIngressPath{
							{
								Path:     "/",
								PathType: &pathTypeImplementationSpecific,
								Backend: networkingv1.IngressBackend{
									Service: &networkingv1.IngressServiceBackend{
										Name: ApplicationName,
										Port: networkingv1.ServiceBackendPort{
											Number: KeycloakServicePort,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return reconciled
}

func KeycloakIngressSelector(cr *kc.Keycloak) client.ObjectKey {
	return client.ObjectKey{
		Name:      ApplicationName,
		Namespace: cr.Namespace,
	}
}
