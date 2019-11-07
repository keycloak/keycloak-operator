package model

import (
	"fmt"
	"math/rand"
	"strings"

	v13 "k8s.io/api/apps/v1"
)

// Copy pasted from https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
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

		// Uppercase letters
		if ascii >= 65 && ascii <= 90 {
			sb.WriteRune(char)
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
