package test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	infisical "github.com/infisical/go-sdk"
)

// Kubernetes rotates projected service account tokens, so re-authentication must
// read the token from disk again instead of reusing the one from the initial login.
func TestKubernetesAuthReAuthReadsRotatedServiceAccountToken(t *testing.T) {
	var mu sync.Mutex
	var receivedJWTs []string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/auth/kubernetes-auth/login" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var body struct {
			JWT string `json:"jwt"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		mu.Lock()
		receivedJWTs = append(receivedJWTs, body.JWT)
		mu.Unlock()

		// max TTL equal to the token TTL forces re-authentication (not renewal) once the token is expiring
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"accessToken":       "access-token",
			"expiresIn":         6,
			"accessTokenMaxTTL": 6,
			"tokenType":         "Bearer",
		}); err != nil {
			t.Errorf("Failed to encode login response: %v", err)
		}
	}))
	defer server.Close()

	tokenPath := filepath.Join(t.TempDir(), "token")
	if err := os.WriteFile(tokenPath, []byte("initial-sa-token"), 0o600); err != nil {
		t.Fatalf("Failed to write service account token: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := infisical.NewInfisicalClient(ctx, infisical.Config{
		SiteUrl:    server.URL,
		SilentMode: true,
	})

	if _, err := client.Auth().KubernetesAuthLogin("identity-id", tokenPath); err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	if err := os.WriteFile(tokenPath, []byte("rotated-sa-token"), 0o600); err != nil {
		t.Fatalf("Failed to rotate service account token: %v", err)
	}

	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		mu.Lock()
		jwts := append([]string{}, receivedJWTs...)
		mu.Unlock()

		if len(jwts) >= 2 {
			if jwts[1] != "rotated-sa-token" {
				t.Fatalf("Expected re-authentication to use the rotated service account token, got %q", jwts[1])
			}
			return
		}
		time.Sleep(100 * time.Millisecond)
	}

	t.Fatalf("Timed out waiting for re-authentication")
}
