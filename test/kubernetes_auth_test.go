package test

import (
	"testing"

	infisical "github.com/infisical/go-sdk"
)

func TestKubernetesAuthLogin(t *testing.T) {

	t.Skip("Skipping Kubernetes Auth test -- requires running in a Kubernetes cluster")

	client, err := infisical.NewInfisicalClient(infisical.Config{
		SiteUrl: "http://localhost:8080",
		Auth: infisical.Authentication{
			Kubernetes: infisical.KubernetesAuth{
				IdentityID:              "K8_MACHINE_IDENTITY_ID",
				ServiceAccountTokenPath: "/var/run/secrets/kubernetes.io/serviceaccount/token", // Optional
			},
		},
	})

	if err != nil {
		t.Fatalf("Failed to create Infisical client: %v", err)
	}

	secrets, err := client.Secrets().List(infisical.ListSecretsOptions{
		ProjectID:   "PROJECT_ID",
		Environment: "ENV_SLUG",
	})

	if err != nil {
		t.Fatalf("Failed to list secrets: %v", err)
	}

	if len(secrets) == 0 {
		t.Fatalf("No secrets found")
	}

	t.Logf("Secrets: %v", secrets)

}
