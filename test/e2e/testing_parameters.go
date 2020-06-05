package e2e

import "flag"

var (
	isProductBuild bool
)

func init() {
	flag.BoolVar(&isProductBuild, "product", false, "Using RHSSO or Keycloak")
}

func currentProfile() string {
	if isProductBuild {
		return "RHSSO"
	}
	return "keycloak"
}
