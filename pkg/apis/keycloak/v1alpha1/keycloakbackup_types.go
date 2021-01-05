package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KeycloakBackupSpec defines the desired state of KeycloakBackup.
// +k8s:openapi-gen=true
type KeycloakBackupSpec struct {
	// Controls automatic restore behavior.
	// Currently not implemented.
	//
	// In the future this will be used to trigger automatic restore for a given KeycloakBackup.
	// Each backup will correspond to a single snapshot of the database (stored either in a
	// Persistent Volume or AWS). If a user wants to restore it, all he/she needs to do is to
	// change this flag to true.
	// Potentially, it will be possible to restore a single backup multiple times.
	// +optional
	Restore bool `json:"restore,omitempty"`
	// If provided, an automatic database backup will be created on AWS S3 instead of
	// a local Persistent Volume. If this property is not provided - a local
	// Persistent Volume backup will be chosen.
	// +optional
	AWS KeycloakAWSSpec `json:"aws,omitempty"`
	// Selector for looking up Keycloak Custom Resources.
	// +kubebuilder:validation:Required
	InstanceSelector *metav1.LabelSelector `json:"instanceSelector,omitempty"`
	// Name of the StorageClass for Postgresql Backup Persistent Volume Claim
	// +optional
	StorageClassName *string `json:"storageClassName,omitempty"`
}

// KeycloakAWSSpec defines the desired state of KeycloakBackupSpec.
// +k8s:openapi-gen=true
type KeycloakAWSSpec struct {
	// If provided, the database backup will be encrypted.
	// Provides a secret name used for encrypting database data.
	// The secret needs to be in the following form:
	//
	//     apiVersion: v1
	//     kind: Secret
	//     metadata:
	//       name: <Secret name>
	//     type: Opaque
	//     stringData:
	//       GPG_PUBLIC_KEY: <GPG Public Key>
	//       GPG_TRUST_MODEL: <GPG Trust Model>
	//       GPG_RECIPIENT: <GPG Recipient>
	//
	// For more information, please refer to the Operator documentation.
	// +optional
	EncryptionKeySecretName string `json:"encryptionKeySecretName,omitempty"`
	// Provides a secret name used for connecting to AWS S3 Service.
	// The secret needs to be in the following form:
	//
	//     apiVersion: v1
	//     kind: Secret
	//     metadata:
	//       name: <Secret name>
	//     type: Opaque
	//     stringData:
	//       AWS_S3_BUCKET_NAME: <S3 Bucket Name>
	//       AWS_ACCESS_KEY_ID: <AWS Access Key ID>
	//       AWS_SECRET_ACCESS_KEY: <AWS Secret Key>
	//
	// For more information, please refer to the Operator documentation.
	// +kubebuilder:validation:Required
	CredentialsSecretName string `json:"credentialsSecretName,omitempty"`
	// If specified, it will be used as a schedule for creating a CronJob.
	// +optional
	Schedule string `json:"schedule,omitempty"`
}

type BackupStatusPhase string

var (
	BackupPhaseNone        BackupStatusPhase
	BackupPhaseReconciling BackupStatusPhase = "reconciling"
	BackupPhaseCreated     BackupStatusPhase = "created"
	BackupPhaseRestored    BackupStatusPhase = "restored"
	BackupPhaseFailing     BackupStatusPhase = "failing"
)

// KeycloakBackupStatus defines the observed state of KeycloakBackup.
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

// KeycloakBackup is the Schema for the keycloakbackups API.
// +genclient
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type KeycloakBackup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KeycloakBackupSpec   `json:"spec,omitempty"`
	Status KeycloakBackupStatus `json:"status,omitempty"`
}

// KeycloakBackupList contains a list of KeycloakBackup.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
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
