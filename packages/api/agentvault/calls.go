package api

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/infisical/go-sdk/packages/errors"
)

const basePath = "/v1/agent-vault"

// Creates and token rotations must be sent with a client that doesn't retry once the request may have reached the
// server: a retry could create a second proxy, session or access bundle, or rotate a token twice. Grants and revokes
// must too: they report which grantees changed, so after a lost response, a retry would report each one as skipped.
// Deletes as well: after a lost response, a retry would fail with a 404 for a delete that worked. The caller picks the
// client; each function below says when it needs that one.

// segment encodes an ID as one path segment. The URL would resolve "." and ".." like directories, so a revoke with
// ".." would reach the route that removes members from Agent Vault, which takes the same body. An empty ID would reach
// a route that isn't the one meant either, so all three are refused before any request.
func segment(id string) (string, error) {
	if id == "" || id == "." || id == ".." {
		return "", fmt.Errorf(`invalid ID %q: an ID can't be empty, "." or ".."`, id)
	}
	return url.PathEscape(id), nil
}

// call sends a request and decodes the response into T, turning a failure into the SDK's request or API error.
func call[T any](operation string, req *resty.Request, method string, path string) (T, error) {
	var result T
	res, err := req.SetResult(&result).Execute(method, path)
	if err != nil {
		var zero T
		return zero, errors.NewRequestError(operation, err)
	}
	if res.IsError() {
		var zero T
		return zero, errors.NewAPIErrorWithResponse(operation, res)
	}
	return result, nil
}

// callOnPath is call for a route under basePath that ends in one ID, such as "/proxies/{id}/revoke".
func callOnPath[T any](operation string, req *resty.Request, method string, format string, id string) (T, error) {
	seg, err := segment(id)
	if err != nil {
		var zero T
		return zero, err
	}
	return call[T](operation, req, method, basePath+fmt.Sprintf(format, seg))
}

// query builds a query string, leaving out empty values so the server applies its defaults.
type query url.Values

func (q query) set(key string, value string) {
	if value != "" {
		url.Values(q).Set(key, value)
	}
}

func (q query) setInt(key string, value int) {
	url.Values(q).Set(key, strconv.Itoa(value))
}

// setBool sends a boolean as "true" or "false", the only forms the API reads, and leaves it out when nil.
func (q query) setBool(key string, value *bool) {
	if value != nil {
		url.Values(q).Set(key, strconv.FormatBool(*value))
	}
}

// setList sends a list as one comma-separated value, the way the API takes a list in a query parameter.
func (q query) setList(key string, values []string) {
	q.set(key, strings.Join(values, ","))
}

// Proxies

func CallListProxiesV1(httpClient *resty.Client, request ListProxiesRequest) (ListProxiesResponse, error) {
	q := query{}
	q.set("projectId", request.ProjectID)
	q.set("search", request.Search)
	q.set("orderBy", request.OrderBy)
	q.set("orderDirection", string(request.OrderDirection))
	q.setInt("limit", request.Limit)
	q.setInt("offset", request.Offset)

	return call[ListProxiesResponse]("CallAgentVaultListProxiesV1",
		httpClient.R().SetQueryParamsFromValues(url.Values(q)), http.MethodGet, basePath+"/proxies")
}

// CallCreateProxyV1 needs a client that doesn't retry once the request may have reached the server.
func CallCreateProxyV1(httpClient *resty.Client, request CreateProxyRequest) (ProxyWithToken, error) {
	return call[ProxyWithToken]("CallAgentVaultCreateProxyV1",
		httpClient.R().SetBody(request), http.MethodPost, basePath+"/proxies")
}

func CallUpdateProxyV1(httpClient *resty.Client, proxyID string, request UpdateProxyRequest) (SingleProxyResponse, error) {
	return callOnPath[SingleProxyResponse]("CallAgentVaultUpdateProxyV1",
		httpClient.R().SetBody(request), http.MethodPatch, "/proxies/%s", proxyID)
}

