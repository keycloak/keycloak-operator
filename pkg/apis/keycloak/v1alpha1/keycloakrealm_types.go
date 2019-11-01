package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KeycloakRealmSpec defines the desired state of KeycloakRealm
// +k8s:openapi-gen=true
type KeycloakRealmSpec struct {
	InstanceSelector                  *metav1.LabelSelector `json:"instanceSelector,omitempty"`
	BrowserRedirectorIdentityProvider string                `json:"browserRedirectorIdentityProvider,omitempty"`
	Realm                             *KeycloakAPIRealm     `json:"realm"`
}

type KeycloakAPIRealm struct {
	ID                string                      `json:"id,omitempty"`
	Realm             string                      `json:"realm,omitempty"`
	Enabled           bool                        `json:"enabled"`
	DisplayName       string                      `json:"displayName"`
	Users             []*KeycloakAPIUser          `json:"users,omitempty"`
	Clients           []*KeycloakAPIClient        `json:"clients,omitempty"`
	IdentityProviders []*KeycloakIdentityProvider `json:"identityProviders,omitempty"`
	EventsListeners   []string                    `json:"eventsListeners,omitempty"`
}

type KeycloakIdentityProvider struct {
	Alias                     string            `json:"alias,omitempty"`
	DisplayName               string            `json:"displayName"`
	InternalID                string            `json:"internalId,omitempty"`
	ProviderID                string            `json:"providerId,omitempty"`
	Enabled                   bool              `json:"enabled"`
	TrustEmail                bool              `json:"trustEmail"`
	StoreToken                bool              `json:"storeToken"`
	AddReadTokenRoleOnCreate  bool              `json:"addReadTokenRoleOnCreate"`
	FirstBrokerLoginFlowAlias string            `json:"firstBrokerLoginFlowAlias"`
	PostBrokerLoginFlowAlias  string            `json:"postBrokerLoginFlowAlias"`
	LinkOnly                  bool              `json:"linkOnly"`
	Config                    map[string]string `json:"config"`
}

type KeycloakAPIUser struct {
	ID                  string               `json:"id,omitempty"`
	UserName            string               `json:"username,omitempty"`
	FirstName           string               `json:"firstName"`
	LastName            string               `json:"lastName"`
	Email               string               `json:"email,omitempty"`
	EmailVerified       bool                 `json:"emailVerified"`
	Enabled             bool                 `json:"enabled"`
	RealmRoles          []string             `json:"realmRoles,omitempty"`
	ClientRoles         map[string][]string  `json:"clientRoles"`
	RequiredActions     []string             `json:"requiredActions,omitempty"`
	Groups              []string             `json:"groups,omitempty"`
	FederatedIdentities []FederatedIdentity  `json:"federatedIdentities,omitempty"`
	Credentials         []KeycloakCredential `json:"credentials,omitempty"`
}

type FederatedIdentity struct {
	IdentityProvider string `json:"identityProvider,omitempty"`
	UserID           string `json:"userId,omitempty"`
	UserName         string `json:"userName,omitempty"`
}

type KeycloakCredential struct {
	Type      string `json:"type,omitempty"`
	Value     string `json:"value,omitempty"`
	Temporary bool   `json:"temporary,omitempty"`
}

