package model

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	logr "log"
	"net"
	"runtime"
	"strconv"
	"strings"
	"unicode"

	"github.com/google/go-cmp/cmp"
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

func GetCallerFunction(depth int) string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(depth, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return frame.Function
}

func LogDiff(was, now interface{}) {
	if diff := cmp.Diff(was, now); diff != "" {
		logr.Printf("%s made changes (-was, +now):\n%s\n", GetCallerFunction(3), diff)
	}
}

func LogHasDiff(was, now interface{}) {
	if !cmp.Equal(was, now) {
		logr.Printf("%s has changes", GetCallerFunction(3))
	}
}

func GetServiceAccountUsername(clientName string) string {
	return fmt.Sprintf("service-account-%s", clientName)
}
