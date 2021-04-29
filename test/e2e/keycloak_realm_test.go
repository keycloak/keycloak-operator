package e2e

import (
	"crypto/tls"
	"net/http"
	"testing"

	keycloakv1alpha1 "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/test"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	realmName                  = "test-realm"
	testOperatorIDPDisplayName = "Test Operator IDP"
)

func NewKeycloakRealmsCRDTestStruct() *CRDTestStruct {
	return &CRDTestStruct{
		prepareEnvironmentSteps: []environmentInitializationStep{
			prepareKeycloaksCR,
		},
		testSteps: map[string]deployedOperatorTestStep{
			"keycloakRealmBasicTest": {
				prepareTestEnvironmentSteps: []environmentInitializationStep{
					prepareKeycloakRealmCR,
				},
				testFunction: keycloakRealmBasicTest,
			},
			"keycloakRealmWithIdentityProviderTest": {
				testFunction: keycloakRealmWithIdentityProviderTest,
			},
			"keycloakRealmWithClientScopesTest": {
				testFunction: keycloakRealmWithClientScopesTest,
			},
			"keycloakRealmWithAuthenticatorFlowTest": {
				testFunction: keycloakRealmWithAuthenticatorFlowTest,
			},
			"keycloakRealmWithUserFederationTest": {
				testFunction: keycloakRealmWithUserFederationTest,
			},
			"unmanagedKeycloakRealmTest": {
				testFunction: keycloakUnmanagedRealmTest,
			},
		},
	}
}

func getKeycloakRealmCR(namespace string) *keycloakv1alpha1.KeycloakRealm {
	return &keycloakv1alpha1.KeycloakRealm{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testKeycloakRealmCRName,
			Namespace: namespace,
			Labels:    CreateLabel(namespace),
		},
		Spec: keycloakv1alpha1.KeycloakRealmSpec{
			InstanceSelector: &metav1.LabelSelector{
				MatchLabels: CreateLabel(namespace),
			},
			Realm: &keycloakv1alpha1.KeycloakAPIRealm{
				ID:                                 realmName,
				Realm:                              realmName,
				Enabled:                            true,
				DisplayName:                        "Operator Testing Realm",
				DisplayNameHTML:                    "<div class='kc-logo-text'><span>Operator Testing Realm</span></div>",
				PasswordPolicy:                     "lowerCase(1)",
				BruteForceProtected:                &[]bool{true}[0],
				PermanentLockout:                   &[]bool{false}[0],
				FailureFactor:                      &[]int32{30}[0],
				WaitIncrementSeconds:               &[]int32{60}[0],
				QuickLoginCheckMilliSeconds:        &[]int64{1000}[0],
				MinimumQuickLoginWaitSeconds:       &[]int32{60}[0],
				MaxFailureWaitSeconds:              &[]int32{900}[0],
				MaxDeltaTimeSeconds:                &[]int32{43200}[0],
				AccessTokenLifespanForImplicitFlow: &[]int32{3600}[0],
				AccessTokenLifespan:                &[]int32{4800}[0],
				SMTPServer: map[string]string{
					"starttls":        "",
					"auth":            "",
					"host":            "smtp.server",
					"from":            "sso@example.com",
					"fromDisplayName": "Example Company",
					"envelopeFrom":    "sso@example.com",
					"ssl":             "",
				},
				InternationalizationEnabled: &[]bool{true}[0],
				SupportedLocales:            []string{"en", "de"},
				DefaultLocale:               "en",
				LoginTheme:                  "keycloak",
				AccountTheme:                "keycloak",
				AdminTheme:                  "keycloak",
				EmailTheme:                  "keycloak",
			},
		},
	}
}

func prepareKeycloakRealmCR(t *testing.T, framework *test.Framework, ctx *test.Context, namespace string) error {
	keycloakRealmCR := getKeycloakRealmCR(namespace)
	return Create(framework, keycloakRealmCR, ctx)
}

func keycloakRealmBasicTest(t *testing.T, framework *test.Framework, ctx *test.Context, namespace string) error {
	return WaitForRealmToBeReady(t, framework, namespace)
}

