package e2e

import (
	"testing"

	keycloakv1alpha1 "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/test"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	groupID = "test-group"
)

func NewKeycloakGroupCRDTestStruct() *CRDTestStruct {
	return &CRDTestStruct{
		prepareEnvironmentSteps: []environmentInitializationStep{
			prepareKeycloaksCR,
			prepareKeycloakRealmCR,
		},
		testSteps: map[string]deployedOperatorTestStep{
			"keycloakGroupBasicTest": {
				prepareTestEnvironmentSteps: []environmentInitializationStep{
					prepareKeycloakGroupCR,
				},
				testFunction: keycloakGroupBasicTest,
			},
		},
	}
}

func getKeycloakGroupCR(namespace string) *keycloakv1alpha1.KeycloakGroup {
	return &keycloakv1alpha1.KeycloakGroup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testKeycloakGroupCRName,
			Namespace: namespace,
			Labels:    CreateLabel(namespace),
		},
		Spec: keycloakv1alpha1.KeycloakGroupSpec{
			RealmSelector: &metav1.LabelSelector{
				MatchLabels: CreateLabel(namespace),
			},
			Group: keycloakv1alpha1.KeycloakAPIGroup{
				ID:   groupID,
				Name: groupID,
			},
		},
	}
}

func prepareKeycloakGroupCR(t *testing.T, framework *test.Framework, ctx *test.Context, namespace string) error {
	keycloakUserCR := getKeycloakGroupCR(namespace)
	return Create(framework, keycloakUserCR, ctx)
}

func keycloakGroupBasicTest(t *testing.T, framework *test.Framework, ctx *test.Context, namespace string) error {
	return WaitForGroupToBeReady(t, framework, namespace)
}
