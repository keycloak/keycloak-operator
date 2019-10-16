package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KeycloakSpec defines the desired state of Keycloak
// +k8s:openapi-gen=true
type KeycloakSpec struct {
	ExternalDatabaseSecret string                 `json:"externalDatabaseSecret,omitempty"`
	AdminCredentialSecret  string                 `json:"adminCredentialSecret,omitempty"`
	Extensions             []string               `json:"extensions,omitempty"`
	Instances              int                    `json:"instances,omitempty"`
	ExternalAccess         KeycloakExternalAccess `json:"externalAccess,omitempty"`
	Profile                string                 `json:"profile,omitempty"`
}

type KeycloakExternalAccess struct {
	Enabled bool `json:"enabled,omitempty"`
}

// KeycloakStatus defines the observed state of Keycloak
// +k8s:openapi-gen=true
type KeycloakStatus struct {
	// Current phase of the operator.
	Phase StatusPhase `json:"phase"`
	// Human-readable message indicating details about current operator phase or error.
	Message string `json:"message"`
	// True if all resources are in a ready state and all work is done.
	Ready bool `json:"ready"`
	// A map of all the secondary resources types and names created for this CR. e.g "Deployment": [ "DeploymentName1", "DeploymentName2" ]
	SecondaryResources map[string][]string `json:"secondaryResources,omitempty"`
	// Version of Keycloak or RHSSO running on the cluster
	Version string `json:"version"`
	// Service IP and Port for in-cluster access to the keycloak instance
	InternalUrl string `json:"url"`
	// The secret where the admin credentials are to be found
	CredentialSecret string `json:"credentialSecret"`
}

type StatusPhase string

var (
	NoPhase           StatusPhase
	PhaseReconciling  StatusPhase = "reconciling"
	PhaseFailing      StatusPhase = "failing"
	PhaseInitialising StatusPhase = "initialising"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Keycloak is the Schema for the keycloaks API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Keycloak struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KeycloakSpec   `json:"spec,omitempty"`
	Status KeycloakStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KeycloakList contains a list of Keycloak
type KeycloakList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Keycloak `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Keycloak{}, &KeycloakList{})
}

func (i *Keycloak) UpdateStatusSecondaryResources(kind string, resourceName string) {
	// If the map is nil, instansiate it
	if i.Status.SecondaryResources == nil {
		i.Status.SecondaryResources = make(map[string][]string)
	}

	// return if the resource name already exists in the slice
	for _, ele := range i.Status.SecondaryResources[kind] {
		if ele == resourceName {
			return
		}
	}
	// add the resource name to the list of secondary resources in the status
	i.Status.SecondaryResources[kind] = append(i.Status.SecondaryResources[kind], resourceName)
}
