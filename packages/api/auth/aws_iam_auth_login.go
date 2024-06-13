package api

import (
	"github.com/go-resty/resty/v2"
	"github.com/infisical/go-sdk/packages/errors"
)

const callAWSIamAuthLoginOperation = "CallAWSIamAuthLogin"

func CallAWSIamAuthLogin(httpClient *resty.Client, request AwsIamAuthLoginRequest) (accessToken string, e error) {
	var responseData GenericAuthLoginResponse

	response, err := httpClient.R().
		SetResult(&responseData).
		SetBody(request).
		Post("/v1/auth/aws-auth/login")

	if err != nil {
		return "", errors.NewRequestError(callAWSIamAuthLoginOperation, err)
	}

	if response.IsError() {
		return "", errors.NewAPIError(callAWSIamAuthLoginOperation, response)
	}

	return responseData.AccessToken, nil
}
