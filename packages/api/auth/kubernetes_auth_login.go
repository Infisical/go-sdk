package api

import (
	"github.com/go-resty/resty/v2"
	"github.com/infisical/go-sdk/packages/errors"
	"github.com/infisical/go-sdk/packages/models"
)

const callKubernetesAuthLoginOperation = "CallKubernetesAuthLogin"

func CallKubernetesAuthLogin(httpClient *resty.Client, request KubernetesAuthLoginRequest) (credential models.MachineIdentityCredential, e error) {
	var responseData MachineIdentityAuthLoginResponse

	response, err := httpClient.R().
		SetResult(&responseData).
		SetBody(request).
		Post("/v1/auth/kubernetes-auth/login")

	if err != nil {
		return models.MachineIdentityCredential{}, errors.NewRequestError(callKubernetesAuthLoginOperation, err)
	}

	if response.IsError() {
		return models.MachineIdentityCredential{}, errors.NewAPIError(callKubernetesAuthLoginOperation, response)
	}

	return responseData.ToMachineIdentity(), nil
}
