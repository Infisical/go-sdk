package api

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/infisical/go-sdk/packages/errors"
)

const callKmsGetKeyByIdOperationV1 = "CallKmsGetKeyByIdV1"

func CallKmsGetKeyByIdV1(httpClient *resty.Client, request KmsGetKeyByIdV1Request) (KmsGetKeyV1Response, error) {
	kmsGetKeyByIdResponse := KmsGetKeyV1Response{}

	res, err := httpClient.R().
		SetResult(&kmsGetKeyByIdResponse).
		Get(fmt.Sprintf("/v1/kms/key?projectId=%s&id=%s", request.ProjectId, request.KeyId))

	if err != nil {
		return KmsGetKeyV1Response{}, errors.NewRequestError(callKmsGetKeyByIdOperationV1, err)
	}

	if res.IsError() {
		return KmsGetKeyV1Response{}, errors.NewAPIErrorWithResponse(callKmsGetKeyByIdOperationV1, res)
	}

	return kmsGetKeyByIdResponse, nil
}
