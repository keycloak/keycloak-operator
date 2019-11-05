package shared

import (
	"context"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/common"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
	config2 "sigs.k8s.io/controller-runtime/pkg/client/config"
)

// Try to get a list of keycloak instances that match the selector specified on the realm
func GetMatchingKeycloaks(c client.Client, ctx context.Context, realm *v1alpha1.KeycloakRealm) (v1alpha1.KeycloakList, error) {
	var list v1alpha1.KeycloakList
	opts := &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(realm.Spec.InstanceSelector.MatchLabels),
	}

	err := c.List(ctx, opts, &list)
	if err != nil {
		return list, err
	}

	return list, nil
}

// Try to get a list of keycloak instances that match the selector specified on the realm
func GetMatchingRealms(c client.Client, ctx context.Context, user *v1alpha1.KeycloakUser) (v1alpha1.KeycloakRealmList, error) {
	var list v1alpha1.KeycloakRealmList
	opts := &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(user.Spec.RealmSelector.MatchLabels),
	}

	err := c.List(ctx, opts, &list)
	if err != nil {
		return list, err
	}

	return list, nil
}

func GetAuthenticatedClient(kc v1alpha1.Keycloak) (common.KeycloakInterface, error) {
	config, err := config2.GetConfig()
	if err != nil {
		return nil, err
	}

	secretClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	factory := common.KeycloakFactory{
		SecretClient: secretClient.CoreV1().Secrets(kc.Namespace),
	}

	return factory.AuthenticatedClient(kc)
}
