package util

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/infisical/go-sdk/packages/models"
	"github.com/rs/zerolog"
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

func PrintWarning(logger zerolog.Logger, message string) {
	logger.Warn().Msgf("[Infisical] Warning: %v", message)
}

func EnsureUniqueSecretsByKey(secrets *[]models.Secret, skipUniqueKey bool) {
	secretMap := make(map[string]models.Secret)

	// Move secrets to a map to ensure uniqueness
	for _, secret := range *secrets {
		var key string
		if skipUniqueKey {
			// Create a composite key using both SecretPath and SecretKey
			key = secret.SecretPath + ":" + secret.SecretKey
		} else {
			// Use only SecretKey for global uniqueness
			key = secret.SecretKey
		}
		secretMap[key] = secret
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

// Helper function to sort the secrets by key so we can create a consistent output
func SortSecretsByKeys(secrets []models.Secret) []models.Secret {
	sort.Slice(secrets, func(i, j int) bool {
		return secrets[i].SecretKey < secrets[j].SecretKey
	})
	return secrets
}

/*
If the status code is 400, there will most likely always be a body.
The body is a json object with a message key. we need to try to parse it, but if it fails, we can just return an empty string.
But if the status code is 500, there may not be a body. if there is, it will be a json object with a message key. we need to try to parse it, but if it fails, we can just return an empty string
*/
func TryParseErrorBody(res *resty.Response) string {
	if res == nil || !res.IsError() {
		return ""
	}

	body := res.String()
	if body == "" {
		return ""
	}

	type ErrorResponse struct {
		Message string `json:"message"`
		ReqId   string `json:"reqId"`
	}

	if res.StatusCode() == 422 {
		if issues, ok := formatValidationIssues(body); ok {
			return issues
		}
		return body
	}

	// now we have a string, we need to try to parse it as json
	var errorResponse ErrorResponse
	err := json.Unmarshal([]byte(body), &errorResponse)

	if err != nil {
		return ""
	}

	return errorResponse.Message
}

// formatValidationIssues turns the issues of a 422 response, such as
// {"message":[{"path":["services",0,"type"],"message":"Invalid discriminator value"}]}, into one entry per issue,
// separated by semicolons: "services.0.type: Invalid discriminator value". It reports false for a body in any other
// shape.
func formatValidationIssues(body string) (string, bool) {
	var response struct {
		Message json.RawMessage `json:"message"`
	}
	if err := json.Unmarshal([]byte(body), &response); err != nil || len(response.Message) == 0 {
		return "", false
	}

	var message string
	if err := json.Unmarshal(response.Message, &message); err == nil {
		return message, message != ""
	}

	var issues []struct {
		Path    []any  `json:"path"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(response.Message, &issues); err != nil || len(issues) == 0 {
		return "", false
	}

	lines := make([]string, 0, len(issues))
	for _, issue := range issues {
		if issue.Message == "" {
			return "", false
		}
		path := make([]string, len(issue.Path))
		for i, segment := range issue.Path {
			path[i] = fmt.Sprint(segment)
		}
		if len(path) == 0 {
			lines = append(lines, issue.Message)
		} else {
			lines = append(lines, strings.Join(path, ".")+": "+issue.Message)
		}
	}
	return strings.Join(lines, "; "), true
}

func TryExtractReqId(res *resty.Response) string {
	if res == nil || !res.IsError() {
		return ""
	}

	type ErrorResponse struct {
		ReqId string `json:"reqId"`
	}

	var errorResponse ErrorResponse

	err := json.Unmarshal([]byte(res.String()), &errorResponse)

	if err != nil {
		return ""
	}

	return errorResponse.ReqId
}

func SleepWithContext(ctx context.Context, duration time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(duration):
		return nil
	}
}

func ComputeCacheKeyFromBytes(bytes []byte, feature string) string {
	key := sha256.Sum256(bytes)
	return fmt.Sprintf("%s-%s", feature, hex.EncodeToString(key[:]))
}
