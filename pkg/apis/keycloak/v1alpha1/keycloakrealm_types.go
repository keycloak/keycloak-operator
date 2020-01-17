package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KeycloakRealmSpec defines the desired state of KeycloakRealm
// +k8s:openapi-gen=true
type KeycloakRealmSpec struct {
	// Selector for looking up Keycloak Custom Resources.
	// +kubebuilder:validation:Required
	InstanceSelector *metav1.LabelSelector `json:"instanceSelector,omitempty"`
	// Keycloak Realm REST object.
	// +kubebuilder:validation:Required
	Realm *KeycloakAPIRealm `json:"realm"`
	// A list of overrides to the default Realm behavior.
	// +listType=map
	RealmOverrides []*RedirectorIdentityProviderOverride `json:"realmOverrides,omitempty"`
}

type KeycloakAPIRealm struct {
	// +kubebuilder:validation:Required
	// +optional
	ID string `json:"id"`
	// Realm name.
	// +kubebuilder:validation:Required
	Realm string `json:"realm"`
	// Realm enabled flag.
	// +optional
	Enabled bool `json:"enabled"`
	// Realm display name.
	// +optional
	DisplayName string `json:"displayName"`
	// A set of Keycloak Users.
	// +optional
	Users []*KeycloakAPIUser `json:"users,omitempty"`
	// A set of Keycloak Clients.
	// +optional
	Clients []*KeycloakAPIClient `json:"clients,omitempty"`
	// A set of Identity Providers.
	// +optional
	IdentityProviders []*KeycloakIdentityProvider `json:"identityProviders,omitempty"`
	// A set of Event Listeners.
	// +optional
	EventsListeners []string `json:"eventsListeners,omitempty"`
	// Enable events recording
	// TODO: change to values and use kubebuilder default annotation once supported
	// +optional
	EventsEnabled *bool `json:"eventsEnabled,omitempty"`
	// Enable events recording
	// TODO: change to values and use kubebuilder default annotation once supported
	// +optional
	AdminEventsEnabled *bool `json:"adminEventsEnabled,omitempty"`
	// Enable admin events details
	// TODO: change to values and use kubebuilder default annotation once supported
	// +optional
	AdminEventsDetailsEnabled *bool `json:"adminEventsDetailsEnabled,omitempty"`
}

type RedirectorIdentityProviderOverride struct {
	// Identity Provider to be overridden.
	// +optional
	IdentityProvider string `json:"identityProvider,omitempty"`
	// Flow to be overridden.
	// +optional
	ForFlow string `json:"forFlow,omitempty"`
}

type KeycloakIdentityProvider struct {
	// Identity Provider Alias.
	// +optional
	Alias string `json:"alias,omitempty"`
	// Identity Provider Display Name.
	// +optional
	DisplayName string `json:"displayName,omitempty"`
	// Identity Provider Internal ID.
	// +optional
	InternalID string `json:"internalId,omitempty"`
	// Identity Provider ID.
	// +optional
	ProviderID string `json:"providerId,omitempty"`
	// Identity Provider enabled flag.
	// +optional
	Enabled bool `json:"enabled,omitempty"`
	// Identity Provider Trust Email.
	// +optional
	TrustEmail bool `json:"trustEmail,omitempty"`
	// Identity Provider Store to Token.
	// +optional
	StoreToken bool `json:"storeToken,omitempty"`
	// Adds Read Token role when creating this Identity Provider.
	// +optional
	AddReadTokenRoleOnCreate bool `json:"addReadTokenRoleOnCreate,omitempty"`
	// Identity Provider First Broker Login Flow Alias.
	// +optional
	FirstBrokerLoginFlowAlias string `json:"firstBrokerLoginFlowAlias,omitempty"`
	// Identity Provider Post Broker Login Flow Alias.
	// +optional
	PostBrokerLoginFlowAlias string `json:"postBrokerLoginFlowAlias,omitempty"`
	// Identity Provider Link Only setting.
	// +optional
	LinkOnly bool `json:"linkOnly,omitempty"`
	// Identity Provider config.
	// +optional
	Config map[string]string `json:"config,omitempty"`
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
	// Authenticator Config Alias.
	// +optional
	Alias string `json:"alias,omitempty"`
	// Authenticator config.
	// +optional
	Config map[string]string `json:"config,omitempty"`
	// Authenticator ID.
	// +optional
	ID string `json:"id,omitempty"`
}

type KeycloakAPIPasswordReset struct {
	// Password Reset Type.
	// +optional
	Type string `json:"type"`
	// Password Reset Value.
	// +optional
	Value string `json:"value"`
	// True if this Password Reset object is temporary.
	// +optional
	Temporary bool `json:"temporary"`
}

type AuthenticationExecutionInfo struct {
	// Authentication Execution Info Alias.
	// +optional
	Alias string `json:"alias,omitempty"`
	// Authentication Execution Info Config.
	// +optional
	AuthenticationConfig string `json:"authenticationConfig,omitempty"`
	// True if Authentication Flow is enabled.
	// +optional
	AuthenticationFlow bool `json:"authenticationFlow,omitempty"`
	// True if Authentication Execution Info is configurable.
	// +optional
	Configurable bool `json:"configurable,omitempty"`
	// Authentication Execution Info Display Name.
	// +optional
	DisplayName string `json:"displayName,omitempty"`
	// Authentication Execution Info Flow ID.
	// +optional
	FlowID string `json:"flowId,omitempty"`
	// Authentication Execution Info ID.
	// +optional
	ID string `json:"id,omitempty"`
	// Authentication Execution Info Index.
	// +optional
	Index int32 `json:"index,omitempty"`
	// Authentication Execution Info Level.
	// +optional
	Level int32 `json:"level,omitempty"`
	// Authentication Execution Info Provider ID.
	// +optional
	ProviderID string `json:"providerId,omitempty"`
	// Authentication Execution Info Requirement.
	// +optional
	Requirement string `json:"requirement,omitempty"`
	// Authentication Execution Info Requirement Choices.
	// +optional
	RequirementChoices []string `json:"requirementChoices,omitempty"`
}

type TokenResponse struct {
	// Token Response Access Token.
	// +optional
	AccessToken string `json:"access_token"`
	// Token Response Expired In setting.
	// +optional
	ExpiresIn int `json:"expires_in"`
	// Token Response Refresh Expires In setting.
	// +optional
	RefreshExpiresIn int `json:"refresh_expires_in"`
	// Token Response Refresh Token.
	// +optional
	RefreshToken string `json:"refresh_token"`
	// Token Response Token Type.
	// +optional
	TokenType string `json:"token_type"`
	// Token Response Not Before Policy setting.
	// +optional
	NotBeforePolicy int `json:"not-before-policy"`
	// Token Response Session State.
	// +optional
	SessionState string `json:"session_state"`
	// Token Response Error.
	// +optional
	Error string `json:"error"`
	// Token Response Error Description.
	// +optional
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
