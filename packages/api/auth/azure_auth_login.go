package api

import (
	"github.com/go-resty/resty/v2"
	"github.com/infisical/go-sdk/packages/errors"
	"github.com/infisical/go-sdk/packages/models"
)

const azureAuthLoginOperation = "CallAzureAuthLogin"

func CallAzureAuthLogin(httpClient *resty.Client, request AzureAuthLoginRequest) (credential models.MachineIdentityCredential, e error) {
	var responseData MachineIdentityAuthLoginResponse

	response, err := httpClient.R().
		SetResult(&responseData).
		SetBody(request).
		Post("/v1/auth/azure-auth/login")

	if err != nil {
		return models.MachineIdentityCredential{}, errors.NewRequestError(azureAuthLoginOperation, err)
	}

	if response.IsError() {
		return models.MachineIdentityCredential{}, errors.NewAPIError(azureAuthLoginOperation, response)
	}

	return responseData.ToMachineIdentity(), nil
}
