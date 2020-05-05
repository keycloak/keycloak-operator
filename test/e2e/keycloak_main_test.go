package e2e

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/keycloak/keycloak-operator/pkg/apis"
	keycloakv1alpha1 "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
)

type deployedOperatorTestStep func(*testing.T, *framework.Framework, *framework.TestCtx, string) error
type environmentInitializationStep func(*testing.T, *framework.Framework, *framework.TestCtx, string) error

type CRDTestStruct struct {
	prepareEnvironmentSteps []environmentInitializationStep
	testSteps               map[string]deployedOperatorTestStep
}

func TestKeycloakCRDS(t *testing.T) {
	keycloakType := &keycloakv1alpha1.Keycloak{}
	err := framework.AddToFrameworkScheme(apis.AddToScheme, keycloakType)
	if err != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	}

	t.Run("KeycloaksCRDTest", func(t *testing.T) {
		runTestsFromCRDInterface(t, NewKeycloaksCRDTestStruct())
	})
	t.Run("KeycloakBackupCRDTest", func(t *testing.T) {
		runTestsFromCRDInterface(t, NewKeycloakBackupCRDTestStruct())
	})
}

func runTestsFromCRDInterface(t *testing.T, crd *CRDTestStruct) {

	for testName, testMethod := range crd.testSteps {
		ctx := framework.NewTestCtx(t)

		err := ctx.InitializeClusterResources(&framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
		if err != nil {
			t.Fatalf("failed to initialize cluster resources: %v", err)
		}
		t.Log("initialized cluster resources")
		namespace, err := ctx.GetNamespace()
		if err != nil {
			t.Fatal(err)
		}
		// get global framework variables
		f := framework.Global
		// wait for Keycloak Operator to be ready
		err = e2eutil.WaitForOperatorDeployment(t, f.KubeClient, namespace, operatorCRName, 1, pollRetryInterval, pollTimeout)
		if err != nil {
			t.Fatal(err)
		}

		if !f.LocalOperator {
			deployment, err := f.KubeClient.AppsV1().Deployments(namespace).Get(operatorCRName, metav1.GetOptions{})
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("Operator deployed from: %s", deployment.Spec.Template.Spec.Containers[0].Image)
		}

		for _, prepareEnvironmentFunction := range crd.prepareEnvironmentSteps {
			err = prepareEnvironmentFunction(t, f, ctx, namespace)

			if err != nil {
				t.Fatal(err)
			}
		}

		t.Run(testName, func(t *testing.T) {
			if err = testMethod(t, f, ctx, namespace); err != nil {
				t.Fatal(err)
			}
		})

		ctx.Cleanup()
	}
}
