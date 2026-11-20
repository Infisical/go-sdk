package infisical

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	agentvaultapi "github.com/infisical/go-sdk/packages/api/agentvault"
	sdkErrors "github.com/infisical/go-sdk/packages/errors"
)

const testProjectID = "6c9c7a3e-0000-4000-8000-000000000001"

// recordedRequest is a request the test server received.
type recordedRequest struct {
	Method string
	// The path as sent, still escaped.
	Path  string
	Query url.Values
	// The JSON body, or nil without one.
	Body map[string]any
}

// agentVaultServer records the requests it receives and answers each with respond.
type agentVaultServer struct {
	t        *testing.T
	mu       sync.Mutex
	requests []recordedRequest
	respond  func(w http.ResponseWriter, r *http.Request)
	url      string
}

func newAgentVaultServer(t *testing.T, respond func(w http.ResponseWriter, r *http.Request)) *agentVaultServer {
	t.Helper()
	s := &agentVaultServer{t: t, respond: respond}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		var body map[string]any
		if len(raw) > 0 {
			if err := json.Unmarshal(raw, &body); err != nil {
				t.Errorf("request body isn't a JSON object: %s", raw)
			}
		}
		s.mu.Lock()
		s.requests = append(s.requests, recordedRequest{
			Method: r.Method,
			Path:   strings.TrimPrefix(r.URL.EscapedPath(), "/api"),
			Query:  r.URL.Query(),
			Body:   body,
		})
		s.mu.Unlock()
		s.respond(w, r)
	}))
	t.Cleanup(server.Close)
	s.url = server.URL
	return s
}

// jsonResponse answers every request with status and body.
func jsonResponse(status int, body any) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(body)
	}
}

// resetConnection drops the connection after the request has reached the server, without a response.
func resetConnection(w http.ResponseWriter, r *http.Request) {
	conn, _, err := w.(http.Hijacker).Hijack()
	if err != nil {
		panic(err)
	}
	if tcp, ok := conn.(*net.TCPConn); ok {
		_ = tcp.SetLinger(0)
	}
	_ = conn.Close()
}

func (s *agentVaultServer) recorded() []recordedRequest {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]recordedRequest(nil), s.requests...)
}

// only returns the one request the server received.
func (s *agentVaultServer) only() recordedRequest {
	s.t.Helper()
	requests := s.recorded()
	if len(requests) != 1 {
		s.t.Fatalf("want 1 request, got %d: %+v", len(requests), requests)
	}
	return requests[0]
}

func newAgentVaultTestClient(t *testing.T, siteURL string) AgentVaultInterface {
	t.Helper()
	// Retries every failed request quickly, so a test can tell the requests that are retried from the ones that
	// aren't.
	return newAgentVaultTestClientWithBackoff(t, siteURL, &ExponentialBackoffStrategy{BaseDelay: time.Millisecond, MaxDelay: time.Millisecond, MaxRetries: 2})
}

func newAgentVaultTestClientWithBackoff(t *testing.T, siteURL string, backoff *ExponentialBackoffStrategy) AgentVaultInterface {
	t.Helper()
	client := NewInfisicalClient(context.Background(), Config{
		SiteUrl:             siteURL,
		AutoTokenRefresh:    BoolPtr(false),
		SilentMode:          true,
		RetryRequestsConfig: &RetryRequestsConfig{ExponentialBackoff: backoff},
	})
	client.Auth().SetAccessToken("test-token")
	return client.AgentVault()
}

func asJSON(t *testing.T, value any) any {
	t.Helper()
	raw, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	var out any
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatal(err)
	}
	return out
}

func assertJSON(t *testing.T, got any, want any) {
	t.Helper()
	if g, w := asJSON(t, got), asJSON(t, want); !reflect.DeepEqual(g, w) {
		gotRaw, _ := json.Marshal(g)
		wantRaw, _ := json.Marshal(w)
		t.Fatalf("got  %s\nwant %s", gotRaw, wantRaw)
	}
}

var exampleProxy = map[string]any{
	"id": "p1", "name": "edge", "description": nil, "heartbeat": nil, "isHealthy": true,
	"rootCaFingerprint": nil, "rootCaExpiresAt": nil, "trafficPolicy": "any-host", "allowedHosts": nil,
	"pollInterval": 60, "createdAt": "2026-11-19T10:00:00.000Z",
}

func proxyWith(fields map[string]any) map[string]any {
	proxy := map[string]any{}
	for k, v := range exampleProxy {
		proxy[k] = v
	}
	for k, v := range fields {
		proxy[k] = v
	}
	return proxy
}

var exampleSession = map[string]any{
	"id": "s1", "label": "nightly deploy", "status": "active", "expiresAt": nil, "revokedAt": nil,
	"createdAt":  "2026-11-19T10:00:00.000Z",
	"actor":      map[string]any{"type": "machineIdentity", "id": "mi1", "name": "orchestrator"},
	"onBehalfOf": nil, "accessBundle": map[string]any{"id": "ab1", "name": "engineering"},
	"services": []any{}, "metadata": map[string]any{"runId": "8812"},
}

func sessionWith(fields map[string]any) map[string]any {
	session := map[string]any{}
	for k, v := range exampleSession {
		session[k] = v
	}
	for k, v := range fields {
		if v == nil {
			delete(session, k)
			continue
		}
		session[k] = v
	}
	return session
}

var exampleCreatedSession = map[string]any{
	"id": "s1", "projectId": testProjectID, "label": nil, "token": "session-token", "expiresAt": nil,
	"createdAt": "2026-11-19T10:00:00.000Z", "accessBundle": nil, "services": []any{}, "metadata": map[string]any{},
}

var exampleDelegation = map[string]any{
	"machineIdentity": map[string]any{"id": "mi1", "name": "orchestrator"}, "scope": "selected",
	"users": []any{}, "groups": []any{}, "createdAt": "2026-11-19T10:00:00.000Z", "updatedAt": "2026-11-19T10:00:00.000Z",
}

// Proxies

func TestProxiesUpdateSendsNullForAnEmptyAllowedHostsList(t *testing.T) {
	server := newAgentVaultServer(t, jsonResponse(200, map[string]any{"proxy": exampleProxy}))
	if _, err := newAgentVaultTestClient(t, server.url).Proxies().Update("p1", AgentVaultUpdateProxyOptions{AllowedHosts: []string{}}); err != nil {
		t.Fatal(err)
	}
	assertJSON(t, server.only().Body, map[string]any{"allowedHosts": nil})
}

func TestProxiesUpdateRefusesAllowedHostsWithOnlyBlankEntries(t *testing.T) {
	// strings.Split("", ",") returns []string{""}, so an empty setting must not remove every allowed host.
	for _, hosts := range [][]string{{""}, {" ", ""}} {
		server := newAgentVaultServer(t, jsonResponse(200, map[string]any{"proxy": exampleProxy}))
		if _, err := newAgentVaultTestClient(t, server.url).Proxies().Update("p1", AgentVaultUpdateProxyOptions{AllowedHosts: hosts}); err == nil {
			t.Fatalf("AllowedHosts %q was accepted", hosts)
		}
		if n := len(server.recorded()); n != 0 {
			t.Fatalf("sent %d requests", n)
		}
	}
}

