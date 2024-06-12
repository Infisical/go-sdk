package api

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/infisical/go-sdk/packages/util"
)

func CallUpdateFolderV1(httpClient *resty.Client, request UpdateFolderV1Request) (UpdateFolderV1Response, error) {

	updateResponse := UpdateFolderV1Response{}

	req := httpClient.R().
		SetResult(&updateResponse).
		SetBody(request)

	res, err := req.Patch(fmt.Sprintf("/v1/folders/%s", request.FolderID))

	if err != nil {
		return UpdateFolderV1Response{}, fmt.Errorf("CallUpdateFolderV1: Unable to complete api request [err=%s]", err)
	}

	if res.IsError() {
		return UpdateFolderV1Response{}, fmt.Errorf("CallUpdateFolderV1: Unsuccessful response [%v %v] [status-code=%v] %s", res.Request.Method, res.Request.URL, res.StatusCode(), fmt.Sprintf("Error: %s", util.TryParseErrorBody(res)))
	}

	return updateResponse, nil
}
