package model

import (
	"fmt"
	"math/rand"
	"strings"
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
