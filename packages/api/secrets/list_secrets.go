package api

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/infisical/go-sdk/packages/util"
)

func CallListSecretsV3(httpClient *resty.Client, request ListSecretsV3RawRequest) (ListSecretsV3RawResponse, error) {

	secretsResponse := ListSecretsV3RawResponse{}

	if request.SecretPath == "" {
		request.SecretPath = "/"
	}

	res, err := httpClient.R().
		SetResult(&secretsResponse).
		SetQueryParams(map[string]string{
			"workspaceId":            request.ProjectID,
			"workspaceSlug":          request.ProjectSlug,
			"environment":            request.Environment,
			"secretPath":             request.SecretPath,
			"expandSecretReferences": fmt.Sprintf("%t", request.ExpandSecretReferences),
			"include_imports":        fmt.Sprintf("%t", request.IncludeImports),
			"recursive":              fmt.Sprintf("%t", request.Recursive),
		}).Get("/v3/secrets/raw")

	if err != nil {
		return ListSecretsV3RawResponse{}, fmt.Errorf("CallListSecretsV3: Unable to complete api request [err=%s]", err)
	}

	if res.IsError() {
		return ListSecretsV3RawResponse{}, fmt.Errorf("CallListSecretsV3: Unsuccessful response [%v %v] [status-code=%v] %s", res.Request.Method, res.Request.URL, res.StatusCode(), fmt.Sprintf("Error: %s", util.TryParseErrorBody(res)))
	}

	return secretsResponse, nil
}