func TestProxiesUpdateSkipsBlankAllowedHosts(t *testing.T) {
	server := newAgentVaultServer(t, jsonResponse(200, map[string]any{"proxy": exampleProxy}))
	if _, err := newAgentVaultTestClient(t, server.url).Proxies().Update("p1", AgentVaultUpdateProxyOptions{AllowedHosts: []string{" api.github.com ", "", "*.slack.com"}}); err != nil {
		t.Fatal(err)
	}
	assertJSON(t, server.only().Body, map[string]any{"allowedHosts": "api.github.com,*.slack.com"})
}

func TestProxiesUpdateLeavesOutValuesYouDoNotSet(t *testing.T) {
	server := newAgentVaultServer(t, jsonResponse(200, map[string]any{"proxy": exampleProxy}))
	if _, err := newAgentVaultTestClient(t, server.url).Proxies().Update("p1", AgentVaultUpdateProxyOptions{}); err != nil {
		t.Fatal(err)
	}
	assertJSON(t, server.only().Body, map[string]any{})
}

func TestProxiesUpdateSendsAnEmptyDescriptionToRemoveIt(t *testing.T) {
	server := newAgentVaultServer(t, jsonResponse(200, map[string]any{"proxy": exampleProxy}))
	if _, err := newAgentVaultTestClient(t, server.url).Proxies().Update("p1", AgentVaultUpdateProxyOptions{Description: StringPtr("")}); err != nil {
		t.Fatal(err)
	}
	assertJSON(t, server.only().Body, map[string]any{"description": ""})
}

func TestProxiesCreateSendsDefaultsAndJoinsAllowedHosts(t *testing.T) {
	server := newAgentVaultServer(t, jsonResponse(200, map[string]any{"proxy": exampleProxy, "proxyToken": "proxy-token"}))
	res, err := newAgentVaultTestClient(t, server.url).Proxies().Create(AgentVaultCreateProxyOptions{
		ProjectID: testProjectID, Name: "edge", AllowedHosts: []string{"api.github.com", " *.slack.com", ""}, Description: "Sidecars",
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.ProxyToken != "proxy-token" {
		t.Fatalf("got proxy token %q", res.ProxyToken)
	}
	assertJSON(t, server.only().Body, map[string]any{
		"projectId": testProjectID, "name": "edge", "trafficPolicy": "any-host",
		"allowedHosts": "api.github.com,*.slack.com", "pollInterval": 60, "description": "Sidecars",
	})
}

func TestProxiesCreateLeavesOutEmptyAllowedHosts(t *testing.T) {
	for _, hosts := range [][]string{nil, {}, {" "}, {"", ""}} {
		server := newAgentVaultServer(t, jsonResponse(200, map[string]any{"proxy": exampleProxy, "proxyToken": "t"}))
		if _, err := newAgentVaultTestClient(t, server.url).Proxies().Create(AgentVaultCreateProxyOptions{ProjectID: testProjectID, Name: "edge", AllowedHosts: hosts}); err != nil {
			t.Fatal(err)
		}
		if _, ok := server.only().Body["allowedHosts"]; ok {
			t.Fatalf("allowedHosts %q was sent", hosts)
		}
	}
}

func TestProxiesSendServiceHostsByItsOldNameWhichEveryVersionTakes(t *testing.T) {
	server := newAgentVaultServer(t, func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(200, map[string]any{"proxy": exampleProxy, "proxyToken": "t"})(w, r)
	})
	av := newAgentVaultTestClient(t, server.url)
	if _, err := av.Proxies().Create(AgentVaultCreateProxyOptions{ProjectID: testProjectID, Name: "edge", TrafficPolicy: AgentVaultTrafficPolicyServiceHosts}); err != nil {
		t.Fatal(err)
	}
	if _, err := av.Proxies().Update("p1", AgentVaultUpdateProxyOptions{TrafficPolicy: AgentVaultTrafficPolicyServiceHosts}); err != nil {
		t.Fatal(err)
	}
	if _, err := av.Proxies().Update("p1", AgentVaultUpdateProxyOptions{TrafficPolicy: AgentVaultTrafficPolicyAnyHost}); err != nil {
		t.Fatal(err)
	}
	requests := server.recorded()
	for i, want := range []string{"bundle-hosts", "bundle-hosts", "any-host"} {
		if got := requests[i].Body["trafficPolicy"]; got != want {
			t.Errorf("request %d sent trafficPolicy %v, want %s", i, got, want)
		}
	}
}

func TestProxiesReportTheOldNameAsServiceHosts(t *testing.T) {
	oldProxy := proxyWith(map[string]any{"trafficPolicy": "bundle-hosts"})
	server := newAgentVaultServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			jsonResponse(200, map[string]any{"proxies": []any{oldProxy}, "totalCount": 1})(w, r)
			return
		}
		jsonResponse(200, map[string]any{"proxy": oldProxy, "proxyToken": "t"})(w, r)
	})
	proxies := newAgentVaultTestClient(t, server.url).Proxies()

	list, err := proxies.List(AgentVaultListProxiesOptions{ProjectID: testProjectID})
	if err != nil {
		t.Fatal(err)
	}
	created, err := proxies.Create(AgentVaultCreateProxyOptions{ProjectID: testProjectID, Name: "edge"})
	if err != nil {
		t.Fatal(err)
	}
	updated, err := proxies.Update("p1", AgentVaultUpdateProxyOptions{})
	if err != nil {
		t.Fatal(err)
	}
	rotated, err := proxies.RotateToken("p1")
	if err != nil {
		t.Fatal(err)
	}
	revoked, err := proxies.Revoke("p1")
	if err != nil {
		t.Fatal(err)
	}
	for _, got := range []AgentVaultTrafficPolicy{list.Proxies[0].TrafficPolicy, created.Proxy.TrafficPolicy, updated.TrafficPolicy, rotated.Proxy.TrafficPolicy, revoked.TrafficPolicy} {
		if got != AgentVaultTrafficPolicyServiceHosts {
			t.Errorf("got traffic policy %q, want service-hosts", got)
		}
	}
}

func TestProxiesListLeavesTheTrafficPolicyUnsetForAMember(t *testing.T) {
	memberView := map[string]any{"id": "p1", "name": "edge", "description": nil, "heartbeat": nil, "isHealthy": true, "rootCaFingerprint": nil, "rootCaExpiresAt": nil}
	server := newAgentVaultServer(t, jsonResponse(200, map[string]any{"proxies": []any{memberView}, "totalCount": 1}))
	res, err := newAgentVaultTestClient(t, server.url).Proxies().List(AgentVaultListProxiesOptions{ProjectID: testProjectID})
	if err != nil {
		t.Fatal(err)
	}
	if res.Proxies[0].TrafficPolicy != "" || res.Proxies[0].CreatedAt != nil {
		t.Fatalf("got settings on a member's view: %+v", res.Proxies[0])
	}
}

// Sessions

