package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ExternalKeycloakSpec defines the desired state of ExternalKeycloak.
// +k8s:openapi-gen=true
type ExternalKeycloakSpec struct {
	// References the secret that stores the credentials for keycloak.
	// +kubebuilder:validation:Required
	CredentialSecret string `json:"credentialSecret,omitempty"`
	// The endpoint to contact the keycloak instance.
	// +kubebuilder:validation:Required
	Endpoint string `json:"endpoint,omitempty"`
}

type ExternalKeycloakStatusPhase string

var (
	ExternalKeycloakPhaseNone ExternalKeycloakStatusPhase
)

// ExternalKeycloakStatus defines the observed state of KeycloakBackup.
// +k8s:openapi-gen=true
type ExternalKeycloakStatus struct {
	// Current phase of the operator.
	Phase ExternalKeycloakStatusPhase `json:"phase"`
	// Human-readable message indicating details about current operator phase or error.
	Message string `json:"message"`
	// True if all resources are in a ready state and all work is done.
	Ready bool `json:"ready"`
}

// ExternalKeycloak is the Schema for the externalkeycloaks API.
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ExternalKeycloak struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ExternalKeycloakSpec   `json:"spec,omitempty"`
	Status ExternalKeycloakStatus `json:"status,omitempty"`
}

// ExternalKeycloakList contains a list of ExternalKeycloak.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ExternalKeycloakList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ExternalKeycloak `json:"items"`
}

func (in *ExternalKeycloak) Endpoint() string {
	return in.Spec.Endpoint
}

func (in *ExternalKeycloak) CredentialSecret() string {
	return in.Spec.CredentialSecret
}

func init() {
	SchemeBuilder.Register(&ExternalKeycloak{}, &ExternalKeycloakList{})
}
