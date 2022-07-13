package e2e

import "flag"

var (
	isProductBuild bool
)

func init() {
	flag.BoolVar(&isProductBuild, "product", false, "Using RHSSO or Keycloak")
}

const keycloakProfile = "keycloak"
const rhssoProfile = "RHSSO"

func currentProfile() string {
	if isProductBuild {
		return rhssoProfile
	}
	return keycloakProfile
}
