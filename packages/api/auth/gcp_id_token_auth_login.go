package api

import (
	"github.com/go-resty/resty/v2"
	"github.com/infisical/go-sdk/packages/errors"
)

const callGCPAuthLoginOperation = "CallGCPAuthLogin"

func CallGCPAuthLogin(httpClient *resty.Client, request GCPAuthLoginRequest) (accessToken string, e error) {
	var responseData GenericAuthLoginResponse

	response, err := httpClient.R().
		SetResult(&responseData).
		SetBody(request).
		Post("/v1/auth/gcp-auth/login")

	if err != nil {
		return "", errors.NewRequestError(callGCPAuthLoginOperation, err)
	}

	if response.IsError() {
		return "", errors.NewAPIError(callGCPAuthLoginOperation, response)
	}

	return responseData.AccessToken, nil
}
