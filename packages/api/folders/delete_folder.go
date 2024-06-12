package api

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/infisical/go-sdk/packages/util"
)

func CallDeleteFolderV1(httpClient *resty.Client, request DeleteFolderV1Request) (DeleteFolderV1Response, error) {

	deleteResponse := DeleteFolderV1Response{}

	// Either folderID or folderName must be provided
	var folderIdOrName string
	if request.folderID != "" {
		folderIdOrName = request.folderID
	} else if request.folderName != "" {
		folderIdOrName = request.folderName
	} else {
		return DeleteFolderV1Response{}, fmt.Errorf("CallDeleteFolderV1: Either folderID or folderName must be provided")
	}

	req := httpClient.R().
		SetResult(&deleteResponse).
		SetBody(request)

	res, err := req.Delete(fmt.Sprintf("/v1/folders/%s", folderIdOrName))

	if err != nil {
		return DeleteFolderV1Response{}, fmt.Errorf("CallDeleteFolderV1: Unable to complete api request [err=%s]", err)
	}

	if res.IsError() {
		return DeleteFolderV1Response{}, fmt.Errorf("CallDeleteFolderV1: Unsuccessful response [%v %v] [status-code=%v] %s", res.Request.Method, res.Request.URL, res.StatusCode(), fmt.Sprintf("Error: %s", util.TryParseErrorBody(res)))
	}

	return deleteResponse, nil
}
