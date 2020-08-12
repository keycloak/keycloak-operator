package keycloak

import (
	"testing"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Checks if all CRs implement runtime.Object
// see: https://www.openshift.com/blog/kubernetes-deep-dive-code-generation-customresources
func TestApi_check_Runtime_Object_Assignment(t *testing.T) {
	var _ runtime.Object = &v1alpha1.Keycloak{}
	var _ runtime.Object = &v1alpha1.KeycloakRealm{}
	var _ runtime.Object = &v1alpha1.KeycloakClient{}
	var _ runtime.Object = &v1alpha1.KeycloakBackup{}
	var _ runtime.Object = &v1alpha1.KeycloakUser{}
}