func keycloakRealmWithIdentityProviderTest(t *testing.T, framework *test.Framework, ctx *test.Context, namespace string) error {
	keycloakRealmCR := getKeycloakRealmCR(namespace)

	identityProvider := &keycloakv1alpha1.KeycloakIdentityProvider{
		Alias:                     "oidc",
		DisplayName:               testOperatorIDPDisplayName,
		InternalID:                "",
		ProviderID:                "oidc",
		Enabled:                   true,
		TrustEmail:                false,
		StoreToken:                false,
		AddReadTokenRoleOnCreate:  false,
		FirstBrokerLoginFlowAlias: "first broker login",
		PostBrokerLoginFlowAlias:  "",
		LinkOnly:                  false,
		Config: map[string]string{
			"useJwksUrl":       "true",
			"loginHint":        "",
			"authorizationUrl": "https://operator.test.url/authorization_url",
			"tokenUrl":         "https://operator.test.url/token_url",
			"clientAuthMethod": "client_secret_jwt",
			"clientId":         "operator-idp",
			"clientSecret":     "test",
			"allowedClockSkew": "5",
		},
	}

	keycloakRealmCR.Spec.Realm.IdentityProviders = []*keycloakv1alpha1.KeycloakIdentityProvider{identityProvider}

	err := Create(framework, keycloakRealmCR, ctx)
	if err != nil {
		return err
	}

	err = WaitForRealmToBeReady(t, framework, namespace)
	if err != nil {
		return err
	}

	keycloakCR := getDeployedKeycloakCR(framework, namespace)
	keycloakURL := keycloakCR.Status.ExternalURL

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint
	return WaitForSuccessResponseToContain(t, framework, keycloakURL+"/auth/realms/"+realmName+"/account", testOperatorIDPDisplayName)
}

func keycloakRealmWithClientScopesTest(t *testing.T, framework *test.Framework, ctx *test.Context, namespace string) error {
	keycloakRealmCR := getKeycloakRealmCR(namespace)

	identityProvider := &keycloakv1alpha1.KeycloakIdentityProvider{
		Alias:                     "oidc",
		DisplayName:               testOperatorIDPDisplayName,
		InternalID:                "",
		ProviderID:                "oidc",
		Enabled:                   true,
		TrustEmail:                false,
		StoreToken:                false,
		AddReadTokenRoleOnCreate:  false,
		FirstBrokerLoginFlowAlias: "first broker login",
		PostBrokerLoginFlowAlias:  "",
		LinkOnly:                  false,
		Config: map[string]string{
			"useJwksUrl":       "true",
			"loginHint":        "",
			"authorizationUrl": "https://operator.test.url/authorization_url",
			"tokenUrl":         "https://operator.test.url/token_url",
			"clientAuthMethod": "client_secret_jwt",
			"clientId":         "operator-idp",
			"clientSecret":     "test",
			"allowedClockSkew": "5",
		},
	}

	keycloakRealmCR.Spec.Realm.IdentityProviders = []*keycloakv1alpha1.KeycloakIdentityProvider{identityProvider}
	keycloakRealmCR.Spec.Realm.ClientScopes = []keycloakv1alpha1.KeycloakClientScope{
		{
			Name:        "profile",
			Description: "subset of the built in profile scope, for e2e testing",
			Protocol:    "openid-connect",
			Attributes: map[string]string{
				"include.in.token.scope":    "true",
				"display.on.consent.screen": "false",
			},
			ProtocolMappers: []keycloakv1alpha1.KeycloakProtocolMapper{
				{
					Name:           "family name",
					Protocol:       "openid-connect",
					ProtocolMapper: "oidc-usermodel-property-mapper",
					Config: map[string]string{
						"access.token.claim":   "true",
						"claim.name":           "family_name",
						"id.token.claim":       "true",
						"jsonType.label":       "String",
						"user.attribute":       "lastName",
						"userinfo.token.claim": "true",
					},
					ConsentRequired: false,
				},
				{
					Name:           "given name",
					Protocol:       "openid-connect",
					ProtocolMapper: "oidc-usermodel-property-mapper",
					Config: map[string]string{
						"access.token.claim":   "true",
						"claim.name":           "given_name",
						"id.token.claim":       "true",
						"jsonType.label":       "String",
						"user.attribute":       "firstName",
						"userinfo.token.claim": "true",
					},
					ConsentRequired: false,
				},
				{
					Name:           "username",
					Protocol:       "openid-connect",
					ProtocolMapper: "oidc-usermodel-property-mapper",
					Config: map[string]string{
						"access.token.claim":   "true",
						"claim.name":           "preferred_username",
						"id.token.claim":       "true",
						"jsonType.label":       "String",
						"user.attribute":       "username",
						"userinfo.token.claim": "true",
					},
					ConsentRequired: false,
				},
			},
		},
		{
			Name:     "groups",
			Protocol: "openid-connect",
			Attributes: map[string]string{
				"include.in.token.scope":    "true",
				"display.on.consent.screen": "false",
			},
			ProtocolMappers: []keycloakv1alpha1.KeycloakProtocolMapper{
				{
					Name:            "groups",
					Protocol:        "openid-connect",
					ProtocolMapper:  "oidc-group-membership-mapper",
					ConsentRequired: false,
					Config: map[string]string{
						"full.path":            "false",
						"id.token.claim":       "true",
						"access.token.claim":   "true",
						"claim.name":           "groups",
						"userinfo.token.claim": "true",
					},
				},
			},
		},
	}

	err := Create(framework, keycloakRealmCR, ctx)
	if err != nil {
		return err
	}

	err = WaitForRealmToBeReady(t, framework, namespace)
	if err != nil {
		return err
	}

	keycloakCR := getDeployedKeycloakCR(framework, namespace)
	keycloakURL := keycloakCR.Status.ExternalURL

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint
	return WaitForSuccessResponseToContain(t, framework, keycloakURL+"/auth/realms/"+realmName+"/account", testOperatorIDPDisplayName)
}