// CallDeleteProxyV1 needs a client that doesn't retry once the request may have reached the server.
func CallDeleteProxyV1(httpClient *resty.Client, proxyID string) (DeleteProxyResponse, error) {
	return callOnPath[DeleteProxyResponse]("CallAgentVaultDeleteProxyV1",
		httpClient.R(), http.MethodDelete, "/proxies/%s", proxyID)
}

// CallRotateProxyTokenV1 needs a client that doesn't retry once the request may have reached the server.
func CallRotateProxyTokenV1(httpClient *resty.Client, proxyID string) (ProxyWithToken, error) {
	return callOnPath[ProxyWithToken]("CallAgentVaultRotateProxyTokenV1",
		httpClient.R(), http.MethodPost, "/proxies/%s/token-auth/proxy-token", proxyID)
}

func CallRevokeProxyV1(httpClient *resty.Client, proxyID string) (SingleProxyResponse, error) {
	return callOnPath[SingleProxyResponse]("CallAgentVaultRevokeProxyV1",
		httpClient.R(), http.MethodPost, "/proxies/%s/revoke", proxyID)
}

// Sessions

func CallListSessionsV1(httpClient *resty.Client, request ListSessionsRequest) (ListSessionsResponse, error) {
	q := query{}
	q.set("projectId", request.ProjectID)
	q.set("scope", string(request.Scope))
	statuses := make([]string, len(request.Statuses))
	for i, status := range request.Statuses {
		statuses[i] = string(status)
	}
	q.setList("status", statuses)
	q.set("search", request.Search)
	// The API takes each metadata filter as its own "key=value" parameter, repeated rather than comma-separated.
	// Sorted, so the same filters always make the same request.
	keys := make([]string, 0, len(request.Metadata))
	for key := range request.Metadata {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		url.Values(q).Add("metadata", key+"="+request.Metadata[key])
	}
	q.setInt("limit", request.Limit)
	q.setInt("offset", request.Offset)

	return call[ListSessionsResponse]("CallAgentVaultListSessionsV1",
		httpClient.R().SetQueryParamsFromValues(url.Values(q)), http.MethodGet, basePath+"/sessions")
}

func CallGetSessionV1(httpClient *resty.Client, sessionID string) (SingleSessionResponse, error) {
	return callOnPath[SingleSessionResponse]("CallAgentVaultGetSessionV1",
		httpClient.R(), http.MethodGet, "/sessions/%s", sessionID)
}

// CallCreateSessionV1 needs a client that doesn't retry once the request may have reached the server.
func CallCreateSessionV1(httpClient *resty.Client, request CreateSessionRequest) (CreateSessionResponse, error) {
	return call[CreateSessionResponse]("CallAgentVaultCreateSessionV1",
		httpClient.R().SetBody(request), http.MethodPost, basePath+"/sessions")
}

func CallUpdateSessionV1(httpClient *resty.Client, sessionID string, request UpdateSessionRequest) (SingleSessionResponse, error) {
	return callOnPath[SingleSessionResponse]("CallAgentVaultUpdateSessionV1",
		httpClient.R().SetBody(request), http.MethodPatch, "/sessions/%s", sessionID)
}

func CallRevokeSessionV1(httpClient *resty.Client, sessionID string) (SingleSessionResponse, error) {
	return callOnPath[SingleSessionResponse]("CallAgentVaultRevokeSessionV1",
		httpClient.R(), http.MethodPost, "/sessions/%s/revoke", sessionID)
}

// Access bundles

func CallListAccessBundlesV1(httpClient *resty.Client, request ListAccessBundlesRequest) (ListAccessBundlesResponse, error) {
	q := query{}
	q.set("projectId", request.ProjectID)
	q.set("search", request.Search)
	q.setBool("requiresSetup", request.RequiresSetup)
	q.setBool("granted", request.Granted)
	q.set("onBehalfOfUserEmail", request.OnBehalfOfUserEmail)
	q.set("onBehalfOfUserId", request.OnBehalfOfUserID)
	q.set("orderBy", request.OrderBy)
	q.set("orderDirection", string(request.OrderDirection))
	q.setInt("limit", request.Limit)
	q.setInt("offset", request.Offset)

	return call[ListAccessBundlesResponse]("CallAgentVaultListAccessBundlesV1",
		httpClient.R().SetQueryParamsFromValues(url.Values(q)), http.MethodGet, basePath+"/access-bundles")
}

