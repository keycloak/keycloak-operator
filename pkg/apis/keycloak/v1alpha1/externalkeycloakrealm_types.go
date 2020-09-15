package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ExternalKeycloakRealmSpec defines the desired state of ExternalKeycloakRealm.
// +k8s:openapi-gen=true
type ExternalKeycloakRealmSpec struct {
	// Selector for looking up Keycloak or ExternalKeycloak Custom Resources.
	// +kubebuilder:validation:Required
	InstanceSelector *metav1.LabelSelector `json:"instanceSelector,omitempty"`
	// Realm name.
	// +kubebuilder:validation:Required
	Realm string `json:"realm"`
}

type ExternalKeycloakRealmStatusPhase string

var (
	ExternalKeycloakRealmPhaseNone ExternalKeycloakRealmStatusPhase
)

// ExternalKeycloakStatus defines the observed state of KeycloakBackup.
// +k8s:openapi-gen=true
type ExternalKeycloakRealmStatus struct {
	// Current phase of the operator.
	Phase ExternalKeycloakRealmStatusPhase `json:"phase"`
	// Human-readable message indicating details about current operator phase or error.
	Message string `json:"message"`
	// True if all resources are in a ready state and all work is done.
	Ready bool `json:"ready"`
}

// ExternalKeycloakRealm is the Schema for the externalkeycloakrealms API.
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ExternalKeycloakRealm struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ExternalKeycloakRealmSpec   `json:"spec,omitempty"`
	Status ExternalKeycloakRealmStatus `json:"status,omitempty"`
}

// ExternalKeycloakRealmList contains a list of ExternalKeycloakRealm.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ExternalKeycloakRealmList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ExternalKeycloakRealm `json:"items"`
}

func (i *ExternalKeycloakRealm) Realm() string {
	return i.Spec.Realm
}

func (i *ExternalKeycloakRealm) InstanceSelector() *metav1.LabelSelector {
	return i.Spec.InstanceSelector.DeepCopy()
}

func init() {
	SchemeBuilder.Register(&ExternalKeycloakRealm{}, &ExternalKeycloakRealmList{})
}
