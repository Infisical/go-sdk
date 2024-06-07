package util

import (
	"encoding/json"
	"strings"

	"github.com/go-resty/resty/v2"
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

/*
If the status code is 400, there will most likely always be a body.
The body is a json object with a message key. we need to try to parse it, but if it fails, we can just return an empty string.
But if the status code is 500, there may not be a body. if there is, it will be a json object with a message key. we need to try to parse it, but if it fails, we can just return an empty string
*/
func TryParseErrorBody(res *resty.Response) string {
	if res == nil || !res.IsError() {
		return "No Error Message"
	}

	body := res.String()
	if body == "" {
		return "No Error Message"
	}

	type ErrorResponse struct {
		Message string `json:"message"`
	}

	// now we have a string, we need to try to parse it as json
	var errorResponse ErrorResponse
	err := json.Unmarshal([]byte(body), &errorResponse)

	if err != nil {
		return "No Error Message"
	}

	return errorResponse.Message

}
