package api

import (
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/infisical/go-sdk/packages/errors"
)

const callListSecretsV3RawOperation = "CallListSecretsV3Raw"
const callListSecretsWithETagV3RawOperation = "CallListSecretsWithETagV3Raw"

func CallListSecretsWithETagV3(httpClient *resty.Client, request ListSecretsV3RawWithETagRequest) (response ListSecretsV3RawResponse, serverETag string, isModified bool, err error) {

	secretsResponse := ListSecretsV3RawResponse{}

	if request.SecretPath == "" {
		request.SecretPath = "/"
	}

	if request.CurrentETag != "" {

		isWeakETag := strings.HasPrefix(request.CurrentETag, "W/")
		if isWeakETag {
			request.CurrentETag = strings.TrimPrefix(request.CurrentETag, "W/")
		}

		request.CurrentETag = strings.TrimPrefix(request.CurrentETag, "\"")
		request.CurrentETag = strings.TrimSuffix(request.CurrentETag, "\"")
		request.CurrentETag = fmt.Sprintf("\"%s\"", request.CurrentETag)

		if isWeakETag {
			request.CurrentETag = fmt.Sprintf("W/%s", request.CurrentETag)
		}
	}

	fmt.Printf("ETAG: %s\n", request.CurrentETag)

	res, err := httpClient.R().
		SetResult(&secretsResponse).
		SetHeader("if-none-match", request.CurrentETag).
		SetQueryParams(map[string]string{
			"workspaceId":            request.ProjectID,
			"workspaceSlug":          request.ProjectSlug,
			"environment":            request.Environment,
			"secretPath":             request.SecretPath,
			"expandSecretReferences": fmt.Sprintf("%t", request.ExpandSecretReferences),
			"include_imports":        fmt.Sprintf("%t", request.IncludeImports),
			"recursive":              fmt.Sprintf("%t", request.Recursive),
		}).Get("/v3/secrets/raw")

	if err != nil {
		return ListSecretsV3RawResponse{}, "", false, errors.NewRequestError(callListSecretsWithETagV3RawOperation, err)
	}

	if res.IsError() {
		return ListSecretsV3RawResponse{}, "", false, errors.NewAPIErrorWithResponse(callListSecretsWithETagV3RawOperation, res)
	}

	var modified = true
	if res.StatusCode() == 304 || (res.Header().Get("etag") == request.CurrentETag && request.CurrentETag != "") {
		modified = false
	}

	return secretsResponse, res.Header().Get("etag"), modified, nil
}

func CallListSecretsV3(httpClient *resty.Client, request ListSecretsV3RawRequest) (ListSecretsV3RawResponse, error) {
	secretsResponse := ListSecretsV3RawResponse{}

	if request.SecretPath == "" {
		request.SecretPath = "/"
	}

	res, err := httpClient.R().
		SetResult(&secretsResponse).
		SetQueryParams(map[string]string{
			"workspaceId":            request.ProjectID,
			"workspaceSlug":          request.ProjectSlug,
			"environment":            request.Environment,
			"secretPath":             request.SecretPath,
			"expandSecretReferences": fmt.Sprintf("%t", request.ExpandSecretReferences),
			"include_imports":        fmt.Sprintf("%t", request.IncludeImports),
			"recursive":              fmt.Sprintf("%t", request.Recursive),
		}).Get("/v3/secrets/raw")

	if err != nil {
		return ListSecretsV3RawResponse{}, errors.NewRequestError(callListSecretsV3RawOperation, err)
	}

	if res.IsError() {
		return ListSecretsV3RawResponse{}, errors.NewAPIErrorWithResponse(callListSecretsV3RawOperation, res)
	}

	return secretsResponse, nil
}