type KeycloakAPIClient struct {
	ID                        string                   `json:"id,omitempty"`
	ClientID                  string                   `json:"clientId,omitempty"`
	Secret                    string                   `json:"secret"`
	Name                      string                   `json:"name"`
	BaseURL                   string                   `json:"baseUrl"`
	AdminURL                  string                   `json:"adminUrl"`
	RootURL                   string                   `json:"rootUrl"`
	Description               string                   `json:"description"`
	SurrogateAuthRequired     bool                     `json:"surrogateAuthRequired"`
	Enabled                   bool                     `json:"enabled"`
	ClientAuthenticatorType   string                   `json:"clientAuthenticatorType"`
	DefaultRoles              []string                 `json:"defaultRoles,omitempty"`
	RedirectUris              []string                 `json:"redirectUris,omitempty"`
	WebOrigins                []string                 `json:"webOrigins,omitempty"`
	NotBefore                 int                      `json:"notBefore"`
	BearerOnly                bool                     `json:"bearerOnly"`
	ConsentRequired           bool                     `json:"consentRequired"`
	StandardFlowEnabled       bool                     `json:"standardFlowEnabled"`
	ImplicitFlowEnabled       bool                     `json:"implicitFlowEnabled"`
	DirectAccessGrantsEnabled bool                     `json:"directAccessGrantsEnabled"`
	ServiceAccountsEnabled    bool                     `json:"serviceAccountsEnabled"`
	PublicClient              bool                     `json:"publicClient"`
	FrontchannelLogout        bool                     `json:"frontchannelLogout"`
	Protocol                  string                   `json:"protocol,omitempty"`
	Attributes                map[string]string        `json:"attributes,omitempty"`
	FullScopeAllowed          bool                     `json:"fullScopeAllowed"`
	NodeReRegistrationTimeout int                      `json:"nodeReRegistrationTimeout"`
	ProtocolMappers           []KeycloakProtocolMapper `json:"protocolMappers,omitempty"`
	UseTemplateConfig         bool                     `json:"useTemplateConfig"`
	UseTemplateScope          bool                     `json:"useTemplateScope"`
	UseTemplateMappers        bool                     `json:"useTemplateMappers"`
	Access                    map[string]bool          `json:"access"`
}

type KeycloakProtocolMapper struct {
	ID              string            `json:"id,omitempty"`
	Name            string            `json:"name,omitempty"`
	Protocol        string            `json:"protocol,omitempty"`
	ProtocolMapper  string            `json:"protocolMapper,omitempty"`
	ConsentRequired bool              `json:"consentRequired,omitempty"`
	ConsentText     string            `json:"consentText"`
	Config          map[string]string `json:"config"`
}

type KeycloakUserRole struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Composite   bool   `json:"composite,omitempty"`
	ClientRole  bool   `json:"clientRole,omitempty"`
	ContainerID string `json:"containerId,omitempty"`
}

type AuthenticatorConfig struct {
	Alias  string            `json:"alias,omitempty"`
	Config map[string]string `json:"config,omitempty"`
	ID     string            `json:"id,omitempty"`
}

type KeycloakAPIPasswordReset struct {
	Type      string `json:"type"`
	Value     string `json:"value"`
	Temporary bool   `json:"temporary"`
}

type AuthenticationExecutionInfo struct {
	Alias                string   `json:"alias,omitempty"`
	AuthenticationConfig string   `json:"authenticationConfig,omitempty"`
	AuthenticationFlow   bool     `json:"authenticationFlow,omitempty"`
	Configurable         bool     `json:"configurable,omitempty"`
	DisplayName          string   `json:"displayName,omitempty"`
	FlowID               string   `json:"flowId,omitempty"`
	ID                   string   `json:"id,omitempty"`
	Index                int32    `json:"index,omitempty"`
	Level                int32    `json:"level,omitempty"`
	ProviderID           string   `json:"providerId,omitempty"`
	Requirement          string   `json:"requirement,omitempty"`
	RequirementChoices   []string `json:"requirementChoices,omitempty"`
}

type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// KeycloakRealmStatus defines the observed state of KeycloakRealm
// +k8s:openapi-gen=true
type KeycloakRealmStatus struct {
	Phase    string `json:"phase"`
	LoginURL string `json:"loginUrl"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KeycloakRealm is the Schema for the keycloakrealms API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type KeycloakRealm struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KeycloakRealmSpec   `json:"spec,omitempty"`
	Status KeycloakRealmStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KeycloakRealmList contains a list of KeycloakRealm
type KeycloakRealmList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KeycloakRealm `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KeycloakRealm{}, &KeycloakRealmList{})
}
