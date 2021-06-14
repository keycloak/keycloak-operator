package model

import (
	"testing"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/stretchr/testify/assert"
	networkingv1 "k8s.io/api/networking/v1"
)

func TestKeycloakIngress_testTLSOverride(t *testing.T) {
	//given
	currentState := &networkingv1.Ingress{
		Spec: networkingv1.IngressSpec{
			TLS: []networkingv1.IngressTLS{
				{
					Hosts: []string{
						IngressDefaultHost,
					},
					SecretName: "keycloak-secret",
				},
			},
			Rules: []networkingv1.IngressRule{
				{
					Host: IngressDefaultHost,
				},
			},
		},
	}
	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			ExternalAccess: v1alpha1.KeycloakExternalAccess{
				Enabled: true,
			},
		},
	}

	//when
	reconciledIngress := KeycloakIngressReconciled(cr, currentState)

	//then
	assert.Equal(t, 1, len(reconciledIngress.Spec.TLS))
	assert.Equal(t, 1, len(reconciledIngress.Spec.TLS[0].Hosts))
	assert.Equal(t, IngressDefaultHost, reconciledIngress.Spec.TLS[0].Hosts[0])
	assert.Equal(t, "keycloak-secret", reconciledIngress.Spec.TLS[0].SecretName)
}

func TestKeycloakIngress_testHost(t *testing.T) {
	//given
	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			ExternalAccess: v1alpha1.KeycloakExternalAccess{
				Enabled: true,
			},
		},
	}

	//when
	ingress := KeycloakIngress(cr)

	//then
	assert.Equal(t, IngressDefaultHost, ingress.Spec.Rules[0].Host)
}

func TestKeycloakIngress_testHostReconciled(t *testing.T) {
	//given
	currentState := &networkingv1.Ingress{
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					Host: IngressDefaultHost,
				},
			},
		},
	}
	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			ExternalAccess: v1alpha1.KeycloakExternalAccess{
				Enabled: true,
			},
		},
	}

	//when
	reconciledIngress := KeycloakIngressReconciled(cr, currentState)

	//then
	assert.Equal(t, IngressDefaultHost, reconciledIngress.Spec.Rules[0].Host)
}

func TestKeycloakIngress_testHostOverride(t *testing.T) {
	//given
	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			ExternalAccess: v1alpha1.KeycloakExternalAccess{
				Enabled: true,
				Host:    "host-override",
			},
		},
	}

	//when
	ingress := KeycloakIngress(cr)

	//then
	assert.Equal(t, "host-override", ingress.Spec.Rules[0].Host)
}

func TestKeycloakIngress_testHostOverrideReconciled(t *testing.T) {
	//given
	currentState := &networkingv1.Ingress{
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					Host: "host-override",
				},
			},
		},
	}
	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			ExternalAccess: v1alpha1.KeycloakExternalAccess{
				Enabled: true,
				Host:    "host-override",
			},
		},
	}

	//when
	reconciledIngress := KeycloakIngressReconciled(cr, currentState)

	//then
	assert.Equal(t, "host-override", reconciledIngress.Spec.Rules[0].Host)
}
