package model

import (
	"testing"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/stretchr/testify/assert"
	v13 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
)

type createDeploymentStatefulSet func(*v1alpha1.Keycloak, *v1.Secret) *v13.StatefulSet

func TestKeycloakDeployment_testExperimentalEnvs(t *testing.T) {
	testExperimentalEnvs(t, KeycloakDeployment)
}

func TestKeycloakDeployment_testExperimentalArgs(t *testing.T) {
	testExperimentalArgs(t, KeycloakDeployment)
}

func TestKeycloakDeployment_testExperimentalCommand(t *testing.T) {
	testExperimentalCommand(t, KeycloakDeployment)
}

func TestKeycloakDeployment_testExperimentalVolumesWithConfigMaps(t *testing.T) {
	testExperimentalVolumesWithConfigMaps(t, KeycloakDeployment)
}

func testExperimentalEnvs(t *testing.T, deploymentFunction createDeploymentStatefulSet) {
	//given
	dbSecret := &v1.Secret{}
	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			KeycloakDeploymentSpec: v1alpha1.KeycloakDeploymentSpec{
				Experimental: v1alpha1.ExperimentalSpec{
					Env: []v1.EnvVar{
						{
							// New value
							Name:  "testName",
							Value: "testValue",
						},
						{
							// An overridden value
							Name:  "DB_SCHEMA",
							Value: "test",
						},
					},
				},
			},
		},
	}

	//when
	envs := deploymentFunction(cr, dbSecret).Spec.Template.Spec.Containers[0].Env

	//then
	hasTestNameKey := false
	testNameValue := ""
	dbSchemaValue := ""
	for _, v := range envs {
		if v.Name == "testName" {
			hasTestNameKey = true
			testNameValue = v.Value
		}
		if v.Name == "DB_SCHEMA" {
			dbSchemaValue = v.Value
		}
	}
	assert.True(t, hasTestNameKey)
	assert.Equal(t, "testValue", testNameValue)
	assert.Equal(t, "test", dbSchemaValue)
	assert.True(t, len(envs) > 1)
}

func testExperimentalArgs(t *testing.T, deploymentFunction createDeploymentStatefulSet) {
	//given
	dbSecret := &v1.Secret{}
	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			KeycloakDeploymentSpec: v1alpha1.KeycloakDeploymentSpec{
				Experimental: v1alpha1.ExperimentalSpec{
					Args: []string{"test"},
				},
			},
		},
	}

	//when
	args := deploymentFunction(cr, dbSecret).Spec.Template.Spec.Containers[0].Args

	//then
	assert.Equal(t, []string{"test"}, args)
}

func testExperimentalCommand(t *testing.T, deploymentFunction createDeploymentStatefulSet) {
	//given
	dbSecret := &v1.Secret{}
	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			KeycloakDeploymentSpec: v1alpha1.KeycloakDeploymentSpec{
				Experimental: v1alpha1.ExperimentalSpec{
					Command: []string{"test"},
				},
			},
		},
	}

	//when
	command := deploymentFunction(cr, dbSecret).Spec.Template.Spec.Containers[0].Command

	//then
	assert.Equal(t, []string{"test"}, command)
}

func testExperimentalVolumesWithConfigMaps(t *testing.T, deploymentFunction createDeploymentStatefulSet) {
	//given
	dbSecret := &v1.Secret{}
	cr := &v1alpha1.Keycloak{
		Spec: v1alpha1.KeycloakSpec{
			KeycloakDeploymentSpec: v1alpha1.KeycloakDeploymentSpec{
				Experimental: v1alpha1.ExperimentalSpec{
					Volumes: v1alpha1.VolumesSpec{
						Items: []v1alpha1.VolumeSpec{
							{
								ConfigMap: &v1alpha1.ConfigMapVolumeSpec{
									Name:      "testName",
									MountPath: "testMountPath",
									Items: []v1.KeyToPath{
										{
											Key:  "testKey",
											Path: "testPath",
										},
									},
								},
							},
						},
						DefaultMode: &[]int32{1}[0],
					},
				},
			},
		},
	}

	//when
	template := deploymentFunction(cr, dbSecret).Spec.Template.Spec
	volumeMount := template.Containers[0].VolumeMounts[3]
	volume := template.Volumes[3]

	//then
	assert.Equal(t, "testName", volumeMount.Name)
	assert.Equal(t, "testMountPath", volumeMount.MountPath)
	assert.Equal(t, "testName", volume.Name)
	assert.Equal(t, "testName", volume.Projected.Sources[0].ConfigMap.Name)
	assert.Equal(t, "testKey", volume.Projected.Sources[0].ConfigMap.Items[0].Key)
	assert.Equal(t, "testPath", volume.Projected.Sources[0].ConfigMap.Items[0].Path)
}
