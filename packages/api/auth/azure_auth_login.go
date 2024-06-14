package api

import (
	"github.com/go-resty/resty/v2"
	"github.com/infisical/go-sdk/packages/errors"
)

const azureAuthLoginOperation = "CallAzureAuthLogin"

func CallAzureAuthLogin(httpClient *resty.Client, request AzureAuthLoginRequest) (credential MachineIdentityAuthLoginResponse, e error) {
	var responseData MachineIdentityAuthLoginResponse

	response, err := httpClient.R().
		SetResult(&responseData).
		SetBody(request).
		Post("/v1/auth/azure-auth/login")

	if err != nil {
		return MachineIdentityAuthLoginResponse{}, errors.NewRequestError(azureAuthLoginOperation, err)
	}

	if response.IsError() {
		return MachineIdentityAuthLoginResponse{}, errors.NewAPIError(azureAuthLoginOperation, response)
	}

	return responseData, nil
}
