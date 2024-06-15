package api

import (
	"github.com/go-resty/resty/v2"
	"github.com/infisical/go-sdk/packages/errors"
	"github.com/infisical/go-sdk/packages/models"
)

const callAWSIamAuthLoginOperation = "CallAWSIamAuthLogin"

func CallAWSIamAuthLogin(httpClient *resty.Client, request AwsIamAuthLoginRequest) (credential models.MachineIdentityCredential, e error) {
	var responseData MachineIdentityAuthLoginResponse

	response, err := httpClient.R().
		SetResult(&responseData).
		SetBody(request).
		Post("/v1/auth/aws-auth/login")

	if err != nil {
		return models.MachineIdentityCredential{}, errors.NewRequestError(callAWSIamAuthLoginOperation, err)
	}

	if response.IsError() {
		return models.MachineIdentityCredential{}, errors.NewAPIError(callAWSIamAuthLoginOperation, response)
	}

	return responseData.ToMachineIdentity(), nil
}
