package api

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/infisical/go-sdk/packages/util"
)

func CallCreateSecretV3(httpClient *resty.Client, request CreateSecretV3RawRequest) (CreateSecretV3RawResponse, error) {

	createResponse := CreateSecretV3RawResponse{}

	req := httpClient.R().
		SetResult(&createResponse).
		SetBody(request)

	res, err := req.Post(fmt.Sprintf("/v3/secrets/raw/%s", request.SecretKey))

	if err != nil {
		return CreateSecretV3RawResponse{}, fmt.Errorf("CallCreateSecretV3: Unable to complete api request [err=%s]", err)
	}

	if res.IsError() {
		return CreateSecretV3RawResponse{}, fmt.Errorf("CallCreateSecretV3: Unsuccessful response [%v %v] [status-code=%v] %s", res.Request.Method, res.Request.URL, res.StatusCode(), fmt.Sprintf("Error: %s", util.TryParseErrorBody(res)))
	}

	return createResponse, nil
}
