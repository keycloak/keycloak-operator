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
	roleName           = "test-client-role"
	scopeName          = "test-client-scope"
	resourceName       = "test-client-resource"
	policyName         = "test-client-policy"
	permissionName     = "test-client-permission"
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
				ID:                           id,
				ClientID:                     id,
				Secret:                       id,
				Name:                         id,
				SurrogateAuthRequired:        false,
				Enabled:                      true,
				BaseURL:                      "https://operator-test.url/client-base-url",
				AdminURL:                     "https://operator-test.url/client-admin-url",
				RootURL:                      "https://operator-test.url/client-root-url",
				Description:                  "Client used within operator tests",
				WebOrigins:                   []string{"https://operator-test.url"},
				BearerOnly:                   false,
				ConsentRequired:              false,
				StandardFlowEnabled:          true,
				ImplicitFlowEnabled:          false,
				DirectAccessGrantsEnabled:    true,
				ServiceAccountsEnabled:       true,
				PublicClient:                 false,
				FrontchannelLogout:           false,
				Protocol:                     "openid-connect",
				FullScopeAllowed:             true,
				NodeReRegistrationTimeout:    -1,
				DefaultClientScopes:          []string{"profile"},
				OptionalClientScopes:         []string{"microprofile-jwt"},
				DefaultRoles:                 []string{roleName},
				ClientAuthenticatorType:      "client-secret",
				AuthorizationServicesEnabled: true,
				AuthorizationSettings: &keycloakv1alpha1.KeycloakResourceServer{
					AllowRemoteResourceManagement: false,
					DecisionStrategy:              "AFFIRMATIVE",
					Policies: []keycloakv1alpha1.KeycloakPolicy{{
						ID:               policyName,
						Name:             policyName,
						Description:      policyName,
						Type:             "role",
						Logic:            "POSITIVE",
						DecisionStrategy: "UNANIMOUS",
						Config: map[string]string{
							"roles": "[{\"id\":\"" + id + "/" + roleName + "\",\"required\":true}]",
						},
					},
						{
							ID:               permissionName,
							Name:             permissionName,
							Description:      permissionName,
							Type:             "scope",
							Logic:            "POSITIVE",
							DecisionStrategy: "UNANIMOUS",
							Config: map[string]string{
								"resources":     "[\"" + resourceName + "\"]",
								"scopes":        "[\"" + scopeName + "\"]",
								"applyPolicies": "[\"" + policyName + "\"]",
							},
						}},
					PolicyEnforcementMode: "ENFORCING",
					Resources: []keycloakv1alpha1.KeycloakResource{{
						ID:                 resourceName,
						DisplayName:        resourceName,
						Name:               resourceName,
						OwnerManagedAccess: false,
						Scopes: []keycloakv1alpha1.KeycloakScope{{
							Name: scopeName,
						}},
						Uris: []string{resourceName + "/*"},
					}},
					Scopes: []keycloakv1alpha1.KeycloakScope{{
						ID:          scopeName,
						Name:        scopeName,
						DisplayName: scopeName,
					}},
				},
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
