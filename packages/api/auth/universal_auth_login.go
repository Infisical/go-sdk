package api

import (
	"github.com/go-resty/resty/v2"
	"github.com/infisical/go-sdk/packages/errors"
)

const callUniversalAuthLoginOperation = "CallUniversalAuthLogin"

func CallUniversalAuthLogin(httpClient *resty.Client, request UniversalAuthLoginRequest) (credential MachineIdentityAuthLoginResponse, e error) {
	var responseData MachineIdentityAuthLoginResponse

	response, err := httpClient.R().
		SetResult(&responseData).
		SetBody(request).
		Post("/v1/auth/universal-auth/login")

	if err != nil {
		return responseData, errors.NewRequestError(callUniversalAuthLoginOperation, err)
	}

	if response.IsError() {
		return responseData, errors.NewAPIError(callUniversalAuthLoginOperation, response)
	}

	return responseData, nil
}