func CallGetAccessBundleV1(httpClient *resty.Client, accessBundleID string) (GetAccessBundleResponse, error) {
	return callOnPath[GetAccessBundleResponse]("CallAgentVaultGetAccessBundleV1",
		httpClient.R(), http.MethodGet, "/access-bundles/%s", accessBundleID)
}

// CallCreateAccessBundleV1 needs a client that doesn't retry once the request may have reached the server.
func CallCreateAccessBundleV1(httpClient *resty.Client, request CreateAccessBundleRequest) (SingleAccessBundleResponse, error) {
	return call[SingleAccessBundleResponse]("CallAgentVaultCreateAccessBundleV1",
		httpClient.R().SetBody(request), http.MethodPost, basePath+"/access-bundles")
}

func CallUpdateAccessBundleV1(httpClient *resty.Client, accessBundleID string, request UpdateAccessBundleRequest) (SingleAccessBundleResponse, error) {
	return callOnPath[SingleAccessBundleResponse]("CallAgentVaultUpdateAccessBundleV1",
		httpClient.R().SetBody(request), http.MethodPatch, "/access-bundles/%s", accessBundleID)
}

// CallDeleteAccessBundleV1 needs a client that doesn't retry once the request may have reached the server.
func CallDeleteAccessBundleV1(httpClient *resty.Client, accessBundleID string) (SingleAccessBundleResponse, error) {
	return callOnPath[SingleAccessBundleResponse]("CallAgentVaultDeleteAccessBundleV1",
		httpClient.R(), http.MethodDelete, "/access-bundles/%s", accessBundleID)
}

func granteesQuery(request ListGranteesRequest) url.Values {
	q := query{}
	q.set("actorType", string(request.ActorType))
	q.set("search", request.Search)
	q.setInt("limit", request.Limit)
	q.setInt("offset", request.Offset)
	return url.Values(q)
}

func CallListAccessBundleGranteesV1(httpClient *resty.Client, accessBundleID string, request ListGranteesRequest) (ListGranteesResponse, error) {
	return callOnPath[ListGranteesResponse]("CallAgentVaultListAccessBundleGranteesV1",
		httpClient.R().SetQueryParamsFromValues(granteesQuery(request)), http.MethodGet, "/access-bundles/%s/members", accessBundleID)
}

func CallListAvailableAccessBundleGranteesV1(httpClient *resty.Client, accessBundleID string, request ListGranteesRequest) (ListMembersResponse, error) {
	return callOnPath[ListMembersResponse]("CallAgentVaultListAvailableAccessBundleGranteesV1",
		httpClient.R().SetQueryParamsFromValues(granteesQuery(request)), http.MethodGet, "/access-bundles/%s/members/available", accessBundleID)
}

// CallGrantAccessBundleV1 needs a client that doesn't retry once the request may have reached the server.
func CallGrantAccessBundleV1(httpClient *resty.Client, accessBundleID string, request ChangeGrantsRequest) (GrantResponse, error) {
	return callOnPath[GrantResponse]("CallAgentVaultGrantAccessBundleV1",
		httpClient.R().SetBody(request), http.MethodPost, "/access-bundles/%s/members", accessBundleID)
}

// CallRevokeAccessBundleV1 needs a client that doesn't retry once the request may have reached the server.
func CallRevokeAccessBundleV1(httpClient *resty.Client, accessBundleID string, request ChangeGrantsRequest) (RevokeGrantsResponse, error) {
	return callOnPath[RevokeGrantsResponse]("CallAgentVaultRevokeAccessBundleV1",
		httpClient.R().SetBody(request), http.MethodDelete, "/access-bundles/%s/members", accessBundleID)
}

