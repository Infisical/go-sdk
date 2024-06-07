package api

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

func CallGCPAuthLogin(httpClient *resty.Client, request GCPAuthLoginRequest) (accessToken string, e error) {
	var responseData GenericAuthLoginResponse

	response, err := httpClient.R().
		SetResult(&responseData).
		SetBody(request).
		Post("/v1/auth/gcp-auth/login")

	if err != nil {
		return "", fmt.Errorf("CallGCPIdTokenAuthLogin: Unable to complete api request [err=%s]", err)
	}

	if response.IsError() {
		return "", fmt.Errorf("CallGCPIdTokenAuthLogin: Unsuccessful response [%v %v] [status-code=%v]", response.Request.Method, response.Request.URL, response.StatusCode())
	}

	return responseData.AccessToken, nil
}
