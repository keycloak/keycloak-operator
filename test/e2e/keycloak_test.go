package e2e

import (
	goctx "context"
	"fmt"
	"testing"
	"time"

	apis "github.com/keycloak/keycloak-operator/pkg/apis"
	keycloakv1alpha1 "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/model"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	testKeycloakCRDName  = "keycloak-test"
	retryInterval        = time.Second * 5
	timeout              = time.Second * 60
	cleanupRetryInterval = time.Second * 1
	cleanupTimeout       = time.Second * 5
)

type testWithDeployedOperator func(*testing.T, *framework.Framework, *framework.TestCtx) error

func TestKeycloak(t *testing.T) {
	keycloakType := &keycloakv1alpha1.Keycloak{}
	err := framework.AddToFrameworkScheme(apis.AddToScheme, keycloakType)
	if err != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	}
	// run subtests
	t.Run("keycloak-group", func(t *testing.T) {
		runTest(t, keycloakDeploymentTest)
	})
}

func keycloakDeploymentTest(t *testing.T, f *framework.Framework, ctx *framework.TestCtx) error {
	namespace, err := ctx.GetNamespace()
	if err != nil {
		return fmt.Errorf("could not get namespace: %v", err)
	}

	keycloakCRD := &keycloakv1alpha1.Keycloak{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testKeycloakCRDName,
			Namespace: namespace,
		},
		Spec: keycloakv1alpha1.KeycloakSpec{
			Instances: 1,
		},
	}

	err = f.Client.Create(goctx.TODO(), keycloakCRD, &framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		return err
	}

	t.Logf("waiting until StatefulSet %v becomes ready", model.ApplicationName)
	err = WaitForStatefulSetReplicasReady(t, f.KubeClient, model.ApplicationName, namespace)
	if err != nil {
		return err
	}

	//TODO: OpenShift platform may additionally test the route.

	return err
}

func runTest(t *testing.T, testCase testWithDeployedOperator) {
	t.Parallel()
	ctx := framework.NewTestCtx(t)
	defer ctx.Cleanup()
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
	err = e2eutil.WaitForOperatorDeployment(t, f.KubeClient, namespace, testKeycloakCRDName, 1, retryInterval, timeout)
	if err != nil {
		t.Fatal(err)
	}

	if err = testCase(t, f, ctx); err != nil {
		t.Fatal(err)
	}
}
