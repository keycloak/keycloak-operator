package model

import (
	"strconv"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func PostgresqlServiceEndpoints(cr *v1alpha1.Keycloak) *v1.Endpoints {
	return &v1.Endpoints{
		ObjectMeta: v12.ObjectMeta{
			Name:      PostgresqlServiceName,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app": ApplicationName,
			},
		},
		Subsets: []v1.EndpointSubset{{
			Addresses: []v1.EndpointAddress{{}},
			Ports:     []v1.EndpointPort{{}},
		}},
	}
}

func PostgresqlServiceEndpointsSelector(cr *v1alpha1.Keycloak) client.ObjectKey {
	return client.ObjectKey{
		Name:      PostgresqlServiceName,
		Namespace: cr.Namespace,
	}
}

func PostgresqlServiceEndpointsReconciled(cr *v1alpha1.Keycloak, currentState *v1.Endpoints, currentDatabaseSecret *v1.Secret) *v1.Endpoints {
	reconciled := currentState.DeepCopy()

	// We don't need any error reporting here. Endpoints are more "internal" K8s objects (not very often
	// configured by users), so they're very picky. If the IP is missing, or is invalid, we'll get
	// proper notification in Keycloak CR (Status.Message). Maybe someday we'll have a Validating WebHook to
	// improve the user experience.

	port := string(currentDatabaseSecret.Data[DatabaseSecretExternalPortProperty])
	portAsInt, err := strconv.ParseInt(port, 10, 32)
	if err != nil {
		// Default Postgresql Port - maybe we'll be lucky...
		portAsInt = 5432
	}

	// Sometimes it happens that Kubernetes doesn't create the slices (it's bad timing I guess).
	// In that case, we need to do it....
	if len(reconciled.Subsets) == 0 {
		reconciled.Subsets = []v1.EndpointSubset{{
			Addresses: []v1.EndpointAddress{{}},
			Ports:     []v1.EndpointPort{{}},
		}}
	}

	reconciled.Subsets[0].Ports = []v1.EndpointPort{{
		Port:     int32(portAsInt),
		Protocol: "TCP",
	}}

	hostname := string(currentDatabaseSecret.Data[DatabaseSecretExternalAddressProperty])
	if hostname != "" {
		reconciled.Subsets[0].Addresses = []v1.EndpointAddress{{
			// According to the comments in K8s' EndpointAddress, this field
			// should work for arbitrary IPs as well as hostnames.
			IP: hostname,
		}}
	}

	return reconciled
}