func TestSessionsCreateSendsOnlyWhatYouSet(t *testing.T) {
	server := newAgentVaultServer(t, jsonResponse(200, map[string]any{"session": exampleCreatedSession}))
	res, err := newAgentVaultTestClient(t, server.url).Sessions().Create(AgentVaultCreateSessionOptions{
		ProjectID: testProjectID, AccessBundle: "engineering",
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.Token != "session-token" || res.ProjectID != testProjectID {
		t.Fatalf("got %+v", res)
	}
	// TTL is left out, so the server picks the default: 7d, or 1h on behalf of a user.
	assertJSON(t, server.only().Body, map[string]any{"projectId": testProjectID, "accessBundle": "engineering"})
}

func TestSessionsCreateSendsMetadataAndReadsItBack(t *testing.T) {
	created := map[string]any{}
	for k, v := range exampleCreatedSession {
		created[k] = v
	}
	created["metadata"] = map[string]any{"runId": "8812", "channel": "C123"}
	server := newAgentVaultServer(t, jsonResponse(200, map[string]any{"session": created}))
	res, err := newAgentVaultTestClient(t, server.url).Sessions().Create(AgentVaultCreateSessionOptions{
		ProjectID: testProjectID, AccessBundle: "engineering", Metadata: map[string]string{"runId": "8812", "channel": "C123"},
	})
	if err != nil {
		t.Fatal(err)
	}
	assertJSON(t, server.only().Body["metadata"], map[string]any{"runId": "8812", "channel": "C123"})
	assertJSON(t, res.Metadata, map[string]any{"runId": "8812", "channel": "C123"})
}

func TestSessionsCreateNamesTheUserItActsFor(t *testing.T) {
	cases := []struct {
		onBehalfOf AgentVaultOnBehalfOf
		want       map[string]any
	}{
		{AgentVaultOnBehalfOf{UserEmail: "jsmith@example.com"}, map[string]any{"userEmail": "jsmith@example.com"}},
		{AgentVaultOnBehalfOf{UserID: "u1"}, map[string]any{"userId": "u1"}},
	}
	for _, c := range cases {
		server := newAgentVaultServer(t, jsonResponse(200, map[string]any{"session": exampleCreatedSession}))
		onBehalfOf := c.onBehalfOf
		if _, err := newAgentVaultTestClient(t, server.url).Sessions().Create(AgentVaultCreateSessionOptions{ProjectID: testProjectID, AccessBundle: "sales", OnBehalfOf: &onBehalfOf}); err != nil {
			t.Fatal(err)
		}
		assertJSON(t, server.only().Body["onBehalfOf"], c.want)
	}
}

func TestSessionsCreateRefusesAUserNamedTwiceOrNotAtAll(t *testing.T) {
	for _, onBehalfOf := range []AgentVaultOnBehalfOf{{UserEmail: "jsmith@example.com", UserID: "u1"}, {}} {
		server := newAgentVaultServer(t, jsonResponse(200, map[string]any{"session": exampleCreatedSession}))
		onBehalfOf := onBehalfOf
		if _, err := newAgentVaultTestClient(t, server.url).Sessions().Create(AgentVaultCreateSessionOptions{ProjectID: testProjectID, AccessBundle: "sales", OnBehalfOf: &onBehalfOf}); err == nil {
			t.Fatalf("OnBehalfOf %+v was accepted", onBehalfOf)
		}
		if n := len(server.recorded()); n != 0 {
			t.Fatalf("sent %d requests", n)
		}
	}
}

func TestSessionsCreateSendsInlineServicesWithTheirTypes(t *testing.T) {
	server := newAgentVaultServer(t, jsonResponse(200, map[string]any{"session": exampleCreatedSession}))
	_, err := newAgentVaultTestClient(t, server.url).Sessions().Create(AgentVaultCreateSessionOptions{
		ProjectID: testProjectID,
		Services: []AgentVaultInlineService{
			AgentVaultInlineHTTPService{
				HostPattern:    "api.github.com",
				AllowedMethods: []string{"GET"},
				Credential:     AgentVaultBearerCredential{Value: "ghp_x"},
			},
			AgentVaultInlineHTTPService{
				Name:        "raw-token",
				HostPattern: "api.example.com",
				Credential:  AgentVaultBearerCredential{Value: "x", HeaderPrefix: StringPtr("")},
				Substitutions: []AgentVaultInlineSubstitution{
					{Placeholder: "__ORG__", Surfaces: []AgentVaultSubstitutionSurface{AgentVaultSubstitutionSurfacePath}, Value: "acme"},
				},
			},
			AgentVaultInlineHTTPService{HostPattern: "pypi.org", Credential: AgentVaultPassthroughCredential{}},
			AgentVaultInlineMCPService{
				ServerURL:    "https://mcp.linear.app/mcp",
				Transport:    AgentVaultMCPTransportStreamableHTTP,
				Credential:   AgentVaultBearerCredential{Value: "lin_x"},
				AllowedTools: []string{"list_issues"},
			},
			AgentVaultInlinePostgresService{
				Name: "analytics", Host: "db.internal", Database: "analytics",
				TLS:        &AgentVaultInlineTLS{Verify: BoolPtr(false)},
				Credential: AgentVaultLoginCredential{Username: "readonly", Password: "pw"},
			},
			AgentVaultInlineMySQLService{
				Name: "shop", Host: "mysql.internal", Database: "shop",
				Credential: AgentVaultLoginCredential{Username: "root"},
			},
			AgentVaultInlineSSHService{
				Name: "bastion", Host: "10.0.0.5", Port: 2222,
				Credential: AgentVaultSSHKeyCredential{Username: "ops", PrivateKey: "-----BEGIN OPENSSH PRIVATE KEY-----"},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	body := server.only().Body
	if _, ok := body["accessBundle"]; ok {
		t.Fatal("accessBundle was sent for a session without one")
	}
	assertJSON(t, body["services"], []any{
		map[string]any{"type": "http", "hostPattern": "api.github.com", "allowedMethods": []any{"GET"}, "credential": map[string]any{"type": "bearer", "value": "ghp_x"}},
		map[string]any{"type": "http", "name": "raw-token", "hostPattern": "api.example.com",
			"credential":    map[string]any{"type": "bearer", "value": "x", "headerPrefix": ""},
			"substitutions": []any{map[string]any{"placeholder": "__ORG__", "surfaces": []any{"path"}, "value": "acme"}}},
		map[string]any{"type": "http", "hostPattern": "pypi.org", "credential": map[string]any{"type": "passthrough"}},
		map[string]any{"type": "mcp", "serverUrl": "https://mcp.linear.app/mcp", "transport": "streamable-http",
			"credential": map[string]any{"type": "bearer", "value": "lin_x"}, "allowedTools": []any{"list_issues"}},
		map[string]any{"type": "postgres", "name": "analytics", "host": "db.internal", "database": "analytics",
			"tls": map[string]any{"verify": false}, "credential": map[string]any{"type": "basic", "username": "readonly", "password": "pw"}},
		map[string]any{"type": "mysql", "name": "shop", "host": "mysql.internal", "database": "shop",
			"credential": map[string]any{"type": "basic", "username": "root"}},
		map[string]any{"type": "ssh", "name": "bastion", "host": "10.0.0.5", "port": 2222,
			"credential": map[string]any{"type": "ssh-key", "username": "ops", "privateKey": "-----BEGIN OPENSSH PRIVATE KEY-----"}},
	})
}

func TestInlineMCPServicesDefaultToStreamableHTTP(t *testing.T) {
	assertJSON(t, AgentVaultInlineMCPService{ServerURL: "https://mcp.example.com/mcp", Credential: AgentVaultBearerCredential{Value: "x"}},
		map[string]any{"type": "mcp", "serverUrl": "https://mcp.example.com/mcp", "transport": "streamable-http", "credential": map[string]any{"type": "bearer", "value": "x"}})
}

func TestLoginCredentialsAlwaysSendTheUsername(t *testing.T) {
	// The API needs a username to log in to a database or SSH server, so an empty one is sent for it to refuse
	// by name rather than dropped.
	assertJSON(t, AgentVaultLoginCredential{Password: "pw"}, map[string]any{"type": "basic", "username": "", "password": "pw"})
	assertJSON(t, AgentVaultInlineSSHService{Name: "bastion", Host: "10.0.0.5", Credential: AgentVaultLoginCredential{Username: "ops", Password: "pw"}},
		map[string]any{"type": "ssh", "name": "bastion", "host": "10.0.0.5", "credential": map[string]any{"type": "basic", "username": "ops", "password": "pw"}})
}

func TestInlineServicesSendAnEmptyAllowListRatherThanAllowingEverything(t *testing.T) {
	assertJSON(t, AgentVaultInlineMCPService{ServerURL: "https://mcp.example.com/mcp", Transport: AgentVaultMCPTransportSSE, Credential: AgentVaultBearerCredential{Value: "x"}, AllowedTools: []string{}},
		map[string]any{"type": "mcp", "serverUrl": "https://mcp.example.com/mcp", "transport": "sse", "credential": map[string]any{"type": "bearer", "value": "x"}, "allowedTools": []any{}})
	assertJSON(t, AgentVaultInlineHTTPService{HostPattern: "a.com", Credential: AgentVaultPassthroughCredential{}, AllowedMethods: []string{}, AllowedPathPrefixes: []string{}},
		map[string]any{"type": "http", "hostPattern": "a.com", "credential": map[string]any{"type": "passthrough"}, "allowedMethods": []any{}, "allowedPathPrefixes": []any{}})
}

func TestSessionsReadServicesAsTheAPISendsThem(t *testing.T) {
	services := []any{
		map[string]any{"id": "sv1", "name": "api.github.com", "description": nil, "type": "http", "hostPattern": "api.github.com",
			"allowedMethods": nil, "allowedPathPrefixes": []any{"/repos"}, "credential": map[string]any{"type": "bearer", "headerName": "Authorization", "headerPrefix": "Bearer"},
			"customHeaders": []any{}, "substitutions": []any{}},
		map[string]any{"id": "sv2", "name": "linear", "type": "mcp", "serverUrl": "https://mcp.linear.app/mcp", "transport": "streamable-http",
			"hostPattern": "mcp.linear.app", "allowedTools": []any{"list_issues"}, "credential": map[string]any{"type": "bearer", "headerName": "Authorization", "headerPrefix": "Bearer"}, "customHeaders": []any{}},
		map[string]any{"id": "sv3", "name": "analytics", "type": "postgres", "host": "db.internal", "port": 5432, "database": "analytics",
			"tls": map[string]any{"enabled": true, "verify": true, "ca": nil}, "credential": map[string]any{"type": "basic"}},
		map[string]any{"id": "sv4", "name": "bastion", "type": "ssh", "host": "10.0.0.5", "port": 22, "credential": map[string]any{"type": "ssh-key", "hasCertificate": true}},
	}
	server := newAgentVaultServer(t, jsonResponse(200, map[string]any{"session": sessionWith(map[string]any{"services": services})}))
	res, err := newAgentVaultTestClient(t, server.url).Sessions().Get("s1")
	if err != nil {
		t.Fatal(err)
	}
	got := res.Services
	if got[0].Credential.HeaderName != "Authorization" || got[0].AllowedMethods != nil || got[0].AllowedPathPrefixes[0] != "/repos" {
		t.Errorf("http service: %+v", got[0])
	}
	if got[1].Transport != AgentVaultMCPTransportStreamableHTTP || got[1].AllowedTools[0] != "list_issues" {
		t.Errorf("mcp service: %+v", got[1])
	}
	if got[2].Port != 5432 || got[2].TLS == nil || !got[2].TLS.Verify || got[2].HostPattern != "" {
		t.Errorf("database service: %+v", got[2])
	}
	if got[3].Type != AgentVaultServiceTypeSSH || !got[3].Credential.HasCertificate {
		t.Errorf("ssh service: %+v", got[3])
	}
}

func TestSessionsFromAVersionWithoutInlineServicesOrMetadataHaveEmptyOnes(t *testing.T) {
	old := sessionWith(map[string]any{"services": nil, "metadata": nil})
	server := newAgentVaultServer(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/sessions") {
			jsonResponse(200, map[string]any{"sessions": []any{old}, "totalCount": 1})(w, r)
			return
		}
		jsonResponse(200, map[string]any{"session": old})(w, r)
	})
	sessions := newAgentVaultTestClient(t, server.url).Sessions()
	list, err := sessions.List(AgentVaultListSessionsOptions{ProjectID: testProjectID})
	if err != nil {
		t.Fatal(err)
	}
	got, err := sessions.Get("s1")
	if err != nil {
		t.Fatal(err)
	}
	updated, err := sessions.Update("s1", AgentVaultUpdateSessionOptions{Label: StringPtr("x")})
	if err != nil {
		t.Fatal(err)
	}
	for _, session := range []AgentVaultSession{list.Sessions[0], got, updated} {
		if session.Services == nil || len(session.Services) != 0 || session.Metadata == nil || len(session.Metadata) != 0 {
			t.Errorf("got services %#v and metadata %#v", session.Services, session.Metadata)
		}
	}
}

func TestSessionsRevokeLeavesTheAccessBundleAndServicesUnsetWhenTheResponseOmitsThem(t *testing.T) {
	server := newAgentVaultServer(t, jsonResponse(200, map[string]any{"session": sessionWith(map[string]any{"accessBundle": nil, "services": nil, "status": "revoked"})}))
	res, err := newAgentVaultTestClient(t, server.url).Sessions().Revoke("s1")
	if err != nil {
		t.Fatal(err)
	}
	if res.AccessBundle != nil || res.Services != nil || res.Status != AgentVaultSessionStatusRevoked {
		t.Fatalf("got %+v", res)
	}
}

func TestSessionsGetReadsTheUserASessionActsForAndADeletedAccessBundle(t *testing.T) {
	server := newAgentVaultServer(t, jsonResponse(200, map[string]any{"session": sessionWith(map[string]any{
		"onBehalfOf":   map[string]any{"type": "user", "id": "u1", "username": "jsmith", "email": "jsmith@example.com", "firstName": "J", "lastName": "Smith"},
		"accessBundle": map[string]any{"id": nil, "name": "deleted-bundle"},
	})}))
	res, err := newAgentVaultTestClient(t, server.url).Sessions().Get("s1")
	if err != nil {
		t.Fatal(err)
	}
	if res.OnBehalfOf == nil || *res.OnBehalfOf.Email != "jsmith@example.com" {
		t.Fatalf("onBehalfOf: %+v", res.OnBehalfOf)
	}
	if res.AccessBundle.ID != nil || res.AccessBundle.Name != "deleted-bundle" {
		t.Fatalf("accessBundle: %+v", res.AccessBundle)
	}
	if res.Actor.Type != AgentVaultMemberTypeMachineIdentity || res.Actor.Name != "orchestrator" {
		t.Fatalf("actor: %+v", res.Actor)
	}
}

func TestSessionsUpdateSendsOnlyWhatYouSet(t *testing.T) {
	cases := []struct {
		options AgentVaultUpdateSessionOptions
		want    map[string]any
	}{
		{AgentVaultUpdateSessionOptions{Label: StringPtr("rerun")}, map[string]any{"label": "rerun"}},
		{AgentVaultUpdateSessionOptions{Label: StringPtr("")}, map[string]any{"label": ""}},
		{AgentVaultUpdateSessionOptions{Metadata: map[string]string{"runId": "9"}}, map[string]any{"metadata": map[string]any{"runId": "9"}}},
		// An empty map removes the metadata, so it's sent rather than dropped.
		{AgentVaultUpdateSessionOptions{Metadata: map[string]string{}}, map[string]any{"metadata": map[string]any{}}},
		{AgentVaultUpdateSessionOptions{Label: StringPtr("a"), Metadata: map[string]string{"k": "v"}}, map[string]any{"label": "a", "metadata": map[string]any{"k": "v"}}},
	}
	for _, c := range cases {
		server := newAgentVaultServer(t, jsonResponse(200, map[string]any{"session": exampleSession}))
		if _, err := newAgentVaultTestClient(t, server.url).Sessions().Update("s1", c.options); err != nil {
			t.Fatal(err)
		}
		assertJSON(t, server.only().Body, c.want)
	}
}

func TestSessionsListSendsEachMetadataPairAsItsOwnParameter(t *testing.T) {
	server := newAgentVaultServer(t, jsonResponse(200, map[string]any{"sessions": []any{}, "totalCount": 0}))
	_, err := newAgentVaultTestClient(t, server.url).Sessions().List(AgentVaultListSessionsOptions{
		ProjectID: testProjectID,
		Statuses:  []AgentVaultSessionStatus{AgentVaultSessionStatusActive, AgentVaultSessionStatusExpired},
		Metadata:  map[string]string{"runId": "8812", "channel": "C1"},
	})
	if err != nil {
		t.Fatal(err)
	}
	q := server.only().Query
	assertJSON(t, q["metadata"], []string{"channel=C1", "runId=8812"})
	if q.Get("status") != "active,expired" || q.Get("scope") != "mine" || q.Get("limit") != "20" || q.Has("offset") {
		t.Fatalf("query: %v", q)
	}
}

func TestSessionsListLeavesOutAnEmptyMetadataFilter(t *testing.T) {
	for _, metadata := range []map[string]string{nil, {}} {
		server := newAgentVaultServer(t, jsonResponse(200, map[string]any{"sessions": []any{}, "totalCount": 0}))
		if _, err := newAgentVaultTestClient(t, server.url).Sessions().List(AgentVaultListSessionsOptions{ProjectID: testProjectID, Metadata: metadata}); err != nil {
			t.Fatal(err)
		}
		if q := server.only().Query; q.Has("metadata") || q.Has("status") {
			t.Fatalf("query: %v", q)
		}
	}
}

func TestListsLeaveOutALimitAndOffsetYouDoNotSet(t *testing.T) {
	server := newAgentVaultServer(t, jsonResponse(200, map[string]any{"sessions": []any{}, "totalCount": 0}))
	httpClient := resty.New().SetBaseURL(server.url + "/api")
	if _, err := agentvaultapi.CallListSessionsV1(httpClient, agentvaultapi.ListSessionsRequest{ProjectID: testProjectID}); err != nil {
		t.Fatal(err)
	}
	if q := server.only().Query; q.Has("limit") || q.Has("offset") {
		t.Fatalf("sent query %s", q.Encode())
	}
}

func TestSessionsCreateKeepsTheServersExplanationWhenItRefuses(t *testing.T) {
	server := newAgentVaultServer(t, jsonResponse(403, map[string]any{"message": "Only Agent Vault admins can send services with a session", "reqId": "req-1"}))
	_, err := newAgentVaultTestClient(t, server.url).Sessions().Create(AgentVaultCreateSessionOptions{
		ProjectID: testProjectID, Services: []AgentVaultInlineService{AgentVaultInlineHTTPService{HostPattern: "a.com", Credential: AgentVaultPassthroughCredential{}}},
	})
	var apiErr *sdkErrors.APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != 403 || apiErr.ErrorMessage != "Only Agent Vault admins can send services with a session" {
		t.Fatalf("got %v", err)
	}
}

// Access bundles

func TestAccessBundlesListSendsBooleansTheWayTheAPIReadsThem(t *testing.T) {
	server := newAgentVaultServer(t, jsonResponse(200, map[string]any{"accessBundles": []any{}, "totalCount": 0}))
	if _, err := newAgentVaultTestClient(t, server.url).AccessBundles().List(AgentVaultListAccessBundlesOptions{
		ProjectID: testProjectID, RequiresSetup: BoolPtr(true), Granted: BoolPtr(false),
	}); err != nil {
		t.Fatal(err)
	}
	q := server.only().Query
	if q.Get("requiresSetup") != "true" || q.Get("granted") != "false" || q.Get("orderBy") != "createdAt" || q.Get("orderDirection") != "desc" {
		t.Fatalf("query: %v", q)
	}
}

func TestAccessBundlesListNamesTheUserASessionWouldActFor(t *testing.T) {
	server := newAgentVaultServer(t, jsonResponse(200, map[string]any{"accessBundles": []any{}, "totalCount": 0}))
	av := newAgentVaultTestClient(t, server.url)
	if _, err := av.AccessBundles().List(AgentVaultListAccessBundlesOptions{ProjectID: testProjectID, OnBehalfOf: &AgentVaultOnBehalfOf{UserEmail: "jsmith@example.com"}}); err != nil {
		t.Fatal(err)
	}
	if q := server.only().Query; q.Get("onBehalfOfUserEmail") != "jsmith@example.com" || q.Has("onBehalfOfUserId") || q.Has("requiresSetup") {
		t.Fatalf("query: %v", q)
	}
	if _, err := av.AccessBundles().List(AgentVaultListAccessBundlesOptions{ProjectID: testProjectID, OnBehalfOf: &AgentVaultOnBehalfOf{UserEmail: "a@b.c", UserID: "u1"}}); err == nil {
		t.Fatal("a user named twice was accepted")
	}
}

func TestAccessBundlesListReadsTheServicesRequiringSetup(t *testing.T) {
	server := newAgentVaultServer(t, jsonResponse(200, map[string]any{"accessBundles": []any{map[string]any{
		"id": "ab1", "name": "engineering", "description": nil, "createdAt": "2026-11-19T10:00:00.000Z", "updatedAt": "2026-11-19T10:00:00.000Z",
		"serviceCount": 2, "memberCount": 3, "hostPatterns": []any{"api.github.com"}, "servicesRequiringSetup": []any{map[string]any{"id": "sv1", "name": "github"}},
	}}, "totalCount": 1}))
	res, err := newAgentVaultTestClient(t, server.url).AccessBundles().List(AgentVaultListAccessBundlesOptions{ProjectID: testProjectID})
	if err != nil {
		t.Fatal(err)
	}
	bundle := res.AccessBundles[0]
	if bundle.Name != "engineering" || bundle.ServiceCount != 2 || bundle.ServicesRequiringSetup[0].Name != "github" {
		t.Fatalf("got %+v", bundle)
	}
}

func TestAccessBundlesGetReadsServicesAndKeepsTheirNestedSettings(t *testing.T) {
	server := newAgentVaultServer(t, jsonResponse(200, map[string]any{"accessBundle": map[string]any{
		"id": "ab1", "name": "engineering", "description": "Eng", "createdAt": "2026-11-19T10:00:00.000Z", "updatedAt": "2026-11-19T10:00:00.000Z",
		"services": []any{
			map[string]any{"id": "sv1", "name": "github", "type": "http", "hostPattern": "api.github.com", "allowedMethods": nil, "allowedPathPrefixes": nil,
				"credential": map[string]any{"type": "bearer", "headerName": "Authorization"}, "personalValues": map[string]any{"token": true},
				"customHeaders": []any{}, "substitutions": []any{}, "createdAt": "2026-11-19T10:00:00.000Z", "updatedAt": "2026-11-19T10:00:00.000Z"},
			map[string]any{"id": "sv2", "name": "linear", "type": "mcp", "serverUrl": "https://mcp.linear.app/mcp", "transport": "streamable-http",
				"hostPattern": "mcp.linear.app", "allowedTools": []any{"list_issues"}, "credential": map[string]any{"type": "oauth"}, "personalValues": nil,
				"customHeaders": []any{}, "createdAt": "2026-11-19T10:00:00.000Z", "updatedAt": "2026-11-19T10:00:00.000Z"},
			map[string]any{"id": "sv3", "name": "db", "type": "mysql", "host": "db", "port": 3306, "database": "app", "tls": map[string]any{"enabled": true, "verify": false, "ca": nil},
				"credential": map[string]any{"type": "basic"}, "personalValues": nil, "createdAt": "2026-11-19T10:00:00.000Z", "updatedAt": "2026-11-19T10:00:00.000Z"},
		},
	}}))
	res, err := newAgentVaultTestClient(t, server.url).AccessBundles().Get("ab1")
	if err != nil {
		t.Fatal(err)
	}
	if credential := res.Services[0].Credential; credential["headerName"] != "Authorization" {
		t.Fatalf("credential: %v", credential)
	}
	if res.Services[0].PersonalValues["token"] != true || res.Services[1].AllowedTools[0] != "list_issues" {
		t.Fatalf("services: %+v", res.Services)
	}
	if res.Services[2].Port != 3306 || res.Services[2].HostPattern != "" || res.Services[2].TLS.Verify {
		t.Fatalf("database service: %+v", res.Services[2])
	}
}

func TestAccessBundlesUpdateSendsAnEmptyDescriptionToRemoveIt(t *testing.T) {
	server := newAgentVaultServer(t, jsonResponse(200, map[string]any{"accessBundle": map[string]any{"id": "ab1", "name": "eng", "description": nil, "createdAt": "2026-11-19T10:00:00.000Z", "updatedAt": "2026-11-19T10:00:00.000Z"}}))
	if _, err := newAgentVaultTestClient(t, server.url).AccessBundles().Update("ab1", AgentVaultUpdateAccessBundleOptions{Description: StringPtr("")}); err != nil {
		t.Fatal(err)
	}
	assertJSON(t, server.only().Body, map[string]any{"description": ""})
}

func TestAccessBundlesGrantAndRevokeSendOnlyTheGranteesYouName(t *testing.T) {
	server := newAgentVaultServer(t, jsonResponse(200, map[string]any{"members": []any{}, "skipped": []any{}}))
	bundles := newAgentVaultTestClient(t, server.url).AccessBundles()
	if _, err := bundles.Grant("ab1", AgentVaultChangeGrantsOptions{MachineIdentityIDs: []string{"mi1"}}); err != nil {
		t.Fatal(err)
	}
	if _, err := bundles.Revoke("ab1", AgentVaultChangeGrantsOptions{UserIDs: []string{"u1"}}); err != nil {
		t.Fatal(err)
	}
	requests := server.recorded()
	if requests[0].Method != http.MethodPost || requests[1].Method != http.MethodDelete || requests[0].Path != "/v1/agent-vault/access-bundles/ab1/members" {
		t.Fatalf("requests: %+v", requests)
	}
	assertJSON(t, requests[0].Body, map[string]any{"machineIdentityIds": []any{"mi1"}})
	assertJSON(t, requests[1].Body, map[string]any{"userIds": []any{"u1"}})
}

func TestAccessBundlesListGranteesReadsAPendingOrganizationMembership(t *testing.T) {
	server := newAgentVaultServer(t, jsonResponse(200, map[string]any{"members": []any{
		map[string]any{"type": "user", "id": "u1", "username": "jsmith", "email": nil, "firstName": nil, "lastName": nil, "isOrgMembershipPending": true, "grantedAt": "2026-11-19T10:00:00.000Z"},
	}, "totalCount": 1}))
	res, err := newAgentVaultTestClient(t, server.url).AccessBundles().ListGrantees("ab1", AgentVaultListGranteesOptions{ActorType: AgentVaultMemberTypeUser})
	if err != nil {
		t.Fatal(err)
	}
	if !res.Members[0].IsOrgMembershipPending || server.only().Query.Get("actorType") != "user" {
		t.Fatalf("got %+v", res)
	}
}

// Delegations

func TestDelegationsSetPutsTheScopeAndSendsTheProjectInTheQuery(t *testing.T) {
	server := newAgentVaultServer(t, jsonResponse(200, map[string]any{"delegation": exampleDelegation}))
	if _, err := newAgentVaultTestClient(t, server.url).Delegations().Set("mi1", AgentVaultSetDelegationOptions{
		ProjectID: testProjectID, Scope: AgentVaultDelegationScopeSelected, GroupIDs: []string{"g1"},
	}); err != nil {
		t.Fatal(err)
	}
	req := server.only()
	if req.Method != http.MethodPut || req.Query.Get("projectId") != testProjectID {
		t.Fatalf("request: %+v", req)
	}
	assertJSON(t, req.Body, map[string]any{"scope": "selected", "groupIds": []any{"g1"}})
}

func TestDelegationsGetReadsUsersGroupsAndMissingAccessBundles(t *testing.T) {
	delegation := map[string]any{}
	for k, v := range exampleDelegation {
		delegation[k] = v
	}
	delegation["users"] = []any{map[string]any{"id": "u1", "username": "jsmith", "email": nil, "firstName": nil, "lastName": nil, "addedAt": "2026-11-19T10:00:00.000Z",
		"missingAccessBundles": []any{map[string]any{"id": "ab2", "name": "sales"}}}}
	delegation["groups"] = []any{map[string]any{"id": "g1", "name": "eng", "addedAt": "2026-11-19T10:00:00.000Z", "missingAccessBundles": []any{}}}
	delegation["missingAccessBundles"] = []any{map[string]any{"id": "ab2", "name": "sales"}}
	server := newAgentVaultServer(t, jsonResponse(200, map[string]any{"delegation": delegation}))
	res, err := newAgentVaultTestClient(t, server.url).Delegations().Get("mi1", AgentVaultProjectOptions{ProjectID: testProjectID})
	if err != nil {
		t.Fatal(err)
	}
	if res.Users[0].MissingAccessBundles[0].Name != "sales" || res.Groups[0].Name != "eng" || res.MissingAccessBundles[0].ID != "ab2" {
		t.Fatalf("got %+v", res)
	}
}

func TestDelegationsListJoinsTheMachineIdentityIDs(t *testing.T) {
	server := newAgentVaultServer(t, jsonResponse(200, map[string]any{"delegations": []any{}, "totalCount": 0}))
	if _, err := newAgentVaultTestClient(t, server.url).Delegations().List(AgentVaultListDelegationsOptions{ProjectID: testProjectID, MachineIdentityIDs: []string{"mi1", "mi2"}}); err != nil {
		t.Fatal(err)
	}
	if q := server.only().Query; q.Get("machineIdentityIds") != "mi1,mi2" || q.Get("limit") != "100" {
		t.Fatalf("query: %v", q)
	}
}

func TestDelegationsRemoveSubjectsReturnsNilOnceTheLastOneIsGone(t *testing.T) {
	server := newAgentVaultServer(t, jsonResponse(200, map[string]any{"delegation": nil}))
	res, err := newAgentVaultTestClient(t, server.url).Delegations().RemoveSubjects("mi1", AgentVaultChangeDelegationSubjectsOptions{ProjectID: testProjectID, UserIDs: []string{"u1"}})
	if err != nil || res != nil {
		t.Fatalf("got %+v, %v", res, err)
	}
	req := server.only()
	if req.Method != http.MethodDelete || req.Query.Get("projectId") != testProjectID {
		t.Fatalf("request: %+v", req)
	}
	if _, ok := req.Body["projectId"]; ok {
		t.Fatal("projectId was sent in the body")
	}
}

// Requests that must not run twice

func TestRequestsThatMustNotRunTwiceAreNotRetriedOnceTheyMayHaveReachedTheServer(t *testing.T) {
	calls := map[string]func(av AgentVaultInterface) error{
		"create proxy": func(av AgentVaultInterface) error {
			_, err := av.Proxies().Create(AgentVaultCreateProxyOptions{ProjectID: testProjectID, Name: "edge"})
			return err
		},
		"delete proxy":         func(av AgentVaultInterface) error { _, err := av.Proxies().Delete("p1"); return err },
		"rotate a proxy token": func(av AgentVaultInterface) error { _, err := av.Proxies().RotateToken("p1"); return err },
		"create session": func(av AgentVaultInterface) error {
			_, err := av.Sessions().Create(AgentVaultCreateSessionOptions{ProjectID: testProjectID, AccessBundle: "eng"})
			return err
		},
		"create access bundle": func(av AgentVaultInterface) error {
			_, err := av.AccessBundles().Create(AgentVaultCreateAccessBundleOptions{ProjectID: testProjectID, Name: "eng"})
			return err
		},
		"delete access bundle": func(av AgentVaultInterface) error { _, err := av.AccessBundles().Delete("ab1"); return err },
		"grant access bundle": func(av AgentVaultInterface) error {
			_, err := av.AccessBundles().Grant("ab1", AgentVaultChangeGrantsOptions{UserIDs: []string{"u1"}})
			return err
		},
		"revoke access bundle": func(av AgentVaultInterface) error {
			_, err := av.AccessBundles().Revoke("ab1", AgentVaultChangeGrantsOptions{UserIDs: []string{"u1"}})
			return err
		},
		"delete delegation": func(av AgentVaultInterface) error {
			_, err := av.Delegations().Delete("mi1", AgentVaultProjectOptions{ProjectID: testProjectID})
			return err
		},
		"remove delegation subjects": func(av AgentVaultInterface) error {
			_, err := av.Delegations().RemoveSubjects("mi1", AgentVaultChangeDelegationSubjectsOptions{ProjectID: testProjectID, UserIDs: []string{"u1"}})
			return err
		},
	}
	for name, call := range calls {
		t.Run(name, func(t *testing.T) {
			for _, respond := range []func(http.ResponseWriter, *http.Request){resetConnection, jsonResponse(500, map[string]any{"message": "boom"})} {
				server := newAgentVaultServer(t, respond)
				if err := call(newAgentVaultTestClient(t, server.url)); err == nil {
					t.Fatal("want an error")
				}
				if n := len(server.recorded()); n != 1 {
					t.Fatalf("sent %d requests, want 1", n)
				}
			}
		})
	}
}

func TestRequestErrorsKeepWhatTheRequestFailedWith(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	refused := "http://" + listener.Addr().String()
	_ = listener.Close()

	_, err = newAgentVaultTestClient(t, refused).Proxies().Create(AgentVaultCreateProxyOptions{ProjectID: testProjectID, Name: "edge"})
	var opErr *net.OpError
	if !errors.As(err, &opErr) || opErr.Op != "dial" {
		t.Fatalf("a refused connection gave %v, want a *net.OpError from dial", err)
	}

	server := newAgentVaultServer(t, resetConnection)
	_, err = newAgentVaultTestClient(t, server.url).Proxies().Create(AgentVaultCreateProxyOptions{ProjectID: testProjectID, Name: "edge"})
	var requestErr *sdkErrors.RequestError
	if !errors.As(err, &requestErr) || errors.Unwrap(requestErr) == nil || (errors.As(err, &opErr) && opErr.Op == "dial") {
		t.Fatalf("a dropped connection gave %v, want a request error that isn't from dial", err)
	}
}

func TestOtherRequestsAreStillRetried(t *testing.T) {
	server := newAgentVaultServer(t, resetConnection)
	if _, err := newAgentVaultTestClient(t, server.url).Sessions().Revoke("s1"); err == nil {
		t.Fatal("want an error")
	}
	// One request and the 2 retries the test client allows.
	if n := len(server.recorded()); n != 3 {
		t.Fatalf("sent %d requests, want 3", n)
	}
}

func TestRequestsThatMustNotRunTwiceAreRetriedWhenTheConnectionNeverOpened(t *testing.T) {
	// Reserve a port, then free it, so the first attempt is refused before any of the request is sent.
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := listener.Addr().String()
	_ = listener.Close()

	server := &agentVaultServer{t: t, respond: jsonResponse(200, map[string]any{"session": exampleCreatedSession})}
	httpServer := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.mu.Lock()
		server.requests = append(server.requests, recordedRequest{Method: r.Method, Path: r.URL.EscapedPath()})
		server.mu.Unlock()
		server.respond(w, r)
	})}
	t.Cleanup(func() { _ = httpServer.Close() })
	go func() {
		time.Sleep(200 * time.Millisecond)
		reopened, err := net.Listen("tcp", addr)
		if err != nil {
			t.Errorf("couldn't reopen %s: %v", addr, err)
			return
		}
		_ = httpServer.Serve(reopened)
	}()

	// Waits 100ms, then 200ms, then 400ms, so a retry comes after the server is back.
	av := newAgentVaultTestClientWithBackoff(t, "http://"+addr, &ExponentialBackoffStrategy{BaseDelay: 100 * time.Millisecond, MaxDelay: time.Second, MaxRetries: 3})
	res, err := av.Sessions().Create(AgentVaultCreateSessionOptions{ProjectID: testProjectID, AccessBundle: "eng"})
	if err != nil {
		t.Fatal(err)
	}
	if res.Token != "session-token" || len(server.recorded()) != 1 {
		t.Fatalf("got %+v after %d requests", res, len(server.recorded()))
	}
}

func TestRequestsThatMustNotRunTwiceSendTheTokenTheClientHasWhenTheyAreSent(t *testing.T) {
	var authorization string
	server := newAgentVaultServer(t, func(w http.ResponseWriter, r *http.Request) {
		authorization = r.Header.Get("Authorization")
		jsonResponse(200, map[string]any{})(w, r)
	})
	client := NewInfisicalClient(context.Background(), Config{SiteUrl: server.url, AutoTokenRefresh: BoolPtr(false), SilentMode: true}).(*InfisicalClient)
	client.Auth().SetAccessToken("first-token")
	once := client.onceHTTPClient()
	// As if the background job signed in again between two attempts.
	client.Auth().SetAccessToken("second-token")

	if _, err := once.R().Post("/v1/agent-vault/sessions"); err != nil {
		t.Fatal(err)
	}
	if authorization != "Bearer second-token" {
		t.Fatalf("sent Authorization %q", authorization)
	}
}

func TestRequestsThatMustNotRunTwiceWaitBetweenRetriesAsConfigured(t *testing.T) {
	client := NewInfisicalClient(context.Background(), Config{
		SiteUrl: "http://127.0.0.1:1", AutoTokenRefresh: BoolPtr(false), SilentMode: true,
		RetryRequestsConfig: &RetryRequestsConfig{
			ExponentialBackoff: &ExponentialBackoffStrategy{BaseDelay: 3 * time.Second, MaxDelay: 3 * time.Second, MaxRetries: 2},
		},
	}).(*InfisicalClient)
	once := client.onceHTTPClient()
	if once.RetryAfter == nil {
		t.Fatal("the client doesn't compute retry delays from RetryRequestsConfig")
	}
	// Resty sets Attempt to 1 for the first send. MaxDelay caps the wait at 3s, with up to 20% jitter either way, which
	// a delay that ignored RetryRequestsConfig wouldn't match.
	req := once.R()
	req.Attempt = 1
	delay, err := once.RetryAfter(once, &resty.Response{Request: req})
	if err != nil || delay < 2400*time.Millisecond || delay > 3600*time.Millisecond {
		t.Fatalf("got a delay of %s (%v), want about 3s", delay, err)
	}
}

// IDs and projects

func TestIDsThatWouldReachAnotherRouteAreRefusedBeforeAnyRequest(t *testing.T) {
	for _, id := range []string{"", ".", ".."} {
		server := newAgentVaultServer(t, jsonResponse(200, map[string]any{}))
		av := newAgentVaultTestClient(t, server.url)
		if _, err := av.AccessBundles().Revoke(id, AgentVaultChangeGrantsOptions{UserIDs: []string{"u1"}}); err == nil {
			t.Errorf("ID %q was accepted", id)
		}
		if _, err := av.Sessions().Get(id); err == nil {
			t.Errorf("ID %q was accepted", id)
		}
		if n := len(server.recorded()); n != 0 {
			t.Errorf("ID %q sent %d requests", id, n)
		}
	}
}

func TestIDsAreSentAsOnePathSegment(t *testing.T) {
	server := newAgentVaultServer(t, jsonResponse(200, map[string]any{"session": exampleSession}))
	if _, err := newAgentVaultTestClient(t, server.url).Sessions().Get("a/../b?c"); err != nil {
		t.Fatal(err)
	}
	if got := server.only().Path; got != "/v1/agent-vault/sessions/a%2F..%2Fb%3Fc" {
		t.Fatalf("path: %s", got)
	}
}

func TestTheProjectIDIsRequired(t *testing.T) {
	server := newAgentVaultServer(t, jsonResponse(200, map[string]any{}))
	av := newAgentVaultTestClient(t, server.url)
	errs := []error{}
	_, err := av.Proxies().List(AgentVaultListProxiesOptions{})
	errs = append(errs, err)
	_, err = av.Sessions().Create(AgentVaultCreateSessionOptions{AccessBundle: "eng"})
	errs = append(errs, err)
	_, err = av.AccessBundles().Create(AgentVaultCreateAccessBundleOptions{Name: "eng"})
	errs = append(errs, err)
	_, err = av.Delegations().Get("mi1", AgentVaultProjectOptions{})
	errs = append(errs, err)
	for i, err := range errs {
		if err == nil {
			t.Errorf("call %d was accepted without a project", i)
		}
	}
	if n := len(server.recorded()); n != 0 {
		t.Fatalf("sent %d requests", n)
	}
}

func TestListsSendTheProjectIDInTheQueryAndCreatesInTheBody(t *testing.T) {
	server := newAgentVaultServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			jsonResponse(200, map[string]any{"sessions": []any{}, "totalCount": 0})(w, r)
			return
		}
		jsonResponse(200, map[string]any{"accessBundle": map[string]any{"id": "ab1", "name": "eng", "description": nil, "createdAt": "2026-11-19T10:00:00.000Z", "updatedAt": "2026-11-19T10:00:00.000Z"}})(w, r)
	})
	av := newAgentVaultTestClient(t, server.url)
	if _, err := av.Sessions().List(AgentVaultListSessionsOptions{ProjectID: testProjectID}); err != nil {
		t.Fatal(err)
	}
	if _, err := av.AccessBundles().Create(AgentVaultCreateAccessBundleOptions{ProjectID: testProjectID, Name: "eng"}); err != nil {
		t.Fatal(err)
	}
	requests := server.recorded()
	if requests[0].Query.Get("projectId") != testProjectID || requests[1].Body["projectId"] != testProjectID || requests[1].Query.Has("projectId") {
		t.Fatalf("requests: %+v", requests)
	}
}
