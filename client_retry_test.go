package infisical

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	sdkErrors "github.com/infisical/go-sdk/packages/errors"
)

// statusServer answers every request with status and counts the requests it receives.
func statusServer(t *testing.T, status int, body any) (*httptest.Server, *atomic.Int32) {
	t.Helper()
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(body)
	}))
	t.Cleanup(server.Close)
	return server, &requests
}

func newBackoffTestClient(t *testing.T, siteURL string, backoff *ExponentialBackoffStrategy) *InfisicalClient {
	t.Helper()
	return NewInfisicalClient(context.Background(), Config{
		SiteUrl:             siteURL,
		AutoTokenRefresh:    BoolPtr(false),
		SilentMode:          true,
		RetryRequestsConfig: &RetryRequestsConfig{ExponentialBackoff: backoff},
	}).(*InfisicalClient)
}

func TestBackoffRetriesServerErrorsAndRateLimitsButNotClientErrors(t *testing.T) {
	for status, want := range map[int]int32{400: 1, 401: 1, 403: 1, 404: 1, 409: 1, 422: 1, 429: 3, 500: 3, 503: 3} {
		server, requests := statusServer(t, status, map[string]any{"message": "no"})
		client := newBackoffTestClient(t, server.URL, &ExponentialBackoffStrategy{BaseDelay: time.Millisecond, MaxDelay: time.Millisecond, MaxRetries: 2})
		if _, err := client.httpClient.R().Get("/v1/secrets/missing"); err != nil {
			t.Fatal(err)
		}
		if got := requests.Load(); got != want {
			t.Errorf("a %d was sent %d times, want %d", status, got, want)
		}
	}
}

func TestBackoffRetriesThreeTimesWhenMaxRetriesIsUnset(t *testing.T) {
	server, requests := statusServer(t, 503, map[string]any{"message": "later"})
	client := newBackoffTestClient(t, server.URL, &ExponentialBackoffStrategy{BaseDelay: time.Millisecond, MaxDelay: time.Millisecond})
	if _, err := client.httpClient.R().Get("/v1/secrets/missing"); err != nil {
		t.Fatal(err)
	}
	// The first send and 3 retries.
	if got := requests.Load(); got != 4 {
		t.Fatalf("sent %d requests, want 4", got)
	}
}

func TestBackoffFirstRetryWaitsBaseDelay(t *testing.T) {
	client := newBackoffTestClient(t, "http://127.0.0.1:1", &ExponentialBackoffStrategy{BaseDelay: 3 * time.Second, MaxDelay: time.Minute, MaxRetries: 2})
	// Resty sets Attempt to 1 for the first send, so the first retry waits BaseDelay and the second twice that, with
	// up to 20% jitter either way.
	for attempt, want := range map[int]time.Duration{1: 3 * time.Second, 2: 6 * time.Second} {
		req := client.httpClient.R()
		req.Attempt = attempt
		delay, err := client.httpClient.RetryAfter(client.httpClient, &resty.Response{Request: req})
		if err != nil || delay < want*8/10 || delay > want*12/10 {
			t.Fatalf("after attempt %d got a delay of %s (%v), want about %s", attempt, delay, err, want)
		}
	}
}

func TestBackoffWaitsLessThanASecondWhenBaseDelayIsShorter(t *testing.T) {
	server, requests := statusServer(t, 503, map[string]any{"message": "later"})
	client := newBackoffTestClient(t, server.URL, &ExponentialBackoffStrategy{BaseDelay: time.Millisecond, MaxDelay: time.Millisecond, MaxRetries: 2})
	start := time.Now()
	if _, err := client.httpClient.R().Get("/v1/secrets/missing"); err != nil {
		t.Fatal(err)
	}
	// Resty would otherwise raise each of the 2 waits to 1s.
	if elapsed := time.Since(start); requests.Load() != 3 || elapsed > 900*time.Millisecond {
		t.Fatalf("sent %d requests in %s", requests.Load(), elapsed)
	}
}

func TestValidationErrorsNameEachField(t *testing.T) {
	server, _ := statusServer(t, 422, map[string]any{
		"reqId": "req-1", "statusCode": 422, "error": "ValidationFailure",
		"message": []any{
			map[string]any{"code": "too_small", "message": "Number must be greater than or equal to 1", "path": []any{"limit"}},
			map[string]any{"code": "invalid_type", "message": "Required", "path": []any{"secrets", 0, "secretKey"}},
		},
	})
	client := newBackoffTestClient(t, server.URL, &ExponentialBackoffStrategy{BaseDelay: time.Millisecond, MaxDelay: time.Millisecond, MaxRetries: 2})
	res, err := client.httpClient.R().Get("/v1/secrets")
	if err != nil {
		t.Fatal(err)
	}
	var apiErr *sdkErrors.APIError
	if !errors.As(sdkErrors.NewAPIErrorWithResponse("CallListSecretsV3", res), &apiErr) ||
		apiErr.ErrorMessage != "limit: Number must be greater than or equal to 1; secrets.0.secretKey: Required" {
		t.Fatalf("got %+v", apiErr)
	}
}

func TestRequestErrorsKeepWhatTheRequestFailedWith(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	refused := "http://" + listener.Addr().String()
	_ = listener.Close()

	client := newBackoffTestClient(t, refused, &ExponentialBackoffStrategy{BaseDelay: time.Millisecond, MaxDelay: time.Millisecond, MaxRetries: 1})
	_, err = client.Secrets().Retrieve(RetrieveSecretOptions{SecretKey: "API_KEY", ProjectID: "p1", Environment: "dev"})
	var requestErr *sdkErrors.RequestError
	var opErr *net.OpError
	if !errors.As(err, &requestErr) || !errors.As(err, &opErr) || opErr.Op != "dial" {
		t.Fatalf("a refused connection gave %v, want a request error holding a *net.OpError from dial", err)
	}

	if !errors.Is(sdkErrors.NewRequestError("CallRetrieveSecretV3", context.Canceled), context.Canceled) {
		t.Fatal("a request error hides the context error it failed with")
	}
}
