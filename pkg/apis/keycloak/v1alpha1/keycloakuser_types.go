package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	UserFinalizer = "user.cleanup"
)

var (
	UserPhaseReconciled StatusPhase = "reconciled"
	UserPhaseFailing    StatusPhase = "failing"
)

// KeycloakUserSpec defines the desired state of KeycloakUser
// +k8s:openapi-gen=true
type KeycloakUserSpec struct {
	RealmSelector *metav1.LabelSelector `json:"realmSelector,omitempty"`
	User          KeycloakAPIUser       `json:"user"`
}

// KeycloakUserStatus defines the observed state of KeycloakUser
// +k8s:openapi-gen=true
type KeycloakUserStatus struct {
	Phase   StatusPhase `json:"phase"`
	Message string      `json:"message"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KeycloakUser is the Schema for the keycloakusers API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=keycloakusers,scope=Namespaced
type KeycloakUser struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KeycloakUserSpec   `json:"spec,omitempty"`
	Status KeycloakUserStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KeycloakUserList contains a list of KeycloakUser
type KeycloakUserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KeycloakUser `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KeycloakUser{}, &KeycloakUserList{})
}
