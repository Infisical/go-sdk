package api

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/infisical/go-sdk/packages/errors"
)

const callIssueSshCredsFromHostOperation = "CallIssueSshCredsFromHostV1"

func CallIssueSshCredsFromHostV1(httpClient *resty.Client, sshHostId string, body IssueSshCredsFromHostV1Request) (IssueSshCredsFromHostV1Response, error) {
	resp := IssueSshCredsFromHostV1Response{}

	res, err := httpClient.R().
		SetBody(body).
		SetResult(&resp).
		Post(fmt.Sprintf("/v1/ssh/hosts/%s/issue", sshHostId))

	if err != nil {
		return IssueSshCredsFromHostV1Response{}, errors.NewRequestError(callIssueSshCredsFromHostOperation, err)
	}

	if res.IsError() {
		return IssueSshCredsFromHostV1Response{}, errors.NewAPIErrorWithResponse(callIssueSshCredsFromHostOperation, res)
	}

	return resp, nil
}
