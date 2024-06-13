package api

import (
	"github.com/go-resty/resty/v2"
	"github.com/infisical/go-sdk/packages/errors"
)

const azureAuthLoginOperation = "CallAzureAuthLogin"

func CallAzureAuthLogin(httpClient *resty.Client, request AzureAuthLoginRequest) (accessToken string, e error) {
	var responseData GenericAuthLoginResponse

	response, err := httpClient.R().
		SetResult(&responseData).
		SetBody(request).
		Post("/v1/auth/azure-auth/login")

	if err != nil {
		return "", errors.NewRequestError(azureAuthLoginOperation, err)
	}

	if response.IsError() {
		return "", errors.NewAPIError(azureAuthLoginOperation, response)
	}

	return responseData.AccessToken, nil
}
