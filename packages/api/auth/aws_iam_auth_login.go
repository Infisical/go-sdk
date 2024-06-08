package api

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

func CallAWSIamAuthLogin(httpClient *resty.Client, request AwsIamAuthLoginRequest) (accessToken string, e error) {
	var responseData GenericAuthLoginResponse

	response, err := httpClient.R().
		SetResult(&responseData).
		SetBody(request).
		Post("/v1/auth/aws-iam/login")

	if err != nil {
		return "", fmt.Errorf("CallAWSIamAuthLogin: Unable to complete api request [err=%s]", err)
	}

	if response.IsError() {
		return "", fmt.Errorf("CallAWSIamAuthLogin: Unsuccessful response [%v %v] [status-code=%v]", response.Request.Method, response.Request.URL, response.StatusCode())
	}

	return responseData.AccessToken, nil
}
