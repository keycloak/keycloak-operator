package common

import (
	"fmt"
	"testing"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestUtil_Test_validateList_With_Previous_Error(t *testing.T) {
	// given
	err := fmt.Errorf("testError")
	var items []v1alpha1.KeycloakRealm

	// when
	returnedError := validateList(err, items)

	// then
	assert.Equal(t, err, returnedError)
}

func TestUtil_Test_validateList_With_Empty_List(t *testing.T) {
	// given
	var items []v1alpha1.KeycloakRealm

	// when
	returnedError := validateList(nil, items)

	// then
	assert.NotNil(t, returnedError)
}

func TestUtil_Test_validateList_With_Proper_Arguments(t *testing.T) {
	// given
	items := []string{"test"}

	// when
	returnedError := validateList(nil, items)

	// then
	assert.Nil(t, returnedError)
}
