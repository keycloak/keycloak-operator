package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KeycloakSpec defines the desired state of Keycloak.
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
	// Specify PodDisruptionBudget configuration.
	// +optional
	PodDisruptionBudget PodDisruptionBudgetConfig `json:"podDisruptionBudget,omitempty"`
	// Resources (Requests and Limits) for KeycloakDeployment.
	// +optional
	KeycloakDeploymentSpec DeploymentSpec `json:"keycloakDeploymentSpec,omitempty"`
	// Resources (Requests and Limits) for PostgresDeployment.
	// +optional
	PostgresDeploymentSpec DeploymentSpec `json:"postgresDeploymentSpec,omitempty"`
	// Specify Migration configuration
	// +optional
	Migration MigrateConfig `json:"migration,omitempty"`
	// Service annotations
	ServiceAnnotations map[string]string `json:"serviceAnnotations,omitempty"`
	// Name of the StorageClass for Postgresql Persistent Volume Claim
	// +optional
	StorageClassName *string `json:"storageClassName,omitempty"`
}

type DeploymentSpec struct {
	// Resources (Requests and Limits) for the Pods.
	// +optional
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	// Probes (Liveness and Readiness) for the Pods
	// +optional
	LivenessProbe corev1.Probe `json:"livenessProbe,omitempty"`
	ReadinessProbe corev1.Probe `json:"readinessProbe,omitempty"`
	// Container Environment variables
	// +optional
	EnvVars map[string]string `json:"envVars,omitempty"`
	// Pod Annotations
	// +optional
	PodAnnotations map[string]string `json:"podAnnotations,omitempty"`
}

type TLSTerminationType string

var (
	DefaultTLSTermintation        TLSTerminationType
	ReencryptTLSTerminationType   TLSTerminationType = "reencrypt"
	PassthroughTLSTerminationType TLSTerminationType = "passthrough"
)

type KeycloakExternalAccess struct {
	// If set to true, the Operator will create an Ingress or a Route
	// pointing to Keycloak.
	Enabled bool `json:"enabled,omitempty"`
	// TLS Termination type for the external access. Setting this field to "reencrypt" will
	// terminate TLS on the Ingress/Route level. Setting this field to "passthrough" will
	// send encrypted traffic to the Pod. If unspecified, defaults to "reencrypt".
	// Note, that this setting has no effect on Ingress
	// as Ingress TLS settings are not reconciled by this operator. In other words,
	// Ingress TLS configuration is the same in both cases and it is up to the user
	// to configure TLS section of the Ingress.
	TLSTermination TLSTerminationType `json:"tlsTermination,omitempty"`
}

type KeycloakExternalDatabase struct {
	// If set to true, the Operator will use an external database.
	// pointing to Keycloak.
	Enabled bool `json:"enabled,omitempty"`
}

type PodDisruptionBudgetConfig struct {
	// If set to true, the operator will create a PodDistruptionBudget for the Keycloak deployment and set its `maxUnavailable` value to 1.
	Enabled bool `json:"enabled,omitempty"`
}

type MigrateConfig struct {
	// Set it to config backup policy for migration
	// +optional
	Backups BackupConfig `json:"backups,omitempty"`
}

type BackupConfig struct {
	// If set to true, the operator will do database backup before doing migration
	Enabled bool `json:"enabled,omitempty"`
}

// KeycloakStatus defines the observed state of Keycloak.
// +k8s:openapi-gen=true
type KeycloakStatus struct {
	// Current phase of the operator.
	Phase StatusPhase `json:"phase"`
	// Human-readable message indicating details about current operator phase or error.
	Message string `json:"message"`
	// True if all resources are in a ready state and all work is done.
	Ready bool `json:"ready"`
	// A map of all the secondary resources types and names created for this CR. e.g "Deployment": [ "DeploymentName1", "DeploymentName2" ].
	SecondaryResources map[string][]string `json:"secondaryResources,omitempty"`
	// Version of Keycloak or RHSSO running on the cluster.
	Version string `json:"version"`
	// Service IP and Port for in-cluster access to the keycloak instance.
	InternalURL string `json:"internalURL"`
	// The secret where the admin credentials are to be found.
	CredentialSecret string `json:"credentialSecret"`
}

type StatusPhase string

var (
	NoPhase           StatusPhase
	PhaseReconciling  StatusPhase = "reconciling"
	PhaseFailing      StatusPhase = "failing"
	PhaseInitialising StatusPhase = "initialising"
)

// Keycloak is the Schema for the keycloaks API.
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Keycloak struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KeycloakSpec   `json:"spec,omitempty"`
	Status KeycloakStatus `json:"status,omitempty"`
}

// KeycloakList contains a list of Keycloak.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
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
