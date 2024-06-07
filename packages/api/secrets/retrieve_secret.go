package api

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/infisical/go-sdk/packages/util"
)

func CallRetrieveSecretV3(httpClient *resty.Client, request RetrieveSecretV3RawRequest) (RetrieveSecretV3RawResponse, error) {

	retrieveResponse := RetrieveSecretV3RawResponse{}

	if request.Type == "" {
		request.Type = "shared"
	}

	if request.SecretPath == "" {
		request.SecretPath = "/"
	}

	req := httpClient.R().
		SetResult(&retrieveResponse).
		SetQueryParams(map[string]string{
			"workspaceId":     request.ProjectID,
			"environment":     request.Environment,
			"secretPath":      request.SecretPath,
			"include_imports": fmt.Sprintf("%t", request.IncludeImports),
			"type":            request.Type,
		})

	res, err := req.Get(fmt.Sprintf("/v3/secrets/raw/%s", request.SecretKey))

	if err != nil {
		return RetrieveSecretV3RawResponse{}, fmt.Errorf("CallRetrieveSecretV3: Unable to complete api request [err=%s]", err)
	}

	if res.IsError() {
		return RetrieveSecretV3RawResponse{}, fmt.Errorf("CallRetrieveSecretV3: Unsuccessful response [%v %v] [status-code=%v] %s", res.Request.Method, res.Request.URL, res.StatusCode(), fmt.Sprintf("Error: %s", util.TryParseErrorBody(res)))
	}

	return retrieveResponse, nil
}
