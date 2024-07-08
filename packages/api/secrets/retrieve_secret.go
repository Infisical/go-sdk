package api

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/infisical/go-sdk/packages/errors"
)

const callRetrieveSecretV3RawOperation = "CallRetrieveSecretV3Raw"

func CallRetrieveSecretV3(httpClient *resty.Client, request RetrieveSecretV3RawRequest) (RetrieveSecretV3RawResponse, error) {

	retrieveResponse := RetrieveSecretV3RawResponse{}

	if request.Type == "" {
		request.Type = "shared"
	}

	if request.SecretPath == "" {
		request.SecretPath = "/"
	}

	var version string
	if request.Version > 0 {
		version = fmt.Sprintf("%d", request.Version)
	}

	req := httpClient.R().
		SetResult(&retrieveResponse).
		SetQueryParams(map[string]string{
			"workspaceId":     request.ProjectID,
			"environment":     request.Environment,
			"secretPath":      request.SecretPath,
			"include_imports": fmt.Sprintf("%t", request.IncludeImports),
			"type":            request.Type,
			"version":         version,
		})

	res, err := req.Get(fmt.Sprintf("/v3/secrets/raw/%s", request.SecretKey))

	if err != nil {
		return RetrieveSecretV3RawResponse{}, errors.NewRequestError(callRetrieveSecretV3RawOperation, err)
	}

	if res.IsError() {
		return RetrieveSecretV3RawResponse{}, errors.NewAPIErrorWithResponse(callRetrieveSecretV3RawOperation, res)
	}

	return retrieveResponse, nil
}
