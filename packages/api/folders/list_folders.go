package api

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/infisical/go-sdk/packages/util"
)

func CallListFoldersV1(httpClient *resty.Client, request ListFoldersV1Request) (ListFoldersV1Response, error) {

	secretsResponse := ListFoldersV1Response{}

	queryParams := map[string]string{
		"workspaceId": request.ProjectID,
		"environment": request.Environment,
	}

	if request.Path != "" {
		queryParams["path"] = request.Path
	}

	res, err := httpClient.R().
		SetResult(&secretsResponse).
		SetQueryParams(queryParams).Get("/v1/folders")

	if err != nil {
		return ListFoldersV1Response{}, fmt.Errorf("CallListFolderV1: Unable to complete api request [err=%s]", err)
	}

	if res.IsError() {
		return ListFoldersV1Response{}, fmt.Errorf("CallListFolderV1: Unsuccessful response [%v %v] [status-code=%v] %s", res.Request.Method, res.Request.URL, res.StatusCode(), fmt.Sprintf("Error: %s", util.TryParseErrorBody(res)))
	}

	return secretsResponse, nil
}
