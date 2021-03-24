package v1alpha1

import (
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KeycloakClientSpec defines the desired state of KeycloakClient.
// +k8s:openapi-gen=true
type KeycloakClientSpec struct {
	// Selector for looking up KeycloakRealm Custom Resources.
	// +kubebuilder:validation:Required
	RealmSelector *metav1.LabelSelector `json:"realmSelector"`
	// Keycloak Client REST object.
	// +kubebuilder:validation:Required
	Client *KeycloakAPIClient `json:"client"`
	// Client Roles
	// +optional
	// +listType=map
	// +listMapKey=name
	Roles []RoleRepresentation `json:"roles,omitempty"`
	// Scope Mappings
	// +optional
	ScopeMappings *MappingsRepresentation `json:"scopeMappings,omitempty"`
}

// https://www.keycloak.org/docs-api/11.0/rest-api/index.html#_mappingsrepresentation
type MappingsRepresentation struct {
	// Client Mappings
	// +optional
	ClientMappings map[string]ClientMappingsRepresentation `json:"clientMappings,omitempty"`

	// Realm Mappings
	// +optional
	RealmMappings []RoleRepresentation `json:"realmMappings,omitempty"`
}

// https://www.keycloak.org/docs-api/11.0/rest-api/index.html#_clientmappingsrepresentation
type ClientMappingsRepresentation struct {
	// Client
	// +optional
	Client string `json:"client,omitempty"`

	// ID
	// +optional
	ID string `json:"id,omitempty"`

	// Mappings
	// +optional
	Mappings []RoleRepresentation `json:"mappings,omitempty"`
}

type KeycloakAPIClient struct {
	// Client ID. If not specified, automatically generated.
	// +optional
	ID string `json:"id,omitempty"`
	// Client ID.
	// +kubebuilder:validation:Required
	ClientID string `json:"clientId"`
	// Client name.
	// +optional
	Name string `json:"name,omitempty"`
	// Surrogate Authentication Required option.
	// +optional
	SurrogateAuthRequired bool `json:"surrogateAuthRequired,omitempty"`
	// Client enabled flag.
	// +optional
	Enabled bool `json:"enabled,omitempty"`
	// What Client authentication type to use.
	// +optional
	ClientAuthenticatorType string `json:"clientAuthenticatorType,omitempty"`
	// Client Secret. The Operator will automatically create a Secret based on this value.
	// +optional
	Secret string `json:"secret,omitempty"`
	// Application base URL.
	// +optional
	BaseURL string `json:"baseUrl,omitempty"`
	// Application Admin URL.
	// +optional
	AdminURL string `json:"adminUrl,omitempty"`
	// Application root URL.
	// +optional
	RootURL string `json:"rootUrl,omitempty"`
	// Client description.
	// +optional
	Description string `json:"description,omitempty"`
	// Default Client roles.
	// +optional
	DefaultRoles []string `json:"defaultRoles,omitempty"`
	// A list of valid Redirection URLs.
	// +optional
	RedirectUris []string `json:"redirectUris,omitempty"`
	// A list of valid Web Origins.
	// +optional
	WebOrigins []string `json:"webOrigins,omitempty"`
	// Not Before setting.
	// +optional
	NotBefore int `json:"notBefore,omitempty"`
	// True if a client supports only Bearer Tokens.
	// +optional
	BearerOnly bool `json:"bearerOnly,omitempty"`
	// True if Consent Screen is required.
	// +optional
	ConsentRequired bool `json:"consentRequired,omitempty"`
	// True if Standard flow is enabled.
	// +optional
	StandardFlowEnabled bool `json:"standardFlowEnabled"`
	// True if Implicit flow is enabled.
	// +optional
	ImplicitFlowEnabled bool `json:"implicitFlowEnabled"`
	// True if Direct Grant is enabled.
	// +optional
	DirectAccessGrantsEnabled bool `json:"directAccessGrantsEnabled"`
	// True if Service Accounts are enabled.
	// +optional
	ServiceAccountsEnabled bool `json:"serviceAccountsEnabled,omitempty"`
	// True if this is a public Client.
	// +optional
	PublicClient bool `json:"publicClient"`
	// True if this client supports Front Channel logout.
	// +optional
	FrontchannelLogout bool `json:"frontchannelLogout,omitempty"`
	// Protocol used for this Client.
	// +optional
	Protocol string `json:"protocol,omitempty"`
	// Client Attributes.
	// +optional
	Attributes map[string]string `json:"attributes,omitempty"`
	// True if Full Scope is allowed.
	// +optional
	FullScopeAllowed *bool `json:"fullScopeAllowed,omitempty"`
	// Node registration timeout.
	// +optional
	NodeReRegistrationTimeout int `json:"nodeReRegistrationTimeout,omitempty"`
	// Protocol Mappers.
	// +optional
	ProtocolMappers []KeycloakProtocolMapper `json:"protocolMappers,omitempty"`
	// True to use a Template Config.
	// +optional
	UseTemplateConfig bool `json:"useTemplateConfig,omitempty"`
	// True to use Template Scope.
	// +optional
	UseTemplateScope bool `json:"useTemplateScope,omitempty"`
	// True to use Template Mappers.
	// +optional
	UseTemplateMappers bool `json:"useTemplateMappers,omitempty"`
	// Access options.
	// +optional
	Access map[string]bool `json:"access,omitempty"`
	// A list of optional client scopes. Optional client scopes are
	// applied when issuing tokens for this client, but only when they
	// are requested by the scope parameter in the OpenID Connect
	// authorization request.
	// +optional
	OptionalClientScopes []string `json:"optionalClientScopes,omitempty"`
	// A list of default client scopes. Default client scopes are
	// always applied when issuing OpenID Connect tokens or SAML
	// assertions for this client.
	// +optional
	DefaultClientScopes []string `json:"defaultClientScopes,omitempty"`
	// True if fine-grained authorization support is enabled for this client.
	// +optional
	AuthorizationServicesEnabled bool `json:"authorizationServicesEnabled,omitempty"`
	// Authorization settings for this resource server.
	// +optional
	AuthorizationSettings *KeycloakResourceServer `json:"authorizationSettings,omitempty"`
}

type KeycloakProtocolMapper struct {
	// Protocol Mapper ID.
	// +optional
	ID string `json:"id,omitempty"`
	// Protocol Mapper Name.
	// +optional
	Name string `json:"name,omitempty"`
	// Protocol to use.
	// +optional
	Protocol string `json:"protocol,omitempty"`
	// Protocol Mapper to use
	// +optional
	ProtocolMapper string `json:"protocolMapper,omitempty"`
	// True if Consent Screen is required.
	// +optional
	ConsentRequired bool `json:"consentRequired,omitempty"`
	// Text to use for displaying Consent Screen.
	// +optional
	ConsentText string `json:"consentText,omitempty"`
	// Config options.
	// +optional
	Config map[string]string `json:"config,omitempty"`
}

// https://www.keycloak.org/docs-api/12.0/rest-api/index.html#_resourceserverrepresentation
type KeycloakResourceServer struct {
	// True if resources should be managed remotely by the resource server.
	// +optional
	AllowRemoteResourceManagement bool `json:"allowRemoteResourceManagement,omitempty"`
	// Client ID.
	// +optional
	ClientID string `json:"clientId,omitempty"`
	// The decision strategy dictates how permissions are evaluated and how a
	// final decision is obtained. 'Affirmative' means that at least one
	// permission must evaluate to a positive decision in order to grant access
	// to a resource and its scopes. 'Unanimous' means that all permissions must
	// evaluate to a positive decision in order for the final decision to be also positive.
	// +optional
	DecisionStrategy string `json:"decisionStrategy,omitempty"`
	// ID.
	// +optional
	ID string `json:"id,omitempty"`
	// Name.
	// +optional
	Name string `json:"name,omitempty"`
	// Policies.
	// +optional
	Policies []KeycloakPolicy `json:"policies,omitempty"`
	// The policy enforcement mode dictates how policies are enforced when evaluating authorization requests.
	// 'Enforcing' means requests are denied by default even when there is no policy associated with a given resource.
	// 'Permissive' means requests are allowed even when there is no policy associated with a given resource.
	// 'Disabled' completely disables the evaluation of policies and allows access to any resource.
	// +optional
	PolicyEnforcementMode string `json:"policyEnforcementMode,omitempty"`
	// Resources.
	// +optional
	Resources []KeycloakResource `json:"resources,omitempty"`
	// Authorization Scopes.
	// +optional
	Scopes []KeycloakScope `json:"scopes,omitempty"`
}

// https://www.keycloak.org/docs-api/12.0/rest-api/index.html#_policyrepresentation
type KeycloakPolicy struct {
	// Config.
	// +optional
	Config map[string]string `json:"config,omitempty"`
	// The decision strategy dictates how the policies associated with a given permission are evaluated and how
	// a final decision is obtained. 'Affirmative' means that at least one policy must evaluate to a positive
	// decision in order for the final decision to be also positive. 'Unanimous' means that all policies must
	// evaluate to a positive decision in order for the final decision to be also positive. 'Consensus' means
	// that the number of positive decisions must be greater than the number of negative decisions. If the number
	// of positive and negative is the same, the final decision will be negative.
	// +optional
	DecisionStrategy string `json:"decisionStrategy,omitempty"`
	// A description for this policy.
	// +optional
	Description string `json:"description,omitempty"`
	// ID.
	// +optional
	ID string `json:"id,omitempty"`
	// The logic dictates how the policy decision should be made. If 'Positive', the resulting effect
	// (permit or deny) obtained during the evaluation of this policy will be used to perform a decision.
	// If 'Negative', the resulting effect will be negated, in other words, a permit becomes a deny and vice-versa.
	// +optional
	Logic string `json:"logic,omitempty"`
	// The name of this policy.
	// +optional
	Name string `json:"name,omitempty"`
	// Owner.
	// +optional
	Owner string `json:"owner,omitempty"`
	// Policies.
	// +optional
	Policies []string `json:"policies,omitempty"`
	// Resources.
	// +optional
	Resources []string `json:"resources,omitempty"`
	// Resources Data.
	// +optional
	ResourcesData []KeycloakResource `json:"resourcesData,omitempty"`
	// Scopes.
	// +optional
	Scopes []string `json:"scopes,omitempty"`
	// Type.
	// +optional
	Type string `json:"type,omitempty"`
	// Scopes Data.
	// +optional
	ScopesData []apiextensionsv1.JSON `json:"scopesData,omitempty"`
	// TODO: JSON struct is a workaround for the lack of support for recursive types
	// in CRD validation schemas. Keycloak will do validation for this field. Read more:
	// https://github.com/kubernetes/kubernetes/issues/62872
}

// https://www.keycloak.org/docs-api/12.0/rest-api/index.html#_resourcerepresentation
type KeycloakResource struct {
	// ID.
	// +optional
	ID string `json:"_id,omitempty"`
	// The attributes associated with the resource.
	// +optional
	Attributes map[string]string `json:"attributes,omitempty"`
	// A unique name for this resource. The name can be used to uniquely identify a resource, useful when
	// querying for a specific resource.
	// +optional
	DisplayName string `json:"displayName,omitempty"`
	// An URI pointing to an icon.
	// +optional
	IconURI string `json:"icon_uri,omitempty"`
	// A unique name for this resource. The name can be used to uniquely identify a resource, useful when
	// querying for a specific resource.
	// +optional
	Name string `json:"name,omitempty"`
	// True if the access to this resource can be managed by the resource owner.
	// +optional
	OwnerManagedAccess bool `json:"ownerManagedAccess,omitempty"`
	// The type of this resource. It can be used to group different resource instances with the same type.
	// +optional
	Type string `json:"type,omitempty"`
	// Set of URIs which are protected by resource.
	// +optional
	Uris []string `json:"uris,omitempty"`
	// The scopes associated with this resource.
	// +optional
	Scopes []apiextensionsv1.JSON `json:"scopes,omitempty"`
	// TODO: JSON struct is a workaround for the lack of support for recursive types
	// in CRD validation schemas. Keycloak will do validation for this field. Read more:
	// https://github.com/kubernetes/kubernetes/issues/62872
}

// https://www.keycloak.org/docs-api/12.0/rest-api/index.html#_scoperepresentation
type KeycloakScope struct {
	// A unique name for this scope. The name can be used to uniquely identify a scope, useful when querying
	// for a specific scope.
	// +optional
	DisplayName string `json:"displayName,omitempty"`
	// An URI pointing to an icon.
	// +optional
	IconURI string `json:"iconUri,omitempty"`
	// ID.
	// +optional
	ID string `json:"id,omitempty"`
	// A unique name for this scope. The name can be used to uniquely identify a scope, useful when querying
	// for a specific scope.
	// +optional
	Name string `json:"name,omitempty"`
	// Policies.
	// +optional
	Policies []KeycloakPolicy `json:"policies,omitempty"`
	// Resources.
	// +optional
	Resources []KeycloakResource `json:"resources,omitempty"`
}

// KeycloakClientStatus defines the observed state of KeycloakClient
// +k8s:openapi-gen=true
type KeycloakClientStatus struct {
	// Current phase of the operator.
	Phase StatusPhase `json:"phase"`
	// Human-readable message indicating details about current operator phase or error.
	Message string `json:"message"`
	// True if all resources are in a ready state and all work is done.
	Ready bool `json:"ready"`
	// A map of all the secondary resources types and names created for this CR. e.g "Deployment": [ "DeploymentName1", "DeploymentName2" ]
	SecondaryResources map[string][]string `json:"secondaryResources,omitempty"`
}

// KeycloakClient is the Schema for the keycloakclients API.
// +genclient
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type KeycloakClient struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KeycloakClientSpec   `json:"spec,omitempty"`
	Status KeycloakClientStatus `json:"status,omitempty"`
}

// KeycloakClientList contains a list of KeycloakClient.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type KeycloakClientList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KeycloakClient `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KeycloakClient{}, &KeycloakClientList{})
}

func (i *KeycloakClient) UpdateStatusSecondaryResources(kind string, resourceName string) {
	i.Status.SecondaryResources = UpdateStatusSecondaryResources(i.Status.SecondaryResources, kind, resourceName)
}
