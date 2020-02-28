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
			DatabaseSecretUsernameProperty: PostgresqlUsername,
			DatabaseSecretPasswordProperty: cr.ObjectMeta.Name + "-" + GenerateRandomString(PostgresqlPasswordLength),
			// The 3 entries below are not used by the Operator itself but rather by the Backup container
			DatabaseSecretDatabaseProperty:  PostgresqlDatabase,
			DatabaseSecretSuperuserProperty: "true",
			DatabaseSecretHostProperty:      PostgresqlServiceName,
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
		reconciled.StringData[DatabaseSecretUsernameProperty] = PostgresqlUsername
	}
	if _, ok := reconciled.Data[DatabaseSecretPasswordProperty]; !ok {
		reconciled.StringData[DatabaseSecretPasswordProperty] = cr.ObjectMeta.Name + "-" + GenerateRandomString(PostgresqlPasswordLength)
	}
	if _, ok := reconciled.Data[DatabaseSecretDatabaseProperty]; !ok {
		reconciled.StringData[DatabaseSecretDatabaseProperty] = PostgresqlDatabase
	}
	if _, ok := reconciled.Data[DatabaseSecretSuperuserProperty]; !ok {
		reconciled.StringData[DatabaseSecretSuperuserProperty] = "true"
	}
	if _, ok := reconciled.Data[DatabaseSecretHostProperty]; !ok {
		reconciled.StringData[DatabaseSecretHostProperty] = PostgresqlServiceName
	}
	return reconciled
}
