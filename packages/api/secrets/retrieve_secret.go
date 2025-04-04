package api

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/infisical/go-sdk/packages/errors"
	"github.com/infisical/go-sdk/packages/util"
)

const callRetrieveSecretV3RawOperation = "CallRetrieveSecretV3Raw"

func CallRetrieveSecretV3(cache *expirable.LRU[string, interface{}], httpClient *resty.Client, request RetrieveSecretV3RawRequest) (RetrieveSecretV3RawResponse, error) {
	var cacheKey string

	if cache != nil {
		reqBytes, err := json.Marshal(request)
		if err != nil {
			return RetrieveSecretV3RawResponse{}, err
		}
		cacheKey = util.ComputeCacheKeyFromBytes(reqBytes, callRetrieveSecretV3RawOperation)
		if cached, found := cache.Get(cacheKey); found {
			if response, ok := cached.(RetrieveSecretV3RawResponse); ok {
				return response, nil
			}
			cache.Remove(cacheKey)
		}
	}

	retrieveResponse := RetrieveSecretV3RawResponse{}

	if request.Type == "" {
		request.Type = "shared"
	}

	if request.SecretPath == "" {
		request.SecretPath = "/"
	}

	var version string
	if request.Version > 0 {
		version = fmt.Sprintf("%d", request.Version)
	}

	req := httpClient.R().
		SetResult(&retrieveResponse).
		SetQueryParams(map[string]string{
			"workspaceId":     request.ProjectID,
			"environment":     request.Environment,
			"secretPath":      request.SecretPath,
			"include_imports": fmt.Sprintf("%t", request.IncludeImports),
			"type":            request.Type,
			"version":         version,
		})

	res, err := req.Get(fmt.Sprintf("/v3/secrets/raw/%s", request.SecretKey))

	if err != nil {
		return RetrieveSecretV3RawResponse{}, errors.NewRequestError(callRetrieveSecretV3RawOperation, err)
	}

	if res.IsError() {
		return RetrieveSecretV3RawResponse{}, errors.NewAPIErrorWithResponse(callRetrieveSecretV3RawOperation, res)
	}

	if cache != nil {
		cache.Add(cacheKey, retrieveResponse)
	}

	return retrieveResponse, nil
}
