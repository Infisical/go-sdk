package api

import (
	"github.com/go-resty/resty/v2"
	"github.com/infisical/go-sdk/packages/errors"
	"github.com/infisical/go-sdk/packages/models"
)

const callTokenRenewOperation = "CallTokenRenew"

func CallTokenRenew(httpClient *resty.Client, request TokenRenewRequest) (credential models.MachineIdentityCredential, e error) {
	var responseData MachineIdentityAuthLoginResponse

	response, err := httpClient.R().
		SetResult(&responseData).
		SetBody(request).
		Post("/v1/auth/token/renew")

	if err != nil {
		return models.MachineIdentityCredential{}, errors.NewRequestError(callTokenRenewOperation, err)
	}

	if response.IsError() {
		return models.MachineIdentityCredential{}, errors.NewAPIError(callTokenRenewOperation, response)
	}

	return responseData.ToMachineIdentity(), nil
}
