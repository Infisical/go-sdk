package api

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

func CallAzureAuthLogin(httpClient *resty.Client, request AzureAuthLoginRequest) (accessToken string, e error) {
	var responseData GenericAuthLoginResponse

	response, err := httpClient.R().
		SetResult(&responseData).
		SetBody(request).
		Post("/v1/auth/azure-auth/login")

	if err != nil {
		return "", fmt.Errorf("CallAzureAuthLogin: Unable to complete api request [err=%s]", err)
	}

	if response.IsError() {
		return "", fmt.Errorf("CallAzureAuthLogin: Unsuccessful response [%v %v] [status-code=%v]", response.Request.Method, response.Request.URL, response.StatusCode())
	}

	return responseData.AccessToken, nil
}
