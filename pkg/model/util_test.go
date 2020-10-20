package model

import (
	"testing"

	v1 "k8s.io/api/core/v1"

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

func TestUtil_testMergeEnvs(t *testing.T) {
	//given
	a := []v1.EnvVar{{
		Name:  "a",
		Value: "a",
	}}
	b := []v1.EnvVar{{
		Name:  "b",
		Value: "b",
	}}

	//when
	c := MergeEnvs(a, b)

	//then
	expected := []v1.EnvVar{
		{
			Name:  "a",
			Value: "a",
		},
		{
			Name:  "b",
			Value: "b",
		},
	}
	assert.Equal(t, expected, c)
}

func TestUtil_testMergeEnvsWithEmptyArguments(t *testing.T) {
	//given
	var a []v1.EnvVar
	var b []v1.EnvVar

	//when
	c := MergeEnvs(a, b)

	//then
	var expected []v1.EnvVar
	assert.Equal(t, expected, c)
}

func TestUtil_testMergeEnvsWithDuplicates1(t *testing.T) {
	//given
	a := []v1.EnvVar{{
		Name:  "a",
		Value: "a",
	}}
	b := []v1.EnvVar{{
		Name:  "a",
		Value: "b",
	}}

	//when
	c := MergeEnvs(a, b)

	//then
	expected := []v1.EnvVar{
		{
			Name:  "a",
			Value: "a",
		},
	}
	assert.Equal(t, expected, c)
}
