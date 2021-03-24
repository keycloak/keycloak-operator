package e2e

import (
	"testing"

	keycloakv1alpha1 "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/common"
	"github.com/keycloak/keycloak-operator/pkg/model"
	"github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	clientName         = "test-client"
	secondClientName   = "test-client-second"
	externalClientName = "test-client-external"
	authZClientName    = "test-client-authz"
)

func NewKeycloakClientsCRDTestStruct() *CRDTestStruct {
	return &CRDTestStruct{
		prepareEnvironmentSteps: []environmentInitializationStep{
			prepareKeycloaksCR,
			prepareExternalKeycloaksCR,
			prepareKeycloakRealmWithRolesCR,
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
			"keycloakClientAuthZSettingsTest": {
				prepareTestEnvironmentSteps: []environmentInitializationStep{
					prepareKeycloakClientAuthZCR,
				},
				testFunction: keycloakClientAuthZTest,
			},
			"keycloakClientRolesTest": {
				testFunction: keycloakClientRolesTest,
			},
			"keycloakClientScopeMappingsTest": {
				prepareTestEnvironmentSteps: []environmentInitializationStep{
					prepareKeycloakClientWithRolesCR,
				},
				testFunction: keycloakClientScopeMappingsTest,
			},
		},
	}
}

func prepareKeycloakRealmWithRolesCR(t *testing.T, framework *test.Framework, ctx *test.Context, namespace string) error {
	keycloakRealmCR := getKeycloakRealmCR(namespace)
	keycloakRealmCR.Spec.Realm.Roles = &keycloakv1alpha1.RolesRepresentation{}
	for _, roleName := range []string{"realmRoleA", "realmRoleB", "realmRoleC"} {
		keycloakRealmCR.Spec.Realm.Roles.Realm = append(keycloakRealmCR.Spec.Realm.Roles.Realm, keycloakv1alpha1.RoleRepresentation{
			ID:   roleName,
			Name: roleName,
		})
	}
	return Create(framework, keycloakRealmCR, ctx)
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

func getKeycloakClientAuthZCR(namespace string) *keycloakv1alpha1.KeycloakClient {
	k8sName := testAuthZKeycloakClientCRName
	id := authZClientName
	labels := CreateLabel(namespace)

	audioResourceType := "urn:" + id + ":resources:audio"
	imageResourceType := "urn:" + id + ":resources:image"

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
				Name:                         id,
				Description:                  "AuthZ Client used within operator tests",
				PublicClient:                 false,
				ServiceAccountsEnabled:       true,
				AuthorizationServicesEnabled: true,
				AuthorizationSettings: &keycloakv1alpha1.KeycloakResourceServer{
					Resources: []keycloakv1alpha1.KeycloakResource{
						{
							Name: "Audio Resource",
							Uris: []string{"/audio"},
							Type: audioResourceType,
							Scopes: []apiextensionsv1.JSON{
								{Raw: []byte(`{"name": "audio:listen"}`)},
							},
						},
						{
							Name: "Image Resource",
							Uris: []string{"/image"},
							Type: imageResourceType,
							Scopes: []apiextensionsv1.JSON{
								{Raw: []byte(`{"name": "image:create"}`)},
								{Raw: []byte(`{"name": "image:read"}`)},
								{Raw: []byte(`{"name": "image:delete"}`)},
							},
						},
					},
					Policies: []keycloakv1alpha1.KeycloakPolicy{
						{
							Name:        "Role Policy",
							Description: "A policy that is role based",
							Type:        "role",
							Logic:       "POSITIVE",
							Config: map[string]string{
								"roles": "[{\"id\":\"" + id + "/uma_protection\",\"required\":true}]",
							},
						},
						{
							Name:             "Aggregate Policy",
							Description:      "A policy that is an aggregate",
							Type:             "aggregate",
							Logic:            "POSITIVE",
							DecisionStrategy: "AFFIRMATIVE",
							Config: map[string]string{
								"applyPolicies": "[\"Role Policy\",\"Deny Policy\"]",
							},
						},
						{
							Name:             "Audio Permission",
							Description:      "An audio permission description",
							Type:             "resource",
							DecisionStrategy: "AFFIRMATIVE",
							Config: map[string]string{
								"defaultResourceType": audioResourceType,
								"default":             "true",
								"applyPolicies":       "[\"Time Policy\"]",
								"scopes":              "[\"audio:listen\"]",
							},
						},
						{
							Name:             "Image Permission",
							Description:      "An image permission description",
							Type:             "scope",
							DecisionStrategy: "UNANIMOUS",
							Config: map[string]string{
								"applyPolicies": "[\"Deny Policy\"]",
								"scopes":        "[\"image:delete\"]",
							},
						},
						{
							Name:        "Deny Policy",
							Description: "A policy that is JS based",
							Type:        "js",
							Config: map[string]string{
								"code": "$evaluation.deny();",
							},
						},
						{
							Name:        "Time Policy",
							Description: "A policy that grants access between 3 and 5 PM",
							Type:        "time",
							Logic:       "POSITIVE",
							Config: map[string]string{
								"hour":    "15",
								"hourEnd": "17",
							},
						},
					},
					Scopes: []keycloakv1alpha1.KeycloakScope{
						{Name: "audio:listen"},
						{Name: "image:create"},
						{Name: "image:read"},
						{Name: "image:delete"},
					},
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

func prepareKeycloakClientAuthZCR(t *testing.T, framework *test.Framework, ctx *test.Context, namespace string) error {
	keycloakClientCR := getKeycloakClientAuthZCR(namespace)
	return Create(framework, keycloakClientCR, ctx)
}

func keycloakClientBasicTest(t *testing.T, framework *test.Framework, ctx *test.Context, namespace string) error {
	return WaitForClientToBeReady(t, framework, namespace, testKeycloakClientCRName)
}

func externalKeycloakClientBasicTest(t *testing.T, framework *test.Framework, ctx *test.Context, namespace string) error {
	return WaitForClientToBeReady(t, framework, namespace, testExternalKeycloakClientCRName)
}

func keycloakClientAuthZTest(t *testing.T, framework *test.Framework, ctx *test.Context, namespace string) error {
	return WaitForClientToBeReady(t, framework, namespace, testAuthZKeycloakClientCRName)
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
	bID, err := getClientRoleID(authenticatedClient, clientName, "b")
	if err != nil {
		return err
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

func getClientRoleID(authenticatedClient common.KeycloakInterface, clientName, roleName string) (string, error) {
	retrievedRoles, err := authenticatedClient.ListClientRoles(clientName, realmName)
	if err != nil {
		return "", err
	}
	return getRole(retrievedRoles, roleName), nil
}

func getRole(retrievedRoles []keycloakv1alpha1.RoleRepresentation, roleName string) string {
	for _, role := range retrievedRoles {
		if role.Name == roleName {
			return role.ID
		}
	}
	return ""
}

func waitForClientRoles(t *testing.T, framework *test.Framework, keycloakCR keycloakv1alpha1.Keycloak, clientCR *keycloakv1alpha1.KeycloakClient, expected []keycloakv1alpha1.RoleRepresentation) error {
	return WaitForConditionWithClient(t, framework, keycloakCR, func(authenticatedClient common.KeycloakInterface) error {
		roles, err := authenticatedClient.ListClientRoles(clientCR.Spec.Client.ID, realmName)
		if err != nil {
			return err
		}

		fail := false
		if len(roles) != len(expected) {
			fail = true
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
					fail = true
					break
				}
			}
		}

		if fail {
			return errors.Errorf("role names not as expected:\n"+
				"expected: %v\n"+
				"actual  : %v",
				expected, roles)
		}
		return nil
	})
}

func prepareKeycloakClientWithRolesCR(t *testing.T, framework *test.Framework, ctx *test.Context, namespace string) error {
	keycloakClientCR := getKeycloakClientCR(namespace, false).DeepCopy()
	keycloakClientCR.Spec.Roles = []keycloakv1alpha1.RoleRepresentation{{Name: "a"}, {Name: "b"}, {Name: "c"}}
	keycloakClientCR.Name = testSecondKeycloakClientCRName
	keycloakClientCR.Spec.Client.ID = secondClientName
	keycloakClientCR.Spec.Client.ClientID = secondClientName
	keycloakClientCR.Spec.Client.Name = secondClientName
	keycloakClientCR.Spec.Client.WebOrigins = []string{"https://operator-test-second.url"}
	return Create(framework, keycloakClientCR, ctx)
}

func getKeycloakClientWithScopeMappingsCR(namespace string, authenticatedClient common.KeycloakInterface,
	realmRoleNames []string, secondClientRoleNames []string) (*keycloakv1alpha1.KeycloakClient, error) {
	client := getKeycloakClientCR(namespace, false)
	mappings, err := getKeycloakClientScopeMappings(authenticatedClient, realmRoleNames, secondClientRoleNames)
	if err != nil {
		return nil, err
	}
	client.Spec.ScopeMappings = mappings
	return client, nil
}

func getKeycloakClientScopeMappings(authenticatedClient common.KeycloakInterface, realmRoleNames []string,
	secondClientRoleNames []string) (*keycloakv1alpha1.MappingsRepresentation, error) {
	var scopeMappings = keycloakv1alpha1.MappingsRepresentation{
		ClientMappings: make(map[string]keycloakv1alpha1.ClientMappingsRepresentation),
	}
	for _, roleName := range realmRoleNames {
		scopeMappings.RealmMappings = append(scopeMappings.RealmMappings, keycloakv1alpha1.RoleRepresentation{
			ID:   roleName,
			Name: roleName,
		})
	}

	secondClient := keycloakv1alpha1.ClientMappingsRepresentation{ID: secondClientName, Client: secondClientName}
	for _, roleName := range secondClientRoleNames {
		roleID, err := getClientRoleID(authenticatedClient, secondClientName, roleName)
		if err != nil {
			return nil, err
		}
		secondClient.Mappings = append(secondClient.Mappings, keycloakv1alpha1.RoleRepresentation{
			ID:   roleID,
			Name: roleName,
		})
	}
	scopeMappings.ClientMappings[secondClientName] = secondClient
	return &scopeMappings, nil
}

func keycloakClientScopeMappingsTest(t *testing.T, framework *test.Framework, ctx *test.Context, namespace string) error {
	err := WaitForClientToBeReady(t, framework, namespace, testSecondKeycloakClientCRName)
	if err != nil {
		return err
	}
	keycloakCR := getDeployedKeycloakCR(framework, namespace)
	authenticatedClient, err := MakeAuthenticatedClient(keycloakCR)
	if err != nil {
		return err
	}

	// create initial client with scope mappings
	client, err := getKeycloakClientWithScopeMappingsCR(
		namespace,
		authenticatedClient,
		[]string{"realmRoleA", "realmRoleB"},
		[]string{"a", "b"})
	if err != nil {
		return err
	}
	err = Create(framework, client, ctx)
	if err != nil {
		return err
	}
	err = WaitForClientToBeReady(t, framework, namespace, testKeycloakClientCRName)
	if err != nil {
		return err
	}

	// add non-existent roles
	var retrievedClient keycloakv1alpha1.KeycloakClient
	err = GetNamespacedObject(framework, namespace, testKeycloakClientCRName, &retrievedClient)
	if err != nil {
		return err
	}
	mappings, err := getKeycloakClientScopeMappings(
		authenticatedClient,
		[]string{"realmRoleB", "realmRoleC", "nonexistent"},
		[]string{"b", "c", "nonexistent"},
	)
	if err != nil {
		return err
	}
	retrievedClient.Spec.ScopeMappings = mappings
	err = Update(framework, &retrievedClient)
	if err != nil {
		return err
	}
	err = WaitForClientToBeFailing(t, framework, namespace, testKeycloakClientCRName)
	if err != nil {
		return err
	}

	// update client: delete/leave/create mappings
	err = GetNamespacedObject(framework, namespace, testKeycloakClientCRName, &retrievedClient)
	if err != nil {
		return err
	}
	mappings, err = getKeycloakClientScopeMappings(authenticatedClient, []string{"realmRoleB", "realmRoleC"}, []string{"b", "c"})
	if err != nil {
		return err
	}
	retrievedClient.Spec.ScopeMappings = mappings
	err = Update(framework, &retrievedClient)
	if err != nil {
		return err
	}
	err = WaitForClientToBeReady(t, framework, namespace, testKeycloakClientCRName)
	if err != nil {
		return err
	}

	retrievedMappings, err := authenticatedClient.ListScopeMappings(clientName, realmName)
	if err != nil {
		return err
	}
	expected := retrievedClient.Spec.ScopeMappings

	difference, intersection := model.RoleDifferenceIntersection(
		retrievedMappings.RealmMappings,
		expected.RealmMappings)
	assert.Equal(t, 0, len(difference))
	assert.Equal(t, 2, len(intersection))

	difference, intersection = model.RoleDifferenceIntersection(
		retrievedMappings.ClientMappings[secondClientName].Mappings,
		expected.ClientMappings[secondClientName].Mappings)
	assert.Equal(t, 0, len(difference))
	assert.Equal(t, 2, len(intersection))

	return nil
}
