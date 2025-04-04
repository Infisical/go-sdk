package api

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/infisical/go-sdk/packages/errors"
)

const callKmsGetKeyByNameOperationV1 = "CallKmsGetKeyByNameV1"

func CallKmsGetKeyByNameV1(httpClient *resty.Client, request KmsGetKeyByNameV1Request) (KmsGetKeyV1Response, error) {
	kmsGetKeyByNameResponse := KmsGetKeyV1Response{}

	res, err := httpClient.R().
		SetResult(&kmsGetKeyByNameResponse).
		Get(fmt.Sprintf("/v1/kms/key?projectId=%s&name=%s", request.ProjectId, request.KeyName))

	if err != nil {
		return KmsGetKeyV1Response{}, errors.NewRequestError(callKmsGetKeyByNameOperationV1, err)
	}

	if res.IsError() {
		return KmsGetKeyV1Response{}, errors.NewAPIErrorWithResponse(callKmsGetKeyByNameOperationV1, res)
	}

	return kmsGetKeyByNameResponse, nil
}