// These flows (by name, not the exact contents here) are built in and required to exist
// See https://issues.redhat.com/browse/KEYCLOAK-14779
func getBrowserFlow() keycloakv1alpha1.KeycloakAPIAuthenticationFlow {
	return keycloakv1alpha1.KeycloakAPIAuthenticationFlow{
		Alias:      "browser",
		ProviderID: "basic-flow",
		TopLevel:   true,
		AuthenticationExecutions: []keycloakv1alpha1.KeycloakAPIAuthenticationExecution{
			{
				Authenticator: "auth-username-password-form",
				Requirement:   "REQUIRED",
			},
		},
	}
}

func getRegistrationFlow() keycloakv1alpha1.KeycloakAPIAuthenticationFlow {
	return keycloakv1alpha1.KeycloakAPIAuthenticationFlow{
		ID:         "d6a87b0e-dfe1-495b-af73-a056f8734b4d",
		Alias:      "registration",
		ProviderID: "basic-flow",
		TopLevel:   true,
		AuthenticationExecutions: []keycloakv1alpha1.KeycloakAPIAuthenticationExecution{
			{
				Authenticator:       "identity-provider-redirector",
				AuthenticatorConfig: "oidc",
				Requirement:         "ALTERNATIVE",
			},
		},
	}
}

func getDirectGrantFlow() keycloakv1alpha1.KeycloakAPIAuthenticationFlow {
	return keycloakv1alpha1.KeycloakAPIAuthenticationFlow{
		Alias:                    "direct grant",
		ProviderID:               "basic-flow",
		TopLevel:                 true,
		AuthenticationExecutions: []keycloakv1alpha1.KeycloakAPIAuthenticationExecution{},
	}
}

func keycloakRealmWithAuthenticatorFlowTest(t *testing.T, framework *test.Framework, ctx *test.Context, namespace string) error {
	keycloakRealmCR := getKeycloakRealmCR(namespace)

	identityProvider := &keycloakv1alpha1.KeycloakIdentityProvider{
		Alias:                     "oidc",
		DisplayName:               testOperatorIDPDisplayName,
		InternalID:                "",
		ProviderID:                "oidc",
		Enabled:                   true,
		TrustEmail:                false,
		StoreToken:                false,
		AddReadTokenRoleOnCreate:  false,
		FirstBrokerLoginFlowAlias: "first broker login",
		PostBrokerLoginFlowAlias:  "",
		LinkOnly:                  false,
		Config: map[string]string{
			"useJwksUrl":       "true",
			"loginHint":        "",
			"authorizationUrl": "https://operator.test.url/authorization_url",
			"tokenUrl":         "https://operator.test.url/token_url",
			"clientAuthMethod": "client_secret_jwt",
			"clientId":         "operator-idp",
			"clientSecret":     "test",
			"allowedClockSkew": "5",
		},
	}

	keycloakRealmCR.Spec.Realm.AuthenticatorConfig = []keycloakv1alpha1.KeycloakAPIAuthenticatorConfig{
		{
			ID:    "ffe3bf1a-5ef0-41af-96b5-c02543dd787a",
			Alias: "oidc",
			Config: map[string]string{
				"defaultProvider": "oidc",
			},
		},
	}

	var autoLinkFlow = keycloakv1alpha1.KeycloakAPIAuthenticationFlow{
		Alias:      "Auto Link",
		ProviderID: "basic-flow",
		TopLevel:   true,
		AuthenticationExecutions: []keycloakv1alpha1.KeycloakAPIAuthenticationExecution{
			{
				Authenticator: "idp-create-user-if-unique",
				Requirement:   "ALTERNATIVE",
				Priority:      0,
			},
			{
				Authenticator: "idp-auto-link",
				Requirement:   "ALTERNATIVE",
				Priority:      1,
			},
		},
	}

	keycloakRealmCR.Spec.Realm.AuthenticationFlows = []keycloakv1alpha1.KeycloakAPIAuthenticationFlow{autoLinkFlow, getBrowserFlow(), getRegistrationFlow(), getDirectGrantFlow()}

	keycloakRealmCR.Spec.Realm.IdentityProviders = []*keycloakv1alpha1.KeycloakIdentityProvider{identityProvider}

	err := Create(framework, keycloakRealmCR, ctx)
	if err != nil {
		return err
	}

	err = WaitForRealmToBeReady(t, framework, namespace)
	if err != nil {
		return err
	}

	keycloakCR := getDeployedKeycloakCR(framework, namespace)
	keycloakURL := keycloakCR.Status.ExternalURL

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint
	return WaitForSuccessResponseToContain(t, framework, keycloakURL+"/auth/realms/"+realmName+"/account", testOperatorIDPDisplayName)
}

