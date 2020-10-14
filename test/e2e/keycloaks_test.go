package e2e

import (
	"crypto/tls"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	keycloakv1alpha1 "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/model"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewKeycloaksCRDTestStruct() *CRDTestStruct {
	return &CRDTestStruct{
		prepareEnvironmentSteps: []environmentInitializationStep{
			prepareKeycloaksCRWithExtension,
		},
		testSteps: map[string]deployedOperatorTestStep{
			"keycloakDeploymentTest": {testFunction: keycloakDeploymentTest},
		},
	}
}

func getKeycloakCR(namespace string) *keycloakv1alpha1.Keycloak {
	return &keycloakv1alpha1.Keycloak{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testKeycloakCRName,
			Namespace: namespace,
			Labels:    CreateLabel(namespace),
		},
		Spec: keycloakv1alpha1.KeycloakSpec{
			Instances:      1,
			ExternalAccess: keycloakv1alpha1.KeycloakExternalAccess{Enabled: true},
			Profile:        currentProfile(),
		},
	}
}

func getDeployedKeycloakCR(framework *framework.Framework, namespace string) keycloakv1alpha1.Keycloak {
	keycloakCR := keycloakv1alpha1.Keycloak{}
	_ = GetNamespacedObject(framework, namespace, testKeycloakCRName, &keycloakCR)
	return keycloakCR
}

func prepareKeycloaksCR(t *testing.T, f *framework.Framework, ctx *framework.Context, namespace string) error {
	return deployKeycloaksCR(t, f, ctx, namespace, getKeycloakCR(namespace))
}

func prepareKeycloaksCRWithExtension(t *testing.T, f *framework.Framework, ctx *framework.Context, namespace string) error {
	keycloakCR := getKeycloakCR(namespace)
	keycloakCR.Spec.Extensions = []string{"https://github.com/aerogear/keycloak-metrics-spi/releases/download/1.0.4/keycloak-metrics-spi-1.0.4.jar"}

	return deployKeycloaksCR(t, f, ctx, namespace, keycloakCR)
}

func deployKeycloaksCR(t *testing.T, f *framework.Framework, ctx *framework.Context, namespace string, keycloakCR *keycloakv1alpha1.Keycloak) error {
	err := doWorkaroundIfNecessary(f, ctx, namespace)
	if err != nil {
		return err
	}

	err = Create(f, keycloakCR, ctx)
	if err != nil {
		return err
	}

	err = WaitForStatefulSetReplicasReady(t, f.KubeClient, model.ApplicationName, namespace)
	if err != nil {
		return err
	}

	return err
}

func keycloakDeploymentTest(t *testing.T, f *framework.Framework, ctx *framework.Context, namespace string) error {
	keycloakCR := getDeployedKeycloakCR(f, namespace)
	keycloakInternalURL := keycloakCR.Status.InternalURL
	assert.NotEmpty(t, keycloakInternalURL)

	// Skipping TLS verification is actually part of the test. In Kubernetes, if there's no signing
	// manager installed, Keycloak will generate its own, self-signed cert. Of course
	// we don't have a matching truststore for it, hence we need to skip TLS verification.
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint

	err := WaitForSuccessResponse(t, f, keycloakInternalURL+"/auth")
	if err != nil {
		return err
	}

	err = WaitForSuccessResponse(t, f, keycloakInternalURL+"/auth/realms/master/metrics")
	return err
}
