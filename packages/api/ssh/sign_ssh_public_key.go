package api

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/infisical/go-sdk/packages/errors"
)

const callSignSshPublicKeyOperation = "CallSignSshPublicKeyV1"

func CallSignSshPublicKeyV1(httpClient *resty.Client, request SignSshPublicKeyV1Request) (SignSshPublicKeyV1Response, error) {
	signSshPublicKeyResponse := SignSshPublicKeyV1Response{}

	res, err := httpClient.R().
		SetResult(&signSshPublicKeyResponse).
		SetBody(request).
		Post(fmt.Sprintf("/v1/ssh/sign"))

	if err != nil {
		return SignSshPublicKeyV1Response{}, errors.NewRequestError(callIssueSshCredsOperation, err)
	}

	if res.IsError() {
		return SignSshPublicKeyV1Response{}, errors.NewAPIErrorWithResponse(callIssueSshCredsOperation, res)
	}

	return signSshPublicKeyResponse, nil
}