func keycloakRealmWithUserFederationTest(t *testing.T, framework *test.Framework, ctx *test.Context, namespace string) error {
	keycloakRealmCR := getKeycloakRealmCR(namespace)

	identityProvider := &keycloakv1alpha1.KeycloakIdentityProvider{
		Alias:                     "oidc",
		DisplayName:               testOperatorIDPDisplayName,
		InternalID:                "",
		ProviderID:                "oidc",
		Enabled:                   true,
		TrustEmail:                false,
		StoreToken:                false,
		AddReadTokenRoleOnCreate:  false,
		FirstBrokerLoginFlowAlias: "first broker login",
		PostBrokerLoginFlowAlias:  "",
		LinkOnly:                  false,
		Config: map[string]string{
			"useJwksUrl":       "true",
			"loginHint":        "",
			"authorizationUrl": "https://operator.test.url/authorization_url",
			"tokenUrl":         "https://operator.test.url/token_url",
			"clientAuthMethod": "client_secret_jwt",
			"clientId":         "operator-idp",
			"clientSecret":     "test",
			"allowedClockSkew": "5",
		},
	}

	userFederationMapper := keycloakv1alpha1.KeycloakAPIUserFederationMapper{
		Config: map[string]string{
			"groups.ldap.filter":                   "(|(CN=Role-*)(CN=Access-*))",
			"groups.dn":                            "OU=groups,DC=example,DC=com",
			"mode":                                 "READ_ONLY",
			"preserve.group.inheritance":           "false",
			"ignore.missing.groups":                "false",
			"group.name.ldap.attribute":            "cn",
			"drop.non.existing.groups.during.sync": "false",
			"user.roles.retrieve.strategy":         "LOAD_GROUPS_BY_MEMBER_ATTRIBUTE_RECURSIVELY",
		},
		Name:                          "group-mapper",
		FederationMapperType:          "group-ldap-mapper",
		FederationProviderDisplayName: "ldap-provider",
	}

	userFederationProvider := keycloakv1alpha1.KeycloakAPIUserFederationProvider{
		Config: map[string]string{
			"vendor":           "ad",
			"connectionUrl":    "ldap://127.0.0.1",
			"bindDn":           "foo",
			"bindCredential":   "p@ssword",
			"useTruststoreSpi": "ldapsOnly",
			"editMode":         "READ_ONLY",
		},
		DisplayName:  "ldap-provider",
		ProviderName: "ldap",
	}

	keycloakRealmCR.Spec.Realm.UserFederationMappers = []keycloakv1alpha1.KeycloakAPIUserFederationMapper{userFederationMapper}
	keycloakRealmCR.Spec.Realm.UserFederationProviders = []keycloakv1alpha1.KeycloakAPIUserFederationProvider{userFederationProvider}
	keycloakRealmCR.Spec.Realm.IdentityProviders = []*keycloakv1alpha1.KeycloakIdentityProvider{identityProvider}

	err := Create(framework, keycloakRealmCR, ctx)
	if err != nil {
		return err
	}

	err = WaitForRealmToBeReady(t, framework, namespace)
	if err != nil {
		return err
	}

	keycloakCR := getDeployedKeycloakCR(framework, namespace)
	keycloakURL := keycloakCR.Status.ExternalURL

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint
	return WaitForSuccessResponseToContain(t, framework, keycloakURL+"/auth/realms/"+realmName+"/account", testOperatorIDPDisplayName)
}

func keycloakUnmanagedRealmTest(t *testing.T, framework *test.Framework, ctx *test.Context, namespace string) error {
	keycloakRealmCR := getKeycloakRealmCR(namespace)
	keycloakRealmCR.Spec.Unmanaged = true

	err := Create(framework, keycloakRealmCR, ctx)
	if err != nil {
		return err
	}

	err = WaitForRealmToBeReady(t, framework, namespace)
	if err != nil {
		return err
	}

	return err
}