// Delegations

// A delegation route names the project in its query string, because a machine identity can be in several.
func projectQuery(projectID string) url.Values {
	q := query{}
	q.set("projectId", projectID)
	return url.Values(q)
}

func CallListDelegationsV1(httpClient *resty.Client, request ListDelegationsRequest) (ListDelegationsResponse, error) {
	q := query{}
	q.set("projectId", request.ProjectID)
	q.setList("machineIdentityIds", request.MachineIdentityIDs)
	q.setInt("limit", request.Limit)
	q.setInt("offset", request.Offset)

	return call[ListDelegationsResponse]("CallAgentVaultListDelegationsV1",
		httpClient.R().SetQueryParamsFromValues(url.Values(q)), http.MethodGet, basePath+"/delegations")
}

func CallGetDelegationV1(httpClient *resty.Client, machineIdentityID string, projectID string) (SingleDelegationResponse, error) {
	return callOnPath[SingleDelegationResponse]("CallAgentVaultGetDelegationV1",
		httpClient.R().SetQueryParamsFromValues(projectQuery(projectID)), http.MethodGet, "/delegations/%s", machineIdentityID)
}

// CallSetDelegationV1 replaces the whole delegation, so it's safe to retry.
func CallSetDelegationV1(httpClient *resty.Client, machineIdentityID string, projectID string, request SetDelegationRequest) (SingleDelegationResponse, error) {
	return callOnPath[SingleDelegationResponse]("CallAgentVaultSetDelegationV1",
		httpClient.R().SetQueryParamsFromValues(projectQuery(projectID)).SetBody(request), http.MethodPut, "/delegations/%s", machineIdentityID)
}

// CallDeleteDelegationV1 needs a client that doesn't retry once the request may have reached the server.
func CallDeleteDelegationV1(httpClient *resty.Client, machineIdentityID string, projectID string) (SingleDelegationResponse, error) {
	return callOnPath[SingleDelegationResponse]("CallAgentVaultDeleteDelegationV1",
		httpClient.R().SetQueryParamsFromValues(projectQuery(projectID)), http.MethodDelete, "/delegations/%s", machineIdentityID)
}

// CallAddDelegationSubjectsV1 skips the users and groups the delegation already names, so it's safe to retry.
func CallAddDelegationSubjectsV1(httpClient *resty.Client, machineIdentityID string, projectID string, request ChangeDelegationSubjectsRequest) (SingleDelegationResponse, error) {
	return callOnPath[SingleDelegationResponse]("CallAgentVaultAddDelegationSubjectsV1",
		httpClient.R().SetQueryParamsFromValues(projectQuery(projectID)).SetBody(request), http.MethodPost, "/delegations/%s/subjects", machineIdentityID)
}

// CallRemoveDelegationSubjectsV1 needs a client that doesn't retry once the request may have reached the server.
func CallRemoveDelegationSubjectsV1(httpClient *resty.Client, machineIdentityID string, projectID string, request ChangeDelegationSubjectsRequest) (RemoveDelegationSubjectsResponse, error) {
	return callOnPath[RemoveDelegationSubjectsResponse]("CallAgentVaultRemoveDelegationSubjectsV1",
		httpClient.R().SetQueryParamsFromValues(projectQuery(projectID)).SetBody(request), http.MethodDelete, "/delegations/%s/subjects", machineIdentityID)
}

func CallListAvailableDelegationSubjectsV1(httpClient *resty.Client, machineIdentityID string, request ListAvailableSubjectsRequest) (ListMembersResponse, error) {
	q := query{}
	q.set("projectId", request.ProjectID)
	q.set("actorType", string(request.ActorType))
	q.set("search", request.Search)
	q.setInt("limit", request.Limit)
	q.setInt("offset", request.Offset)

	return callOnPath[ListMembersResponse]("CallAgentVaultListAvailableDelegationSubjectsV1",
		httpClient.R().SetQueryParamsFromValues(url.Values(q)), http.MethodGet, "/delegations/%s/subjects/available", machineIdentityID)
}
