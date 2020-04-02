package model

import (
	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func KeycloakMigrationOneTimeBackup(cr *v1alpha1.KeycloakBackup) *v1alpha1.KeycloakBackup {
	labelSelect := metav1.LabelSelector{
		MatchLabels: cr.Labels,
	}
	return &v1alpha1.KeycloakBackup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    cr.Labels,
		},
		Spec: v1alpha1.KeycloakBackupSpec{
			InstanceSelector: &labelSelect,
		},
	}
}

func KeycloakMigrationOneTimeBackupSelector(cr *v1alpha1.KeycloakBackup) client.ObjectKey {
	return client.ObjectKey{
		Name:      cr.Name,
		Namespace: cr.Namespace,
	}
}
