package api

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

func CallListSecretsV3(httpClient *resty.Client, request ListSecretsRequest) (ListSecretsResponse, error) {

	secretsResponse := ListSecretsResponse{}

	req := httpClient.R().
		SetResult(&secretsResponse).
		SetQueryParams(map[string]string{
			"workspaceId":            request.ProjectId,
			"workspaceSlug":          request.ProjectSlug,
			"environment":            request.Environment,
			"expandSecretReferences": fmt.Sprintf("%t", request.ExpandSecretReferences),
			"include_imports":        fmt.Sprintf("%t", request.IncludeImports),
			"recursive":              fmt.Sprintf("%t", request.Recursive),
		})

	if request.SecretPath != "" {
		req.SetQueryParam("secretPath", request.SecretPath)
	}

	res, err := req.Get("/v3/secrets/raw")

	if err != nil {
		return ListSecretsResponse{}, fmt.Errorf("CallListSecretsV3: Unable to complete api request [err=%s]", err)
	}

	if res.IsError() {
		return ListSecretsResponse{}, fmt.Errorf("CallListSecretsV3: Unsuccessful response [%v %v] [status-code=%v]", res.Request.Method, res.Request.URL, res.StatusCode())
	}

	return secretsResponse, nil
}
