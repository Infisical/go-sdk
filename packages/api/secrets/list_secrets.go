package api

import (
	"github.com/go-resty/resty/v2"
)

func CallListSecretsV3(httpClient *resty.Client, request ListSecretsOptions) (ListSecretsResponse, error) {

	secretsResponse := ListSecretsResponse{}

	response, err := httpClient.R().
		SetResult(&secretsResponse).
		SetBody(request).
		Post("/v1/auth/universal-auth/login")

	if err != nil {
		return ListSecretsResponse{}, err
	}

	if response.IsError() {
		return ListSecretsResponse{}, response.Error().(error)
	}

	return secretsResponse, nil
}
