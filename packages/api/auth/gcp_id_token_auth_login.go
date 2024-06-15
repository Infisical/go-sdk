package api

import (
	"github.com/go-resty/resty/v2"
	"github.com/infisical/go-sdk/packages/errors"
	"github.com/infisical/go-sdk/packages/models"
)

const callGCPAuthLoginOperation = "CallGCPAuthLogin"

func CallGCPAuthLogin(httpClient *resty.Client, request GCPAuthLoginRequest) (credential models.MachineIdentityCredential, e error) {
	var responseData MachineIdentityAuthLoginResponse

	response, err := httpClient.R().
		SetResult(&responseData).
		SetBody(request).
		Post("/v1/auth/gcp-auth/login")

	if err != nil {
		return models.MachineIdentityCredential{}, errors.NewRequestError(callGCPAuthLoginOperation, err)
	}

	if response.IsError() {
		return models.MachineIdentityCredential{}, errors.NewAPIError(callGCPAuthLoginOperation, response)
	}

	return responseData.ToMachineIdentity(), nil
}
