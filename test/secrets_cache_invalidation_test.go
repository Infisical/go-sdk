package test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	infisical "github.com/infisical/go-sdk"
)

func TestCacheIsInvalidatedByWrites(t *testing.T) {
	var gets int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodGet {
			atomic.AddInt32(&gets, 1)
		}
		_, _ = w.Write([]byte(`{"secret":{"secretKey":"K","secretValue":"v"},"secrets":[]}`))
	}))
	defer srv.Close()

	newClient := func() infisical.InfisicalClientInterface {
		c := infisical.NewInfisicalClient(context.Background(), infisical.Config{
			SiteUrl:              srv.URL,
			AutoTokenRefresh:     infisical.BoolPtr(false),
			CacheExpiryInSeconds: 60,
		})
		c.Auth().SetAccessToken("token")
		return c
	}
	var client infisical.InfisicalClientInterface

	retrieve := func() {
		t.Helper()
		_, err := client.Secrets().Retrieve(infisical.RetrieveSecretOptions{
			SecretKey: "K", Environment: "dev", SecretPath: "/", ProjectID: "p",
		})
		if err != nil {
			t.Fatalf("retrieve: %v", err)
		}
	}

	writes := map[string]func() error{
		"update": func() error {
			_, err := client.Secrets().Update(infisical.UpdateSecretOptions{
				SecretKey: "K", Environment: "dev", SecretPath: "/", ProjectID: "p", NewSecretValue: "w",
			})
			return err
		},
		"create": func() error {
			_, err := client.Secrets().Create(infisical.CreateSecretOptions{
				SecretKey: "K", SecretValue: "w", Environment: "dev", SecretPath: "/", ProjectID: "p",
			})
			return err
		},
		"delete": func() error {
			_, err := client.Secrets().Delete(infisical.DeleteSecretOptions{
				SecretKey: "K", Environment: "dev", SecretPath: "/", ProjectID: "p",
			})
			return err
		},
		"batch create": func() error {
			_, err := client.Secrets().Batch().Create(infisical.BatchCreateSecretsOptions{
				Environment: "dev", SecretPath: "/", ProjectID: "p",
				Secrets: []infisical.BatchCreateSecret{{SecretKey: "K", SecretValue: "w"}},
			})
			return err
		},
	}

	for name, write := range writes {
		client = newClient()
		atomic.StoreInt32(&gets, 0)
		retrieve()
		retrieve()
		if got := atomic.LoadInt32(&gets); got != 1 {
			t.Fatalf("%s: expected the second read to be served from cache, got %d requests", name, got)
		}
		if err := write(); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		retrieve()
		if got := atomic.LoadInt32(&gets); got != 2 {
			t.Fatalf("%s: read after write must reach the API, got %d requests", name, got)
		}
	}
}
