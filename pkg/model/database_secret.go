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
		Data: map[string][]byte{
			DatabaseSecretUsernameProperty: []byte(PostgresqlUsername),
			DatabaseSecretPasswordProperty: []byte(cr.ObjectMeta.Name + "-" + GenerateRandomString(PostgresqlPasswordLength)),
			// The 3 entries below are not used by the Operator itself but rather by the Backup container
			DatabaseSecretDatabaseProperty: []byte(PostgresqlDatabase),
			DatabaseSecretHostProperty:     []byte(PostgresqlServiceName),
			DatabaseSecretVersionProperty:  []byte("10"),
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
	if reconciled.Data == nil || len(reconciled.Data) == 0 {
		reconciled.Data = make(map[string][]byte)
	}

	if _, ok := reconciled.Data[DatabaseSecretUsernameProperty]; !ok {
		reconciled.Data[DatabaseSecretUsernameProperty] = []byte(PostgresqlUsername)
	}
	if _, ok := reconciled.Data[DatabaseSecretPasswordProperty]; !ok {
		reconciled.Data[DatabaseSecretPasswordProperty] = []byte(cr.ObjectMeta.Name + "-" + GenerateRandomString(PostgresqlPasswordLength))
	}
	if _, ok := reconciled.Data[DatabaseSecretDatabaseProperty]; !ok {
		reconciled.Data[DatabaseSecretDatabaseProperty] = []byte(PostgresqlDatabase)
	}
	if _, ok := reconciled.Data[DatabaseSecretHostProperty]; !ok {
		reconciled.Data[DatabaseSecretHostProperty] = []byte(PostgresqlServiceName)
	}
	if _, ok := reconciled.Data[DatabaseSecretVersionProperty]; !ok {
		reconciled.Data[DatabaseSecretVersionProperty] = []byte("10")
	}
	return reconciled
}
