package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KeycloakBackupSpec defines the desired state of KeycloakBackup
// +k8s:openapi-gen=true
type KeycloakBackupSpec struct {
	Restore bool            `json:"restore,omitempty"`
	AWS     KeycloakAWSSpec `json:"aws,omitempty"`
}

// KeycloakAWSSpec defines the desired state of KeycloakBackupSpec
// +k8s:openapi-gen=true
type KeycloakAWSSpec struct {
	EncryptionKeySecretName string `json:"encryptionKeySecretName,omitempty"`
	CredentialsSecretName   string `json:"credentialsSecretName,omitempty"`
	Schedule                string `json:"schedule,omitempty"`
}

type BackupStatusPhase string

var (
	BackupPhaseNone        BackupStatusPhase
	BackupPhaseReconciling BackupStatusPhase = "reconciling"
	BackupPhaseCreated     BackupStatusPhase = "created"
	BackupPhaseRestored    BackupStatusPhase = "restored"
	BackupPhaseFailing     BackupStatusPhase = "failing"
)

// KeycloakBackupStatus defines the observed state of KeycloakBackup
// +k8s:openapi-gen=true
type KeycloakBackupStatus struct {
	// Current phase of the operator.
	Phase BackupStatusPhase `json:"phase"`
	// Human-readable message indicating details about current operator phase or error.
	Message string `json:"message"`
	// True if all resources are in a ready state and all work is done.
	Ready bool `json:"ready"`
	// A map of all the secondary resources types and names created for this CR. e.g "Deployment": [ "DeploymentName1", "DeploymentName2" ]
	SecondaryResources map[string][]string `json:"secondaryResources,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KeycloakBackup is the Schema for the keycloakbackups API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=keycloakbackups,scope=Namespaced
type KeycloakBackup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KeycloakBackupSpec   `json:"spec,omitempty"`
	Status KeycloakBackupStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KeycloakBackupList contains a list of KeycloakBackup
type KeycloakBackupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KeycloakBackup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KeycloakBackup{}, &KeycloakBackupList{})
}

func (i *KeycloakBackup) UpdateStatusSecondaryResources(kind string, resourceName string) {
	i.Status.SecondaryResources = UpdateStatusSecondaryResources(i.Status.SecondaryResources, kind, resourceName)
}
