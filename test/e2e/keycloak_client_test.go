package e2e

import (
	"testing"

	keycloakv1alpha1 "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/common"
	"github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/pkg/errors"
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
			"keycloakClientRolesTest": {
				testFunction: keycloakClientRolesTest,
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
				FullScopeAllowed:          &[]bool{true}[0],
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

func keycloakClientRolesTest(t *testing.T, framework *test.Framework, ctx *test.Context, namespace string) error {
	// create
	client := getKeycloakClientCR(namespace, false)
	client.Spec.Roles = []keycloakv1alpha1.RoleRepresentation{{Name: "a"}, {Name: "b"}, {Name: "c"}}
	err := Create(framework, client, ctx)
	if err != nil {
		return err
	}
	err = WaitForClientToBeReady(t, framework, namespace, testKeycloakClientCRName)
	if err != nil {
		return err
	}

	// update client: delete/rename/leave/add role
	keycloakCR := getDeployedKeycloakCR(framework, namespace)
	authenticatedClient, err := MakeAuthenticatedClient(keycloakCR)
	if err != nil {
		return err
	}
	retrievedRoles, err := authenticatedClient.ListClientRoles(clientName, realmName)
	if err != nil {
		return err
	}
	var bID string
	for _, role := range retrievedRoles {
		if role.Name == "b" {
			bID = role.ID
			break
		}
	}
	err = GetNamespacedObject(framework, namespace, testKeycloakClientCRName, client)
	if err != nil {
		return err
	}
	client.Spec.Roles = []keycloakv1alpha1.RoleRepresentation{{ID: bID, Name: "b2"}, {Name: "c"}, {Name: "d"}}
	err = Update(framework, client)
	if err != nil {
		return err
	}
	// check role presence directly as a "ready" CR might just be stale
	err = waitForClientRoles(t, framework, keycloakCR, client, client.Spec.Roles)
	if err != nil {
		return err
	}
	return WaitForClientToBeReady(t, framework, namespace, testKeycloakClientCRName)
}

func waitForClientRoles(t *testing.T, framework *test.Framework, keycloakCR keycloakv1alpha1.Keycloak, clientCR *keycloakv1alpha1.KeycloakClient, expected []keycloakv1alpha1.RoleRepresentation) error {
	return WaitForConditionWithClient(t, framework, keycloakCR, func(authenticatedClient common.KeycloakInterface) error {
		roles, err := authenticatedClient.ListClientRoles(clientCR.Spec.Client.ID, realmName)
		if err != nil {
			return err
		}

		error := false
		if len(roles) != len(expected) {
			error = true
		} else {
			for _, expectedRole := range expected {
				found := false
				for _, role := range roles {
					if role.Name == expectedRole.Name && (expectedRole.ID == "" || role.ID == expectedRole.ID) {
						found = true
						break
					}
				}
				if !found {
					error = true
					break
				}
			}
		}

		if error {
			return errors.Errorf("role names not as expected:\n"+
				"expected: %v\n"+
				"actual  : %v",
				expected, roles)
		}
		return nil
	})
}
