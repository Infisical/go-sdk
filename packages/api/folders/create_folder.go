package api

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/infisical/go-sdk/packages/util"
)

func CallCreateFolderV1(httpClient *resty.Client, request CreateFolderV1Request) (CreateFolderV1Response, error) {

	createResponse := CreateFolderV1Response{}

	req := httpClient.R().
		SetResult(&createResponse).
		SetBody(request)

	res, err := req.Post("/v1/folders")

	if err != nil {
		return CreateFolderV1Response{}, fmt.Errorf("CallCreateFolderV1: Unable to complete api request [err=%s]", err)
	}

	if res.IsError() {
		return CreateFolderV1Response{}, fmt.Errorf("CallCreateFolderV1: Unsuccessful response [%v %v] [status-code=%v] %s", res.Request.Method, res.Request.URL, res.StatusCode(), fmt.Sprintf("Error: %s", util.TryParseErrorBody(res)))
	}

	return createResponse, nil
}
