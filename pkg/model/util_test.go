package model

import (
	"testing"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
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

func TestKeycloakClientReconciler_Test_Role_DifferenceIntersection(t *testing.T) {
	// given
	a := []v1alpha1.RoleRepresentation{
		{Name: "a"},
		{ID: "ignored", Name: "b"},
		{ID: "cID", Name: "c"},
	}
	b := []v1alpha1.RoleRepresentation{
		{Name: "b"},
		{ID: "cID", Name: "differentName"},
		{Name: "d"},
	}

	// when
	difference, intersection := RoleDifferenceIntersection(a, b)

	// then
	expectedDifference := []v1alpha1.RoleRepresentation{
		{Name: "a"},
	}
	expectedIntersection := []v1alpha1.RoleRepresentation{
		{ID: "ignored", Name: "b"},
		{ID: "cID", Name: "c"},
	}
	assert.Equal(t, expectedDifference, difference)
	assert.Equal(t, expectedIntersection, intersection)
}

func TestKeycloakClientReconciler_Test_ClientScope_DifferenceIntersection(t *testing.T) {
	// given
	a := []v1alpha1.KeycloakClientScope{
		{Name: "a"},
		{ID: "ignored", Name: "b"},
		{ID: "cID", Name: "c"},
	}
	b := []v1alpha1.KeycloakClientScope{
		{Name: "b"},
		{ID: "cID", Name: "differentName"},
		{Name: "d"},
	}

	// when
	difference, intersection := ClientScopeDifferenceIntersection(a, b)

	// then
	expectedDifference := []v1alpha1.KeycloakClientScope{
		{Name: "a"},
	}
	expectedIntersection := []v1alpha1.KeycloakClientScope{
		{ID: "ignored", Name: "b"},
		{ID: "cID", Name: "c"},
	}
	assert.Equal(t, expectedDifference, difference)
	assert.Equal(t, expectedIntersection, intersection)
}

func TestPodLabels_When_EnvVars_Then_FullListOfLabels(t *testing.T) {
	cr := v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			KeycloakDeploymentSpec: v1alpha1.KeycloakDeploymentSpec{
				PodLabels: map[string]string{
					"LabelToTest":       "thisisthelabelvalue",
					"SecondLabelToTest": "anotherthisisthelabelvalue",
				},
			},
		}}

	PodLabels = map[string]string{
		"FirstLabel":  "firstLabelValue",
		"SecondLabel": "secondLabelValue",
	}

	labels := map[string]string{
		"app":       ApplicationName,
		"component": KeycloakDeploymentComponent,
	}
	totalLabels := AddPodLabels(&cr, labels)
	assert.Equal(t, 6, len(totalLabels))
	assert.Equal(t, 2, len(labels))
	assert.Contains(t, totalLabels, "LabelToTest")
	assert.Contains(t, totalLabels, "SecondLabelToTest")
	assert.Contains(t, totalLabels, "FirstLabel")
	assert.Contains(t, totalLabels, "SecondLabel")
	assert.Contains(t, totalLabels, "app")
	assert.Contains(t, totalLabels, "component")
}
func TestPodAnnotations_When_EnvVars_Then_FullListOfAnnotations(t *testing.T) {
	cr := v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			KeycloakDeploymentSpec: v1alpha1.KeycloakDeploymentSpec{
				PodAnnotations: map[string]string{
					"AnnotationToTest":       "thisistheannotationvalue",
					"SecondAnnotationToTest": "anotherthisistheannotationvalue",
				},
			},
		}}

	annotations := map[string]string{
		"app":       ApplicationName,
		"component": KeycloakDeploymentComponent,
	}
	totalAnnotations := AddPodAnnotations(&cr, annotations)
	assert.Equal(t, 4, len(totalAnnotations))
	assert.Contains(t, totalAnnotations, "AnnotationToTest")
	assert.Contains(t, totalAnnotations, "SecondAnnotationToTest")
	assert.Contains(t, totalAnnotations, "app")
	assert.Contains(t, totalAnnotations, "component")
}
