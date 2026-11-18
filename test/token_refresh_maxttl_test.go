package test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	infisical "github.com/infisical/go-sdk"
)

// An accessTokenMaxTTL of 0 means the token has no maximum lifetime, so the SDK
// should keep renewing it instead of performing a full re-authentication every
// time the access token expires.
func TestTokenRefreshRenewsWhenMaxTTLIsZero(t *testing.T) {
	var mu sync.Mutex
	var logins, renewals int

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v1/auth/universal-auth/login":
			mu.Lock()
			logins++
			mu.Unlock()
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"accessToken":       "access-token",
				"expiresIn":         1,
				"accessTokenMaxTTL": 0,
				"tokenType":         "Bearer",
			})
		case "/api/v1/auth/token/renew":
			mu.Lock()
			renewals++
			mu.Unlock()
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"accessToken":       "access-token",
				"expiresIn":         1,
				"accessTokenMaxTTL": 0,
				"tokenType":         "Bearer",
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := infisical.NewInfisicalClient(ctx, infisical.Config{
		SiteUrl:    server.URL,
		SilentMode: true,
	})

	if _, err := client.Auth().UniversalAuthLogin("client-id", "client-secret"); err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	// Let the access token expire so the background refresh runs at least once.
	time.Sleep(2500 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if renewals == 0 {
		t.Fatalf("Expected the SDK to renew the token, but it never called /token/renew (logins=%d)", logins)
	}
	if logins != 1 {
		t.Fatalf("Expected a single login, but the SDK re-authenticated %d times instead of renewing", logins)
	}
}
