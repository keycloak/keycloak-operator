package model

import (
	"testing"
)

func TestRHSSODeployment_testExperimentalEnvs(t *testing.T) {
	testExperimentalEnvs(t, RHSSODeployment)
}

func TestRHSSODeployment_testExperimentalArgs(t *testing.T) {
	testExperimentalArgs(t, RHSSODeployment)
}

func TestRHSSODeployment_testExperimentalCommand(t *testing.T) {
	testExperimentalCommand(t, RHSSODeployment)
}

func TestRHSSODeployment_testExperimentalVolumesWithConfigMaps(t *testing.T) {
	testExperimentalVolumesWithConfigMaps(t, RHSSODeployment)
}
