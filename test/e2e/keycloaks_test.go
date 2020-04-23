package e2e

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"

	"k8s.io/client-go/kubernetes"

	keycloakv1alpha1 "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/model"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func NewKeycloaksCRDTestStruct() *CRDTestStruct {
	return &CRDTestStruct{
		prepareEnvironmentSteps: []environmentInitializationStep{
			prepareKeycloaksCR,
		},
		testSteps: map[string]deployedOperatorTestStep{
			"keycloakDeploymentTest": keycloakDeploymentTest,
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

func prepareKeycloaksCR(t *testing.T, f *framework.Framework, ctx *framework.TestCtx, namespace string) error {
	keycloakCR := getKeycloakCR(namespace)
	err := Create(f, keycloakCR, ctx)
	if err != nil {
		return err
	}

	err = WaitForStatefulSetReplicasReady(t, f.KubeClient, model.ApplicationName, namespace)
	if err != nil {
		return err
	}

	return err
}

func keycloakDeploymentTest(t *testing.T, f *framework.Framework, ctx *framework.TestCtx, namespace string) error {
	keycloakCRName := types.NamespacedName{
		Namespace: namespace,
		Name:      testKeycloakCRName,
	}

	keycloakCR := &keycloakv1alpha1.Keycloak{}

	err := Get(f, keycloakCRName, keycloakCR)
	if err != nil {
		return err
	}

	// Skipping TLS verification is actually part of the test. In Kubernetes, if there's no signing
	// manager installed, Keycloak will generate its own, self-signed cert. Of course
	// we don't have a matching truststore for it, hence we need to skip TLS verification.
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint
	err = WaitForCondition(t, f.KubeClient, func(t *testing.T, c kubernetes.Interface) error {
		response, err := http.Get(keycloakCR.Status.InternalURL + "/auth")
		if err != nil {
			return err
		}
		response.Body.Close()
		if response.StatusCode == 200 {
			return nil
		}
		return fmt.Errorf("invalid response from Keycloak (%v)", response.Status)
	})
	return err
}
