package util

import (
	"strings"

	"github.com/infisical/go-sdk/packages/models"
)

func AppendAPIEndpoint(siteUrl string) string {
	// Ensure the address does not already end with "/api"
	if strings.HasSuffix(siteUrl, "/api") {
		return siteUrl
	}

	// Check if the address ends with a slash and append accordingly
	if siteUrl[len(siteUrl)-1] == '/' {
		return siteUrl + "api"
	}
	return siteUrl + "/api"
}

func EnsureUniqueSecretsByKey(secrets *[]models.Secret) {
	secretMap := make(map[string]models.Secret)

	// Move secrets to a map to ensure uniqueness
	for _, secret := range *secrets {
		secretMap[secret.SecretKey] = secret // Maps will overwrite existing entry with the same key
	}

	// Clear the slice
	*secrets = (*secrets)[:0]

	// Refill the slice from the map
	for _, secret := range secretMap {
		*secrets = append(*secrets, secret)
	}
}

// containsSecret checks if the given key exists in the slice of secrets
func ContainsSecret(secrets []models.Secret, key string) bool {
	for _, secret := range secrets {
		if secret.SecretKey == key {
			return true
		}
	}
	return false
}
