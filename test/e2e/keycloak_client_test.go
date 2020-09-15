package e2e

import (
	"testing"

	keycloakv1alpha1 "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/test"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	clientName         = "test-client"
	externalClientName = "test-client-external"
)

func NewKeycloakClientsCRDTestStruct() *CRDTestStruct {
	return &CRDTestStruct{
		prepareEnvironmentSteps: []environmentInitializationStep{
			prepareKeycloaksCR,
			prepareExternalKeycloaksCR,
			prepareKeycloakRealmCR,
		},
		testSteps: map[string]deployedOperatorTestStep{
			"keycloakClientBasicTest": {
				prepareTestEnvironmentSteps: []environmentInitializationStep{
					prepareKeycloakClientCR,
				},
				testFunction: keycloakClientBasicTest,
			},
			"externalKeycloakClientBasicTest": {
				prepareTestEnvironmentSteps: []environmentInitializationStep{
					prepareExternalKeycloakClientCR,
				},
				testFunction: externalKeycloakClientBasicTest,
			},
		},
	}
}

func getKeycloakClientCR(namespace string, external bool) *keycloakv1alpha1.KeycloakClient {
	k8sName := testKeycloakClientCRName
	id := clientName
	labels := CreateLabel(namespace)

	if external {
		k8sName = testExternalKeycloakClientCRName
		id = externalClientName
		labels = CreateExternalLabel(namespace)
	}

	return &keycloakv1alpha1.KeycloakClient{
		ObjectMeta: metav1.ObjectMeta{
			Name:      k8sName,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: keycloakv1alpha1.KeycloakClientSpec{
			RealmSelector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Client: &keycloakv1alpha1.KeycloakAPIClient{
				ID:                        id,
				ClientID:                  id,
				Name:                      id,
				SurrogateAuthRequired:     false,
				Enabled:                   true,
				BaseURL:                   "https://operator-test.url/client-base-url",
				AdminURL:                  "https://operator-test.url/client-admin-url",
				RootURL:                   "https://operator-test.url/client-root-url",
				Description:               "Client used within operator tests",
				WebOrigins:                []string{"https://operator-test.url"},
				BearerOnly:                false,
				ConsentRequired:           false,
				StandardFlowEnabled:       true,
				ImplicitFlowEnabled:       false,
				DirectAccessGrantsEnabled: true,
				ServiceAccountsEnabled:    false,
				PublicClient:              true,
				FrontchannelLogout:        false,
				Protocol:                  "openid-connect",
				FullScopeAllowed:          true,
				NodeReRegistrationTimeout: -1,
				DefaultClientScopes:       []string{"profile"},
				OptionalClientScopes:      []string{"microprofile-jwt"},
			},
		},
	}
}

func prepareKeycloakClientCR(t *testing.T, framework *test.Framework, ctx *test.Context, namespace string) error {
	keycloakClientCR := getKeycloakClientCR(namespace, false)
	return Create(framework, keycloakClientCR, ctx)
}

func prepareExternalKeycloakClientCR(t *testing.T, framework *test.Framework, ctx *test.Context, namespace string) error {
	keycloakClientCR := getKeycloakClientCR(namespace, true)
	return Create(framework, keycloakClientCR, ctx)
}

func keycloakClientBasicTest(t *testing.T, framework *test.Framework, ctx *test.Context, namespace string) error {
	return WaitForClientToBeReady(t, framework, namespace, testKeycloakClientCRName)
}

func externalKeycloakClientBasicTest(t *testing.T, framework *test.Framework, ctx *test.Context, namespace string) error {
	return WaitForClientToBeReady(t, framework, namespace, testExternalKeycloakClientCRName)
}
