package api

import (
	"github.com/go-resty/resty/v2"
	"github.com/infisical/go-sdk/packages/errors"
)

const callUniversalAuthLoginOperation = "CallUniversalAuthLogin"

func CallUniversalAuthLogin(httpClient *resty.Client, request UniversalAuthLoginRequest) (accessToken string, e error) {
	var responseData GenericAuthLoginResponse

	response, err := httpClient.R().
		SetResult(&responseData).
		SetBody(request).
		Post("/v1/auth/universal-auth/login")

	if err != nil {
		return "", errors.NewRequestError(callUniversalAuthLoginOperation, err)
	}

	if response.IsError() {
		return "", errors.NewAPIError(callUniversalAuthLoginOperation, response)
	}

	return responseData.AccessToken, nil
}
