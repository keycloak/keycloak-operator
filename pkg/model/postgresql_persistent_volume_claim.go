package model

import (
	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func PostgresqlPersistentVolumeClaim(cr *v1alpha1.Keycloak) *v1.PersistentVolumeClaim {
	return &v1.PersistentVolumeClaim{
		ObjectMeta: v12.ObjectMeta{
			Name:      PostgresqlPersistentVolumeName,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app": ApplicationName,
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

func PostgresqlPersistentVolumeClaimSelector(cr *v1alpha1.Keycloak) client.ObjectKey {
	return client.ObjectKey{
		Name:      PostgresqlPersistentVolumeName,
		Namespace: cr.Namespace,
	}
}

func PostgresqlPersistentVolumeClaimReconciled(cr *v1alpha1.Keycloak, currentState *v1.PersistentVolumeClaim) *v1.PersistentVolumeClaim {
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
