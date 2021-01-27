package common

import (
	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/keycloak/keycloak-operator/pkg/model"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
	config2 "sigs.k8s.io/controller-runtime/pkg/client/config"
)

type LoginType int32

const (
	Unknown                       LoginType = -1
	UsingClientCredentials        LoginType = 0
	UsingAdminUsernameAndPassword LoginType = 1
)

type KeycloakConnectionFactory interface {
	CreateConnection() (KeycloakInterface, error)
}

type DefaultKeycloakConnectionFactory struct {
	keycloakCR       *v1alpha1.Keycloak
	realmCR          *v1alpha1.KeycloakRealm
	clientCR         *v1alpha1.KeycloakClient
	ctx              context.Context
	controllerClient client.Client
}

func NewDefaultKeycloakConnectionFactory(ctx context.Context, c client.Client, keycloakCR *v1alpha1.Keycloak, realmCR *v1alpha1.KeycloakRealm) KeycloakConnectionFactory {
	operatorCLIClient, err := getMatchingCLIClient(ctx, c, keycloakCR, realmCR)
	if err != nil {
		log.V(-1).Info("error while obtaining Operator CLI Client. Falling back to Username and Password", "err", err)
	}

	return &DefaultKeycloakConnectionFactory{
		keycloakCR:       keycloakCR,
		realmCR:          realmCR,
		clientCR:         operatorCLIClient,
		ctx:              ctx,
		controllerClient: c,
	}
}

// Try to get a matching CLI Client for a Realm
func getMatchingCLIClient(ctx context.Context, c client.Client, keycloakCR *v1alpha1.Keycloak, realmCR *v1alpha1.KeycloakRealm) (*v1alpha1.KeycloakClient, error) {
	owningRealm := realmCR

	if owningRealm == nil {
		masterRealm := model.KeycloakMasterRealm(keycloakCR)
		masterRealmSelector := model.KeycloakMasterRealmSelector(keycloakCR)

		err := c.Get(ctx, masterRealmSelector, masterRealm)
		if err != nil {
			return nil, err
		}

		owningRealm = masterRealm
	}

	keycloakCLIClient := model.KeycloakOperatorCLIClient(owningRealm)
	keycloakCLIClientSelector := model.KeycloakOperatorCLIClientSelector(owningRealm)

	err := c.Get(ctx, keycloakCLIClientSelector, keycloakCLIClient)
	if err != nil {
		return nil, err
	}

	return keycloakCLIClient, nil
}

func (k *DefaultKeycloakConnectionFactory) loginDecision() (LoginType, error) {
	if model.Profiles.UseDefaultAuthenticationMode() {
		if k.keycloakCR == nil {
			return Unknown, errors.Errorf("cannot perform realm create when client is nil")
		}
		if k.clientCR != nil && k.clientCR.Status.Ready {
			return UsingClientCredentials, nil
		}
	}
	return UsingAdminUsernameAndPassword, nil
}

func (k *DefaultKeycloakConnectionFactory) getKubeClient() (*kubernetes.Clientset, error) {
	config, err := config2.GetConfig()
	if err != nil {
		return nil, err
	}

	secretClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return secretClient, nil
}

func (k *DefaultKeycloakConnectionFactory) CreateConnection() (KeycloakInterface, error) {
	secretClient, err := k.getKubeClient()
	if err != nil {
		return nil, err
	}
	loginDecision, err := k.loginDecision()
	if err != nil {
		return nil, err
	}
	switch loginDecision {
	case UsingClientCredentials:
		log.V(1).Info("Logging in using Client Credentials")
		clientSecret, err := secretClient.CoreV1().Secrets(k.clientCR.Namespace).Get(k.ctx, "keycloak-client-secret-"+k.clientCR.Name, v12.GetOptions{})
		if err != nil {
			return nil, errors.Wrap(err, "failed to get the client secret")
		}
		clientID := string(clientSecret.Data[model.ClientSecretClientIDProperty])
		secret := string(clientSecret.Data[model.ClientSecretClientSecretProperty])
		url := k.keycloakCR.Status.InternalURL
		client := &Client{
			URL:       url,
			requester: defaultRequester(),
		}
		if err := client.clientCredentialsLogin(clientID, secret, k.realmCR.Spec.Realm.Realm); err != nil {
			return nil, err
		}
		return client, nil
	case UsingAdminUsernameAndPassword:
		log.V(1).Info("Logging in using Admin username and password")
		adminCreds, err := secretClient.CoreV1().Secrets(k.keycloakCR.Namespace).Get(k.ctx, k.keycloakCR.Status.CredentialSecret, v12.GetOptions{})
		if err != nil {
			return nil, errors.Wrap(err, "failed to get the admin credentials")
		}
		user := string(adminCreds.Data[model.AdminUsernameProperty])
		pass := string(adminCreds.Data[model.AdminPasswordProperty])
		url := k.keycloakCR.Status.InternalURL
		client := &Client{
			URL:       url,
			requester: defaultRequester(),
		}
		if err := client.usernameAndPasswordLogin(user, pass); err != nil {
			return nil, err
		}
		return client, nil
	}
	return nil, errors.Errorf("cannot obtain authenticated client")
}
