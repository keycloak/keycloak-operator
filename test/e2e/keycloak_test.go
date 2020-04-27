package e2e

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"

	apis "github.com/keycloak/keycloak-operator/pkg/apis"
	keycloakv1alpha1 "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/model"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

type testWithDeployedOperator func(*testing.T, *framework.Framework, *framework.TestCtx, string) error

func TestKeycloak(t *testing.T) {
	keycloakType := &keycloakv1alpha1.Keycloak{}
	err := framework.AddToFrameworkScheme(apis.AddToScheme, keycloakType)
	if err != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	}

	// run subtests
	t.Run("keycloakDeploymentTest", func(t *testing.T) {
		runTest(t, keycloakDeploymentTest)
	})
	t.Run("keycloakBackupTest", func(t *testing.T) {
		runTest(t, keycloakBackupTest)
	})
}

func keycloakDeploymentTest(t *testing.T, f *framework.Framework, ctx *framework.TestCtx, namespace string) error {
	//given
	keycloakCR := &keycloakv1alpha1.Keycloak{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testKeycloakCRName,
			Namespace: namespace,
		},
		Spec: keycloakv1alpha1.KeycloakSpec{
			Instances: 1,
			ExternalAccess: keycloakv1alpha1.KeycloakExternalAccess{
				Enabled: true,
			},
		},
	}

	//when - then
	err := Create(f, keycloakCR, ctx)
	if err != nil {
		return err
	}

	err = WaitForStatefulSetReplicasReady(t, f.KubeClient, model.ApplicationName, namespace)
	if err != nil {
		return err
	}

	keycloakCRName := types.NamespacedName{
		Namespace: keycloakCR.Namespace,
		Name:      keycloakCR.Name,
	}
	err = Get(f, keycloakCRName, keycloakCR, ctx)
	if err != nil {
		return err
	}

	// Kubernetes case - Ingress is on
	if _, ok := keycloakCR.Status.SecondaryResources["Ingress"]; ok {
		// Skipping TLS verification is actually part of the test. In Kubernetes, if there's no signing
		// manager installed, Keycloak will generate its own, self-signed cert. Of course
		// we don't have a matching truststore for it, hence we need to skip TLS verification.
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint
		err = WaitForCondition(t, f.KubeClient, func(t *testing.T, c kubernetes.Interface) error {
			response, err := http.Get("https://keycloak.local/auth")
			if err != nil {
				return err
			}
			response.Body.Close()
			if response.StatusCode == 200 {
				return nil
			}
			return fmt.Errorf("invalid response from Keycloak (%v)", response.Status)
		})
	}

	//TODO: OpenShift platform may additionally test the route.

	return err
}

func keycloakBackupTest(t *testing.T, f *framework.Framework, ctx *framework.TestCtx, namespace string) error {
	//given
	lab := map[string]string{"app": "sso"}
	labSel := metav1.LabelSelector{
		MatchLabels: lab,
	}

	keycloakCR := &keycloakv1alpha1.Keycloak{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testKeycloakCRName,
			Namespace: namespace,
			Labels:    lab,
		},
		Spec: keycloakv1alpha1.KeycloakSpec{
			Instances: 1,
		},
	}

	keycloakBackupCR := &keycloakv1alpha1.KeycloakBackup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testKeycloakCRName,
			Namespace: namespace,
		},
		Spec: keycloakv1alpha1.KeycloakBackupSpec{
			InstanceSelector: &labSel,
		},
	}

	//when - then
	err := Create(f, keycloakCR, ctx)
	if err != nil {
		return err
	}

	err = WaitForStatefulSetReplicasReady(t, f.KubeClient, model.ApplicationName, namespace)
	if err != nil {
		return err
	}

	err = f.Client.Create(context.TODO(), keycloakBackupCR, &framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		return err
	}

	err = WaitForPersistentVolumeClaimCreated(t, f.KubeClient, model.PostgresqlBackupPersistentVolumeName+"-"+testKeycloakCRName, namespace)
	if err != nil {
		return err
	}

	return err
}

func runTest(t *testing.T, testCase testWithDeployedOperator) {
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
	err = e2eutil.WaitForOperatorDeployment(t, f.KubeClient, namespace, testKeycloakCRName, 1, pollRetryInterval, pollTimeout)
	if err != nil {
		t.Fatal(err)
	}

	if err = testCase(t, f, ctx, namespace); err != nil {
		t.Fatal(err)
	}
}
