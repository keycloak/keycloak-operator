package e2e

import (
	"testing"

	keycloakv1alpha1 "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/test"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	userID = "test-user"
)

func NewKeycloakUserCRDTestStruct() *CRDTestStruct {
	return &CRDTestStruct{
		prepareEnvironmentSteps: []environmentInitializationStep{
			prepareKeycloaksCR,
			prepareKeycloakRealmCR,
		},
		testSteps: map[string]deployedOperatorTestStep{
			"keycloakUserBasicTest": {
				prepareTestEnvironmentSteps: []environmentInitializationStep{
					prepareKeycloakUserCR,
				},
				testFunction: keycloakUserBasicTest,
			},
		},
	}
}

func getKeycloakUserCR(namespace string) *keycloakv1alpha1.KeycloakUser {
	return &keycloakv1alpha1.KeycloakUser{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testKeycloakUserCRName,
			Namespace: namespace,
			Labels:    CreateLabel(namespace),
		},
		Spec: keycloakv1alpha1.KeycloakUserSpec{
			RealmSelector: &metav1.LabelSelector{
				MatchLabels: CreateLabel(namespace),
			},
			User: keycloakv1alpha1.KeycloakAPIUser{
				ID:            userID,
				UserName:      userID,
				FirstName:     "First name",
				LastName:      "Last name",
				Email:         "test-user@operator-test.email",
				EmailVerified: true,
				Enabled:       true,
			},
		},
	}
}

func prepareKeycloakUserCR(t *testing.T, framework *test.Framework, ctx *test.Context, namespace string) error {
	keycloakUserCR := getKeycloakUserCR(namespace)
	return Create(framework, keycloakUserCR, ctx)
}

func keycloakUserBasicTest(t *testing.T, framework *test.Framework, ctx *test.Context, namespace string) error {
	return WaitForUserToBeReady(t, framework, namespace)
}
