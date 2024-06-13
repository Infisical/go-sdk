package api

import (
	"github.com/go-resty/resty/v2"
	"github.com/infisical/go-sdk/packages/errors"
	"github.com/infisical/go-sdk/packages/models"
)

const callUniversalAuthLoginOperation = "CallUniversalAuthLogin"

func CallUniversalAuthLogin(httpClient *resty.Client, request UniversalAuthLoginRequest) (models.UniversalAuthCredential, error) {
	var responseData models.UniversalAuthCredential

	response, err := httpClient.R().
		SetResult(&responseData).
		SetBody(request).
		Post("/v1/auth/universal-auth/login")

	if err != nil {
		return responseData, errors.NewRequestError(callUniversalAuthLoginOperation, err)
	}

	if response.IsError() {
		return responseData, errors.NewAPIError(callUniversalAuthLoginOperation, response)
	}

	return responseData, nil
}
