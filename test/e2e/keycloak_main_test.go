package e2e

import (
	"context"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/keycloak/keycloak-operator/pkg/apis"
	keycloakv1alpha1 "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
)

type deployedOperatorTestStep struct {
	prepareTestEnvironmentSteps []environmentInitializationStep
	testFunction                func(*testing.T, *framework.Framework, *framework.Context, string) error
}

type environmentInitializationStep func(*testing.T, *framework.Framework, *framework.Context, string) error

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
		runTestsFromCRDInterface(t, NewUnmanagedKeycloaksCRDTestStruct())
	})
	t.Run("KeycloakBackupCRDTest", func(t *testing.T) {
		runTestsFromCRDInterface(t, NewKeycloakBackupCRDTestStruct())
	})
	t.Run("KeycloakRealmsCRDTest", func(t *testing.T) {
		runTestsFromCRDInterface(t, NewKeycloakRealmsCRDTestStruct())
	})
	t.Run("KeycloakUsersCRDTest", func(t *testing.T) {
		runTestsFromCRDInterface(t, NewKeycloakUserCRDTestStruct())
	})
	t.Run("KeycloakClientsCRDTest", func(t *testing.T) {
		runTestsFromCRDInterface(t, NewKeycloakClientsCRDTestStruct())
	})
}

func runTestsFromCRDInterface(t *testing.T, crd *CRDTestStruct) {
	globalCTX := framework.NewContext(t)
	defer globalCTX.Cleanup()

	err := globalCTX.InitializeClusterResources(&framework.CleanupOptions{TestContext: globalCTX, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		t.Fatalf("failed to initialize cluster resources: %v", err)
	}
	t.Log("initialized cluster resources")
	namespace, err := globalCTX.GetOperatorNamespace()
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
		deployment, err := f.KubeClient.AppsV1().Deployments(namespace).Get(context.TODO(), operatorCRName, metav1.GetOptions{})
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Operator deployed from: %s", deployment.Spec.Template.Spec.Containers[0].Image)
	}

	t.Log("prepare global CRD environment")
	for _, prepareEnvironmentFunction := range crd.prepareEnvironmentSteps {
		err = prepareEnvironmentFunction(t, f, globalCTX, namespace)

		if err != nil {
			t.Fatal(err)
		}
	}

	for testName, testStep := range crd.testSteps {
		testName := testName
		testStep := testStep
		t.Run(testName, func(t *testing.T) {
			t.Logf("test %s started", testName)
			testCTX := framework.NewContext(t)

			t.Logf("prepare test environment for test %s", testName)
			for _, prepareEnvironmentFunction := range testStep.prepareTestEnvironmentSteps {
				err = prepareEnvironmentFunction(t, f, testCTX, namespace)

				if err != nil {
					t.Logf("preparation step for test %s failed, cleaning context", testName)
					testCTX.Cleanup()
					t.Fatal(err)
				}
			}

			t.Logf("running test %s", testName)
			if err = testStep.testFunction(t, f, testCTX, namespace); err != nil {
				t.Logf("test %s failed, cleaning context", testName)
				testCTX.Cleanup()
				t.Fatal(err)
			}

			t.Logf("cleanup for test %s", testName)
			testCTX.Cleanup()

			t.Logf("test finished %s", testName)
		})
	}
}
