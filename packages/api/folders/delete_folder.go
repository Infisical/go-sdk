package api

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/infisical/go-sdk/packages/util"
)

func CallDeleteFolderV1(httpClient *resty.Client, request DeleteFolderV1Request) (DeleteFolderV1Response, error) {

	deleteResponse := DeleteFolderV1Response{}

	req := httpClient.R().
		SetResult(&deleteResponse).
		SetBody(request)

	res, err := req.Delete(fmt.Sprintf("/v1/folders//%s", request.folderID))

	if err != nil {
		return DeleteFolderV1Response{}, fmt.Errorf("CallDeleteFolderV1: Unable to complete api request [err=%s]", err)
	}

	if res.IsError() {
		return DeleteFolderV1Response{}, fmt.Errorf("CallDeleteFolderV1: Unsuccessful response [%v %v] [status-code=%v] %s", res.Request.Method, res.Request.URL, res.StatusCode(), fmt.Sprintf("Error: %s", util.TryParseErrorBody(res)))
	}

	return deleteResponse, nil
}
