package api

import (
	"github.com/go-resty/resty/v2"
	"github.com/infisical/go-sdk/packages/errors"
)

const callKubernetesAuthLoginOperation = "CallKubernetesAuthLogin"

func CallKubernetesAuthLogin(httpClient *resty.Client, request KubernetesAuthLoginRequest) (accessToken string, e error) {
	var responseData GenericAuthLoginResponse

	response, err := httpClient.R().
		SetResult(&responseData).
		SetBody(request).
		Post("/v1/auth/kubernetes-auth/login")

	if err != nil {
		return "", errors.NewRequestError(callKubernetesAuthLoginOperation, err)
	}

	if response.IsError() {
		return "", errors.NewAPIError(callKubernetesAuthLoginOperation, response)
	}

	return responseData.AccessToken, nil
}
