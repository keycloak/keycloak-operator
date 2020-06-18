package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtil_Test_GetServiceEnvVar(t *testing.T) {
	assert.Equal(t, GetServiceEnvVar("SERVICE_HOST"), "KEYCLOAK_POSTGRESQL_SERVICE_HOST")
	assert.Equal(t, GetServiceEnvVar("SERVICE_PORT"), "KEYCLOAK_POSTGRESQL_SERVICE_PORT")
}

func TestUtil_SanitizeResourceName(t *testing.T) {
	expected := map[string]string{
		// Allowed characters
		"test123-_.": "test123--.",
		// Mixed of allowed characters and disallowed characters
		"testTEST[(/%^&*,)]123-_.": "testtest123--.",
	}

	for input, output := range expected {
		actual := SanitizeResourceName(input)
		assert.Equal(t, output, actual)
	}
}

func TestIsIP(t *testing.T) {
	assert.True(t, IsIP([]byte("54.154.171.84")))
	assert.False(t, IsIP([]byte("this.is.a.hostname")))
	assert.False(t, IsIP([]byte("http://www.database.url")))
}
