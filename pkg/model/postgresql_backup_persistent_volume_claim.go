package model

import (
	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func PostgresqlBackupPersistentVolumeClaim(cr *v1alpha1.KeycloakBackup) *v1.PersistentVolumeClaim {
	return &v1.PersistentVolumeClaim{
		ObjectMeta: v12.ObjectMeta{
			Name:      PostgresqlBackupPersistentVolumeName + "-" + cr.Name,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app":       ApplicationName,
				"component": PostgresqlBackupComponent,
			},
		},
		Spec: v1.PersistentVolumeClaimSpec{
			AccessModes: []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{
					v1.ResourceStorage: resource.MustParse(PostgresqlPersistentVolumeCapacity),
				}},
			StorageClassName: cr.Spec.StorageClassName,
		},
	}
}

func PostgresqlBackupPersistentVolumeClaimSelector(cr *v1alpha1.KeycloakBackup) client.ObjectKey {
	return client.ObjectKey{
		Name:      PostgresqlBackupPersistentVolumeName + "-" + cr.Name,
		Namespace: cr.Namespace,
	}
}

func PostgresqlBackupPersistentVolumeClaimReconciled(cr *v1alpha1.KeycloakBackup, currentState *v1.PersistentVolumeClaim) *v1.PersistentVolumeClaim {
	reconciled := currentState.DeepCopy()
	reconciled.Spec.AccessModes = []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce}
	reconciled.Spec.Resources = v1.ResourceRequirements{
		Requests: v1.ResourceList{
			v1.ResourceStorage: resource.MustParse(PostgresqlPersistentVolumeCapacity),
		}}
	if cr.Spec.StorageClassName != nil {
		reconciled.Spec.StorageClassName = cr.Spec.StorageClassName
	}
	return reconciled
}
