package api

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/infisical/go-sdk/packages/errors"
)

const callListSecretsV3RawOperation = "CallListSecretsV3Raw"

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
		return ListSecretsV3RawResponse{}, errors.NewRequestError(callListSecretsV3RawOperation, err)
	}

	if res.IsError() {
		return ListSecretsV3RawResponse{}, errors.NewAPIErrorWithResponse(callListSecretsV3RawOperation, res)
	}

	return secretsResponse, nil
}
