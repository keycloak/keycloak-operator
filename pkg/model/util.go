package model

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net"
	"strconv"
	"strings"
	"unicode"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	v1 "k8s.io/api/core/v1"
)

// Copy pasted from https://blog.questionable.services/article/generating-secure-random-numbers-crypto-rand/

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(s int) string {
	b := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b)
}

func GetRealmUserSecretName(keycloakNamespace, realmName, userName string) string {
	return SanitizeResourceName(fmt.Sprintf("credential-%v-%v-%v",
		realmName,
		userName,
		keycloakNamespace))
}

func SanitizeNumberOfReplicas(numberOfReplicas int, isCreate bool) *int32 {
	numberOfReplicasCasted := int32(numberOfReplicas)
	if isCreate && numberOfReplicasCasted < 1 {
		numberOfReplicasCasted = 1
	}
	return &[]int32{numberOfReplicasCasted}[0]
}

func SanitizeResourceName(name string) string {
	sb := strings.Builder{}
	for _, char := range name {
		ascii := int(char)
		// number
		if ascii >= 48 && ascii <= 57 {
			sb.WriteRune(char)
			continue
		}

		// Uppercase letters are transformed to lowercase
		if ascii >= 65 && ascii <= 90 {
			sb.WriteRune(unicode.ToLower(char))
			continue
		}

		// Lowercase letters
		if ascii >= 97 && ascii <= 122 {
			sb.WriteRune(char)
			continue
		}

		// dash
		if ascii == 45 {
			sb.WriteRune(char)
			continue
		}

		if char == '.' {
			sb.WriteRune(char)
			continue
		}

		if ascii == '_' {
			sb.WriteRune('-')
			continue
		}

		// Ignore all invalid chars
		continue
	}

	return sb.String()
}

func IsIP(host []byte) bool {
	return net.ParseIP(string(host)) != nil
}

func GetExternalDatabaseHost(secret *v1.Secret) string {
	host := secret.Data[DatabaseSecretExternalAddressProperty]
	return string(host)
}

func GetExternalDatabaseName(secret *v1.Secret) string {
	if secret == nil {
		return PostgresqlDatabase
	}

	name := secret.Data[DatabaseSecretDatabaseProperty]
	return string(name)
}

func GetExternalDatabasePort(secret *v1.Secret) int32 {
	if secret == nil {
		return PostgresDefaultPort
	}

	port := secret.Data[DatabaseSecretExternalPortProperty]
	parsed, err := strconv.ParseInt(string(port), 10, 32)
	if err != nil {
		return PostgresDefaultPort
	}
	return int32(parsed)
}

// This function favors values in "a".
func MergeEnvs(a []v1.EnvVar, b []v1.EnvVar) []v1.EnvVar {
	for _, bb := range b {
		found := false
		for _, aa := range a {
			if aa.Name == bb.Name {
				aa.Value = bb.Value
				found = true
				break
			}
		}
		if !found {
			a = append(a, bb)
		}
	}
	return a
}

// returned roles are always from a
func RoleDifferenceIntersection(a []v1alpha1.RoleRepresentation, b []v1alpha1.RoleRepresentation) (d []v1alpha1.RoleRepresentation, i []v1alpha1.RoleRepresentation) {
	for _, role := range a {
		if hasMatchingRole(b, role) {
			i = append(i, role)
		} else {
			d = append(d, role)
		}
	}
	return d, i
}

func hasMatchingRole(roles []v1alpha1.RoleRepresentation, otherRole v1alpha1.RoleRepresentation) bool {
	for _, role := range roles {
		if roleMatches(role, otherRole) {
			return true
		}
	}
	return false
}

func roleMatches(a v1alpha1.RoleRepresentation, b v1alpha1.RoleRepresentation) bool {
	if a.ID != "" && b.ID != "" {
		return a.ID == b.ID
	}
	return a.Name == b.Name
}

// FIXME Find a better way to refactor this code with role difference part above
// returned clientScopes are always from a
func ClientScopeDifferenceIntersection(a []v1alpha1.KeycloakClientScope, b []v1alpha1.KeycloakClientScope) (d []v1alpha1.KeycloakClientScope, i []v1alpha1.KeycloakClientScope) {
	for _, clientScope := range a {
		if hasMatchingClientScope(b, clientScope) {
			i = append(i, clientScope)
		} else {
			d = append(d, clientScope)
		}
	}
	return d, i
}

func hasMatchingClientScope(clientScopes []v1alpha1.KeycloakClientScope, otherClientScope v1alpha1.KeycloakClientScope) bool {
	for _, clientScope := range clientScopes {
		if clientScopeMatches(clientScope, otherClientScope) {
			return true
		}
	}
	return false
}

func clientScopeMatches(a v1alpha1.KeycloakClientScope, b v1alpha1.KeycloakClientScope) bool {
	if a.ID != "" && b.ID != "" {
		return a.ID == b.ID
	}
	return a.Name == b.Name
}

func FilterClientScopesByNames(clientScopes []v1alpha1.KeycloakClientScope, names []string) (filteredScopes []v1alpha1.KeycloakClientScope) {
	hashMap := make(map[string]v1alpha1.KeycloakClientScope)

	for _, scope := range clientScopes {
		hashMap[scope.Name] = scope
	}

	for _, name := range names {
		if scope, retrieved := hashMap[name]; retrieved {
			filteredScopes = append(filteredScopes, scope)
		}
	}

	return filteredScopes
}
