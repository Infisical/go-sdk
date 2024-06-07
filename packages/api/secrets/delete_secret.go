package api

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/infisical/go-sdk/packages/util"
)

func CallDeleteSecretV3(httpClient *resty.Client, request DeleteSecretV3RawRequest) (DeleteSecretV3RawResponse, error) {

	deleteResponse := DeleteSecretV3RawResponse{}

	req := httpClient.R().
		SetResult(&deleteResponse).
		SetBody(request)

	res, err := req.Delete(fmt.Sprintf("/v3/secrets/raw/%s", request.SecretKey))

	if err != nil {
		return DeleteSecretV3RawResponse{}, fmt.Errorf("CallDeleteSecretV3: Unable to complete api request [err=%s]", err)
	}

	if res.IsError() {
		return DeleteSecretV3RawResponse{}, fmt.Errorf("CallDeleteSecretV3: Unsuccessful response [%v %v] [status-code=%v] %s", res.Request.Method, res.Request.URL, res.StatusCode(), fmt.Sprintf("Error: %s", util.TryParseErrorBody(res)))
	}

	return deleteResponse, nil
}
