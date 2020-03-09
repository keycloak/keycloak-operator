package model

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net"
	"strconv"
	"strings"
	"unicode"

	v1 "k8s.io/api/core/v1"

	v13 "k8s.io/api/apps/v1"
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

// Get image string from the statefulset. Default to RHSSOImage string
func GetCurrentKeycloakImage(currentState *v13.StatefulSet) string {
	for _, ele := range currentState.Spec.Template.Spec.Containers {
		if ele.Name == KeycloakDeploymentName {
			return ele.Image
		}
	}
	return RHSSOImage
}

// Split a full image string (e.g. quay.io/keycloak/keycloak:7.0.1 or registry.access.redhat.com/redhat-sso-7/sso73-openshift:1.0 ) into it's repo and individual versions
func GetImageRepoAndVersion(image string) (string, string, string, string) {
	imageRepo, imageMajor, imageMinor, imagePatch := "", "", "", ""

	// Split the string on : which will leave the repo and tag
	imageStrings := strings.Split(image, ":")

	if len(imageStrings) > 0 {
		imageRepo = imageStrings[0]
	}

	// If somehow the tag doesn't exist, return with empty strings for the versions
	if len(imageStrings) == 1 {
		return imageRepo, imageMajor, imageMinor, imagePatch
	}

	// Split the image tag on . to separate the version numbers
	imageTagStrings := strings.Split(imageStrings[1], ".")

	if len(imageTagStrings) > 0 {
		imageMajor = imageTagStrings[0]
	}

	if len(imageTagStrings) > 1 {
		imageMinor = imageTagStrings[1]
	}

	if len(imageTagStrings) > 2 {
		imagePatch = imageTagStrings[2]
	}

	return imageRepo, imageMajor, imageMinor, imagePatch
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
	parsed, err := strconv.Atoi(string(port))
	if err != nil {
		return PostgresDefaultPort
	}
	return int32(parsed)
}
