package api

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/infisical/go-sdk/packages/util"
)

func CallUpdateSecretV3(httpClient *resty.Client, request UpdateSecretV3RawRequest) (UpdateSecretV3RawResponse, error) {

	updateResponse := UpdateSecretV3RawResponse{}

	req := httpClient.R().
		SetResult(&updateResponse).
		SetBody(request)

	res, err := req.Patch(fmt.Sprintf("/v3/secrets/raw/%s", request.SecretKey))

	if err != nil {
		return UpdateSecretV3RawResponse{}, fmt.Errorf("CallUpdateSecretV3: Unable to complete api request [err=%s]", err)
	}

	if res.IsError() {
		return UpdateSecretV3RawResponse{}, fmt.Errorf("CallUpdateSecretV3: Unsuccessful response [%v %v] [status-code=%v] %s", res.Request.Method, res.Request.URL, res.StatusCode(), fmt.Sprintf("Error: %s", util.TryParseErrorBody(res)))
	}

	return updateResponse, nil
}
