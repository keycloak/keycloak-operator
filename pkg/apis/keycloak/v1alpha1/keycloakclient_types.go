package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KeycloakClientSpec defines the desired state of KeycloakClient
// +k8s:openapi-gen=true
type KeycloakClientSpec struct {
	// +kubebuilder:validation:Required
	RealmSelector *metav1.LabelSelector `json:"realmSelector"`
	// +kubebuilder:validation:Required
	Client *KeycloakAPIClient `json:"client"`
}

type KeycloakAPIClient struct {
	ID string `json:"id,omitempty"`
	// +kubebuilder:validation:Required
	ClientID                  string                   `json:"clientId"`
	Name                      string                   `json:"name,omitempty"`
	SurrogateAuthRequired     bool                     `json:"surrogateAuthRequired,omitempty"`
	Enabled                   bool                     `json:"enabled,omitempty"`
	ClientAuthenticatorType   string                   `json:"clientAuthenticatorType,omitempty"`
	Secret                    string                   `json:"secret,omitempty"`
	BaseURL                   string                   `json:"baseUrl,omitempty"`
	AdminURL                  string                   `json:"adminUrl,omitempty"`
	RootURL                   string                   `json:"rootUrl,omitempty"`
	Description               string                   `json:"description,omitempty"`
	DefaultRoles              []string                 `json:"defaultRoles,omitempty"`
	RedirectUris              []string                 `json:"redirectUris,omitempty"`
	WebOrigins                []string                 `json:"webOrigins,omitempty"`
	NotBefore                 int                      `json:"notBefore,omitempty"`
	BearerOnly                bool                     `json:"bearerOnly,omitempty"`
	ConsentRequired           bool                     `json:"consentRequired,omitempty"`
	StandardFlowEnabled       bool                     `json:"standardFlowEnabled,omitempty"`
	ImplicitFlowEnabled       bool                     `json:"implicitFlowEnabled,omitempty"`
	DirectAccessGrantsEnabled bool                     `json:"directAccessGrantsEnabled,omitempty"`
	ServiceAccountsEnabled    bool                     `json:"serviceAccountsEnabled,omitempty"`
	PublicClient              bool                     `json:"publicClient,omitempty"`
	FrontchannelLogout        bool                     `json:"frontchannelLogout,omitempty"`
	Protocol                  string                   `json:"protocol,omitempty"`
	Attributes                map[string]string        `json:"attributes,omitempty"`
	FullScopeAllowed          bool                     `json:"fullScopeAllowed,omitempty"`
	NodeReRegistrationTimeout int                      `json:"nodeReRegistrationTimeout,omitempty"`
	ProtocolMappers           []KeycloakProtocolMapper `json:"protocolMappers,omitempty"`
	UseTemplateConfig         bool                     `json:"useTemplateConfig,omitempty"`
	UseTemplateScope          bool                     `json:"useTemplateScope,omitempty"`
	UseTemplateMappers        bool                     `json:"useTemplateMappers,omitempty"`
	Access                    map[string]bool          `json:"access,omitempty"`
}

type KeycloakProtocolMapper struct {
	ID              string            `json:"id,omitempty"`
	Name            string            `json:"name,omitempty"`
	Protocol        string            `json:"protocol,omitempty"`
	ProtocolMapper  string            `json:"protocolMapper,omitempty"`
	ConsentRequired bool              `json:"consentRequired,omitempty"`
	ConsentText     string            `json:"consentText,omitempty"`
	Config          map[string]string `json:"config,omitempty"`
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

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KeycloakClient is the Schema for the keycloakclients API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=keycloakclients,scope=Namespaced
type KeycloakClient struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KeycloakClientSpec   `json:"spec,omitempty"`
	Status KeycloakClientStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KeycloakClientList contains a list of KeycloakClient
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
