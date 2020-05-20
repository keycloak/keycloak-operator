package model

import (
	"os"
	"strings"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
)

const (
	RHSSOProfile                 = "RHSSO"
	ProfileEnvironmentalVariable = "PROFILE"
)

var Profiles = NewProfileManager()

type ProfileManager struct {
	Profiles []string
}

func NewProfileManager() ProfileManager {
	ret := ProfileManager{}
	ret.Profiles = ret.getProfiles()
	return ret
}

func (p *ProfileManager) IsRHSSO(cr *v1alpha1.Keycloak) bool {
	for _, profile := range p.Profiles {
		if profile == RHSSOProfile {
			return true
		}
	}
	if cr != nil && cr.Spec.Profile == RHSSOProfile {
		return true
	}
	return false
}

func (p *ProfileManager) GetKeycloakOrRHSSOImage(cr *v1alpha1.Keycloak) string {
	if p.IsRHSSO(cr) {
		return Images.Images[RHSSOImage]
	}
	return Images.Images[KeycloakImage]
}

func (p *ProfileManager) GetInitContainerImage(cr *v1alpha1.Keycloak) string {
	if p.IsRHSSO(cr) {
		return Images.Images[RHSSOInitContainer]
	}
	return Images.Images[KeycloakInitContainer]
}

func (p *ProfileManager) getProfiles() []string {
	env := os.Getenv(ProfileEnvironmentalVariable)
	if env == "" {
		return []string{}
	}
	return strings.Split(env, ",")
}
