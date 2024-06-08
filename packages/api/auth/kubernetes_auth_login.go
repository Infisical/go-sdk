package api

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

func CallKubernetesAuthLogin(httpClient *resty.Client, request KubernetesAuthLoginRequest) (accessToken string, e error) {
	var responseData GenericAuthLoginResponse

	response, err := httpClient.R().
		SetResult(&responseData).
		SetBody(request).
		Post("/v1/auth/kubernetes-auth/login")

	if err != nil {
		return "", fmt.Errorf("CallKubernetesAuthLogin: Unable to complete api request [err=%s]", err)
	}

	if response.IsError() {
		return "", fmt.Errorf("CallKubernetesAuthLogin: Unsuccessful response [%v %v] [status-code=%v]", response.Request.Method, response.Request.URL, response.StatusCode())
	}

	return responseData.AccessToken, nil
}
