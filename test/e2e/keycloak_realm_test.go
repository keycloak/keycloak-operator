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
				ID:          realmName,
				Realm:       realmName,
				Enabled:     true,
				DisplayName: "Operator Testing Realm",
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
	keycloakInternalURL := keycloakCR.Status.InternalURL

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint
	return WaitForSuccessResponseToContain(t, framework, keycloakInternalURL+"/auth/realms/"+realmName+"/account", testOperatorIDPDisplayName)
}
