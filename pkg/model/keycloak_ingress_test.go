package model

import (
	"testing"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/api/extensions/v1beta1"
)

func TestKeycloakIngress_testTLSOverride(t *testing.T) {
	//given
	currentState := &v1beta1.Ingress{
		Spec: v1beta1.IngressSpec{
			TLS: []v1beta1.IngressTLS{
				{
					Hosts: []string{
						IngressDefaultHost,
					},
					SecretName: "keycloak-secret",
				},
			},
			Rules: []v1beta1.IngressRule{
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
