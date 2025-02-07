package api

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

func CallListKmsKeysV1(client *resty.Client, request KmsListKeysV1Request) (*ListKmsKeysV1Response, error) {
	var result ListKmsKeysV1Response
	req := client.R().SetResult(&result)
	req.SetQueryParam("projectId", request.ProjectID)

	if request.Offset != 0 {
		req.SetQueryParam("offset", fmt.Sprint(request.Offset))
	}
	if request.Limit != 0 {
		req.SetQueryParam("limit", fmt.Sprint(request.Limit))
	}
	if request.OrderBy != "" {
		req.SetQueryParam("orderBy", request.OrderBy)
	}
	if request.OrderDir != "" {
		req.SetQueryParam("orderDir", request.OrderDir)
	}
	if request.Search != "" {
		req.SetQueryParam("search", request.Search)
	}

	resp, err := req.Get("/api/v1/kms/keys")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error listing KMS keys: %s", resp.String())
	}

	return &result, nil
}

func CallCreateKmsKeyV1(client *resty.Client, request KmsCreateKeyV1Request) (*CreateKmsKeyV1Response, error) {
	var result CreateKmsKeyV1Response
	resp, err := client.R().
		SetBody(request).
		SetResult(&result).
		Post("/api/v1/kms/keys")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error creating KMS key: %s", resp.String())
	}

	return &result, nil
}

func CallUpdateKmsKeyV1(client *resty.Client, request KmsUpdateKeyV1Request) (*UpdateKmsKeyV1Response, error) {
	var result UpdateKmsKeyV1Response
	resp, err := client.R().
		SetBody(request).
		SetResult(&result).
		Patch(fmt.Sprintf("/api/v1/kms/keys/%s", request.ID))

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error updating KMS key: %s", resp.String())
	}

	return &result, nil
}

func CallDeleteKmsKeyV1(client *resty.Client, request KmsDeleteKeyV1Request) (*DeleteKmsKeyV1Response, error) {
	var result DeleteKmsKeyV1Response
	resp, err := client.R().
		SetResult(&result).
		Delete(fmt.Sprintf("/api/v1/kms/keys/%s", request.ID))

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error deleting KMS key: %s", resp.String())
	}

	return &result, nil
}
