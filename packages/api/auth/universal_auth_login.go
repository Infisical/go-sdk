package api

import (
	"github.com/go-resty/resty/v2"
)

func CallUniversalAuthLogin(httpClient *resty.Client, request UniversalAuthLoginRequest) (accessToken string, e error) {
	var responseData UniversalAuthLoginResponse

	response, err := httpClient.R().
		SetResult(&responseData).
		SetBody(request).
		Post("/v1/auth/universal-auth/login")

	if err != nil {
		return "", err
	}

	if response.IsError() {
		return "", response.Error().(error)
	}

	return responseData.AccessToken, nil
}
