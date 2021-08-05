package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	GroupFinalizer = "group.cleanup"
)

var (
	GroupPhaseReconciled StatusPhase = "reconciled"
	GroupPhaseFailing    StatusPhase = "failing"
)

// KeycloakGroupSpec defines the desired state of KeycloakGroup.
// +k8s:openapi-gen=true
type KeycloakGroupSpec struct {
	// Selector for looking up KeycloakRealm Custom Resources.
	// +kubebuilder:validation:Required
	RealmSelector *metav1.LabelSelector `json:"realmSelector,omitempty"`
	// Keycloak Group REST object.
	// +kubebuilder:validation:Required
	Group KeycloakAPIGroup `json:"group"`
}

// KeycloakUserStatus defines the observed state of KeycloakUser.
// +k8s:openapi-gen=true
type KeycloakGroupStatus struct {
	// Current phase of the operator.
	Phase StatusPhase `json:"phase"`
	// Human-readable message indicating details about current operator phase or error.
	Message string `json:"message"`
}

// KeycloakGroup is the Schema for the keycloakusers API.
// +genclient
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type KeycloakGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KeycloakGroupSpec   `json:"spec,omitempty"`
	Status KeycloakGroupStatus `json:"status,omitempty"`
}

type KeycloakAPIGroup struct {
	// Group ID.
	// +optional
	ID string `json:"id,omitempty"`

	// Group Name.
	// +kubebuilder:validation:Required
	Name string `json:"name,omitempty"`

	// Group Path.
	// +optional
	Path string `json:"path,omitempty"`

	// Realm Roles
	// +optional
	RealmRoles []string `json:"realmRoles,omitempty"`

	// Client Roles
	// +optional
	ClientRoles map[string][]string `json:"clientRoles,omitempty"`

	//
	//// A set of Attributes.
	//// +optional
	//Attributes map[string][]string `json:"attributes,omitempty"`
}

// KeycloakGroupList contains a list of KeycloakGroup
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type KeycloakGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KeycloakGroup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KeycloakGroup{}, &KeycloakGroupList{})
}
