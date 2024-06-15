package api

import (
	"github.com/go-resty/resty/v2"
	"github.com/infisical/go-sdk/packages/errors"
	"github.com/infisical/go-sdk/packages/models"
)

const callUniversalAuthLoginOperation = "CallUniversalAuthLogin"

func CallUniversalAuthLogin(httpClient *resty.Client, request UniversalAuthLoginRequest) (credential models.MachineIdentityCredential, e error) {
	var responseData MachineIdentityAuthLoginResponse

	response, err := httpClient.R().
		SetResult(&responseData).
		SetBody(request).
		Post("/v1/auth/universal-auth/login")

	if err != nil {
		return models.MachineIdentityCredential{}, errors.NewRequestError(callUniversalAuthLoginOperation, err)
	}

	if response.IsError() {
		return models.MachineIdentityCredential{}, errors.NewAPIError(callUniversalAuthLoginOperation, response)
	}

	return responseData.ToMachineIdentity(), nil
}
