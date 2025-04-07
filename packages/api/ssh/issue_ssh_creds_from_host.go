package api

import (
	"github.com/go-resty/resty/v2"
	"github.com/infisical/go-sdk/packages/errors"
)

const callIssueSshCredsFromHostOperation = "CallIssueSshCredsFromHostV1"

func CallIssueSshCredsFromHostV1(httpClient *resty.Client, sshHostId string) (IssueSshCredsFromHostV1Response, error) {
	var result IssueSshCredsFromHostV1Response

	res, err := httpClient.R().
		SetResult(&result).
		Post("/v1/ssh/hosts/" + sshHostId + "/issue")

	if err != nil {
		return IssueSshCredsFromHostV1Response{}, errors.NewRequestError(callIssueSshCredsFromHostOperation, err)
	}

	if res.IsError() {
		return IssueSshCredsFromHostV1Response{}, errors.NewAPIErrorWithResponse(callIssueSshCredsFromHostOperation, res)
	}

	return result, nil
}