package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KeycloakRealmSpec defines the desired state of KeycloakRealm
// +k8s:openapi-gen=true
type KeycloakRealmSpec struct {
	InstanceSelector *metav1.LabelSelector `json:"instanceSelector,omitempty"`
	Realm            *KeycloakAPIRealm     `json:"realm"`
	// +listType=map
	RealmOverrides []*RedirectorIdentityProviderOverride `json:"realmOverrides,omitempty"`
}

type KeycloakAPIRealm struct {
	ID                string                      `json:"id"`
	Realm             string                      `json:"realm"`
	Enabled           bool                        `json:"enabled"`
	DisplayName       string                      `json:"displayName"`
	Users             []*KeycloakAPIUser          `json:"users,omitempty"`
	Clients           []*KeycloakAPIClient        `json:"clients,omitempty"`
	IdentityProviders []*KeycloakIdentityProvider `json:"identityProviders,omitempty"`
	EventsListeners   []string                    `json:"eventsListeners,omitempty"`
}

type RedirectorIdentityProviderOverride struct {
	IdentityProvider string `json:"identityProvider,omitempty"`
	ForFlow          string `json:"forFlow,omitempty"`
}

type KeycloakIdentityProvider struct {
	Alias                     string            `json:"alias,omitempty"`
	DisplayName               string            `json:"displayName,omitempty"`
	InternalID                string            `json:"internalId,omitempty"`
	ProviderID                string            `json:"providerId,omitempty"`
	Enabled                   bool              `json:"enabled,omitempty"`
	TrustEmail                bool              `json:"trustEmail,omitempty"`
	StoreToken                bool              `json:"storeToken,omitempty"`
	AddReadTokenRoleOnCreate  bool              `json:"addReadTokenRoleOnCreate,omitempty"`
	FirstBrokerLoginFlowAlias string            `json:"firstBrokerLoginFlowAlias,omitempty"`
	PostBrokerLoginFlowAlias  string            `json:"postBrokerLoginFlowAlias,omitempty"`
	LinkOnly                  bool              `json:"linkOnly,omitempty"`
	Config                    map[string]string `json:"config,omitempty"`
}

type KeycloakAPIUser struct {
	ID                  string               `json:"id,omitempty"`
	UserName            string               `json:"username,omitempty"`
	FirstName           string               `json:"firstName,omitempty"`
	LastName            string               `json:"lastName,omitempty"`
	Email               string               `json:"email,omitempty"`
	EmailVerified       bool                 `json:"emailVerified,omitempty"`
	Enabled             bool                 `json:"enabled,omitempty"`
	RealmRoles          []string             `json:"realmRoles,omitempty"`
	ClientRoles         map[string][]string  `json:"clientRoles,omitempty"`
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
	// Current phase of the operator.
	Phase StatusPhase `json:"phase"`
	// Human-readable message indicating details about current operator phase or error.
	Message string `json:"message"`
	// True if all resources are in a ready state and all work is done.
	Ready bool `json:"ready"`
	// A map of all the secondary resources types and names created for this CR. e.g "Deployment": [ "DeploymentName1", "DeploymentName2" ]
	SecondaryResources map[string][]string `json:"secondaryResources,omitempty"`
	// TODO
	LoginURL string `json:"loginURL"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KeycloakRealm is the Schema for the keycloakrealms API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=keycloakrealms,scope=Namespaced
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

func (i *KeycloakRealm) UpdateStatusSecondaryResources(kind string, resourceName string) {
	i.Status.SecondaryResources = UpdateStatusSecondaryResources(i.Status.SecondaryResources, kind, resourceName)
}
