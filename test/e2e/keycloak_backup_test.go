package e2e

import (
	"context"
	"testing"

	keycloakv1alpha1 "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/model"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewKeycloakBackupCRDTestStruct() *CRDTestStruct {
	return &CRDTestStruct{
		prepareEnvironmentSteps: []environmentInitializationStep{
			prepareKeycloaksCR,
		},
		testSteps: map[string]deployedOperatorTestStep{
			"keycloakDeploymentTest": {
				prepareTestEnvironmentSteps: []environmentInitializationStep{
					prepareKeycloakBackupCR,
				},
				testFunction: keycloakBackupTest,
			},
		},
	}
}

func getKeycloakBackupCR(namespace string) *keycloakv1alpha1.KeycloakBackup {
	labSel := metav1.LabelSelector{
		MatchLabels: CreateLabel(namespace),
	}

	return &keycloakv1alpha1.KeycloakBackup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testKeycloakCRName,
			Namespace: namespace,
		},
		Spec: keycloakv1alpha1.KeycloakBackupSpec{
			InstanceSelector: &labSel,
		},
	}
}

func prepareKeycloakBackupCR(t *testing.T, f *framework.Framework, ctx *framework.Context, namespace string) error {
	keycloakBackupCR := getKeycloakBackupCR(namespace)

	err := f.Client.Create(context.TODO(), keycloakBackupCR, &framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		return err
	}

	return err
}

func keycloakBackupTest(t *testing.T, f *framework.Framework, ctx *framework.Context, namespace string) error {
	err := WaitForPersistentVolumeClaimCreated(t, f.KubeClient, model.PostgresqlBackupPersistentVolumeName+"-"+testKeycloakCRName, namespace)
	if err != nil {
		return err
	}

	return err
}
