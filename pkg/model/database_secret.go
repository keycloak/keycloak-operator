package model

import (
	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func DatabaseSecret(cr *v1alpha1.Keycloak) *v1.Secret {
	return &v1.Secret{
		ObjectMeta: v12.ObjectMeta{
			Name:      DatabaseSecretName,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app": ApplicationName,
			},
		},
		StringData: map[string]string{
			DatabaseSecretUsernameProperty: cr.ObjectMeta.Name + "-" + randStringRunes(4),
			DatabaseSecretPasswordProperty: cr.ObjectMeta.Name + "-" + randStringRunes(4),
		},
	}
}

func DatabaseSecretSelector(cr *v1alpha1.Keycloak) client.ObjectKey {
	return client.ObjectKey{
		Name:      DatabaseSecretName,
		Namespace: cr.Namespace,
	}
}

func DatabaseSecretReconciled(cr *v1alpha1.Keycloak, currentState *v1.Secret) *v1.Secret {
	reconciled := currentState.DeepCopy()
	// K8s automatically converts StringData to Data when getting the resource
	if _, ok := reconciled.Data[DatabaseSecretUsernameProperty]; !ok {
		reconciled.StringData[DatabaseSecretUsernameProperty] = cr.ObjectMeta.Name + "-" + randStringRunes(4)
	}
	if _, ok := reconciled.Data[DatabaseSecretPasswordProperty]; !ok {
		reconciled.StringData[DatabaseSecretPasswordProperty] = cr.ObjectMeta.Name + "-" + randStringRunes(4)
	}
	return reconciled
}
