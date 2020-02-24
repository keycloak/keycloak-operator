package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KeycloakSpec defines the desired state of Keycloak
// +k8s:openapi-gen=true
type KeycloakSpec struct {
	// A list of extensions, where each one is a URL to a JAR files that will be deployed in Keycloak.
	// +listType=set
	// +optional
	Extensions []string `json:"extensions,omitempty"`
	// Number of Keycloak instances in HA mode. Default is 1.
	// +optional
	Instances int `json:"instances,omitempty"`
	// Controls external Ingress/Route settings.
	// +optional
	ExternalAccess KeycloakExternalAccess `json:"externalAccess,omitempty"`
	// Controls external database settings.
	// Using an external database requires providing a secret containing credentials
	// as well as connection details. Here's an example of such secret:
	//
	//     apiVersion: v1
	//     kind: Secret
	//     metadata:
	//         name: keycloak-db-secret
	//         namespace: keycloak
	//     stringData:
	//         POSTGRES_DATABASE: <Database Name>
	//         POSTGRES_EXTERNAL_ADDRESS: <External Database IP or URL (resolvable by K8s)>
	//         POSTGRES_EXTERNAL_PORT: <External Database Port>
	//         # Strongly recommended to use <'Keycloak CR Name'-postgresql>
	//         POSTGRES_HOST: <Database Service Name>
	//         POSTGRES_PASSWORD: <Database Password>
	//         # Required for AWS Backup functionality
	//         POSTGRES_SUPERUSER: true
	//         POSTGRES_USERNAME: <Database Username>
	//      type: Opaque
	//
	// Both POSTGRES_EXTERNAL_ADDRESS and POSTGRES_EXTERNAL_PORT are specifically required for creating
	// connection to the external database. The secret name is created using the following convention:
	//       <Custom Resource Name>-db-secret
	//
	// For more information, please refer to the Operator documentation.
	// +optional
	ExternalDatabase KeycloakExternalDatabase `json:"externalDatabase,omitempty"`
	// Profile used for controlling Operator behavior. Default is empty.
	// +optional
	Profile string `json:"profile,omitempty"`
	// Specify PodDisruptionBudget configuration
	// +optional
	PodDisruptionBudget PodDisruptionBudgetConfig `json:"podDisruptionBudget,omitempty"`
}

type KeycloakExternalAccess struct {
	// If set to true, the Operator will create an Ingress or a Route
	// pointing to Keycloak.
	Enabled bool `json:"enabled,omitempty"`
}

type KeycloakExternalDatabase struct {
	// If set to true, the Operator will use an external database.
	// pointing to Keycloak.
	Enabled bool `json:"enabled,omitempty"`
}

type PodDisruptionBudgetConfig struct {
	// If set to true, the operator will create a PodDistruptionBudget for the Keycloak deployment and set its `maxUnavailable` value to 1
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
	InternalURL string `json:"internalURL"`
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
// +kubebuilder:resource:path=keycloaks,scope=Namespaced
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
	i.Status.SecondaryResources = UpdateStatusSecondaryResources(i.Status.SecondaryResources, kind, resourceName)
}
