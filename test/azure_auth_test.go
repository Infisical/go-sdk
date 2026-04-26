package test

import (
	"context"
	"fmt"
	"os"
	"testing"

	infisical "github.com/infisical/go-sdk"
)

// These tests are skipped by default; they are integration-style smoke tests intended to
// be run manually on an Azure VM / AKS pod with the corresponding credentials available.
// Set INFISICAL_AZURE_TEST_RUN=1 (and the relevant inputs) to opt in.

const azureTestRunEnv = "INFISICAL_AZURE_TEST_RUN"

func skipUnlessAzureIntegration(t *testing.T) {
	t.Helper()
	if os.Getenv(azureTestRunEnv) == "" {
		t.Skipf("set %s=1 to run Azure integration tests", azureTestRunEnv)
	}
}

// TestAzureAuthLogin_Default exercises the legacy IMDS path with no options. Requires a
// VM with a single managed identity attached.
func TestAzureAuthLogin_Default(t *testing.T) {
	skipUnlessAzureIntegration(t)

	siteURL := os.Getenv("INFISICAL_TEST_SITE_URL")
	identityID := os.Getenv("INFISICAL_AZURE_AUTH_IDENTITY_ID")
	if siteURL == "" || identityID == "" {
		t.Skip("INFISICAL_TEST_SITE_URL and INFISICAL_AZURE_AUTH_IDENTITY_ID are required")
	}

	client := infisical.NewInfisicalClient(context.Background(), infisical.Config{SiteUrl: siteURL})
	cred, err := client.Auth().AzureAuthLogin(identityID, "")
	if err != nil {
		t.Fatalf("AzureAuthLogin returned error: %v", err)
	}
	fmt.Println("Obtained access token via IMDS (default):", cred.AccessToken)
}

// TestAzureAuthLogin_IMDSWithClientID pins IMDS to a specific user-assigned managed
// identity. Required when multiple identities are attached to the host.
func TestAzureAuthLogin_IMDSWithClientID(t *testing.T) {
	skipUnlessAzureIntegration(t)

	siteURL := os.Getenv("INFISICAL_TEST_SITE_URL")
	identityID := os.Getenv("INFISICAL_AZURE_AUTH_IDENTITY_ID")
	imdsClientID := os.Getenv("INFISICAL_AZURE_AUTH_CLIENT_ID")
	if siteURL == "" || identityID == "" || imdsClientID == "" {
		t.Skip("INFISICAL_TEST_SITE_URL, INFISICAL_AZURE_AUTH_IDENTITY_ID and INFISICAL_AZURE_AUTH_CLIENT_ID are required")
	}

	client := infisical.NewInfisicalClient(context.Background(), infisical.Config{SiteUrl: siteURL})
	cred, err := client.Auth().AzureAuthLogin(identityID, "", infisical.AzureAuthLoginOptions{
		IMDSClientID: imdsClientID,
	})
	if err != nil {
		t.Fatalf("AzureAuthLogin returned error: %v", err)
	}
	fmt.Println("Obtained access token via IMDS (client_id pinned):", cred.AccessToken)
}

// TestAzureAuthLogin_WorkloadIdentity forces the Workload Identity path. Requires the
// AZURE_CLIENT_ID, AZURE_TENANT_ID and AZURE_FEDERATED_TOKEN_FILE env vars to be set
// (typically injected by the Azure Workload Identity admission webhook).
func TestAzureAuthLogin_WorkloadIdentity(t *testing.T) {
	skipUnlessAzureIntegration(t)

	siteURL := os.Getenv("INFISICAL_TEST_SITE_URL")
	identityID := os.Getenv("INFISICAL_AZURE_AUTH_IDENTITY_ID")
	if siteURL == "" || identityID == "" {
		t.Skip("INFISICAL_TEST_SITE_URL and INFISICAL_AZURE_AUTH_IDENTITY_ID are required")
	}

	useWI := true
	client := infisical.NewInfisicalClient(context.Background(), infisical.Config{SiteUrl: siteURL})
	cred, err := client.Auth().AzureAuthLogin(identityID, "", infisical.AzureAuthLoginOptions{
		UseWorkloadIdentity: &useWI,
	})
	if err != nil {
		t.Fatalf("AzureAuthLogin returned error: %v", err)
	}
	fmt.Println("Obtained access token via Workload Identity:", cred.AccessToken)
}

// TestAzureAuthLogin_ForceIMDSWhenWIEnvsPresent verifies the explicit opt-out path: even
// if the AZURE_* WI env vars are set, passing UseWorkloadIdentity=false routes through IMDS.
func TestAzureAuthLogin_ForceIMDSWhenWIEnvsPresent(t *testing.T) {
	skipUnlessAzureIntegration(t)

	siteURL := os.Getenv("INFISICAL_TEST_SITE_URL")
	identityID := os.Getenv("INFISICAL_AZURE_AUTH_IDENTITY_ID")
	if siteURL == "" || identityID == "" {
		t.Skip("INFISICAL_TEST_SITE_URL and INFISICAL_AZURE_AUTH_IDENTITY_ID are required")
	}

	useWI := false
	client := infisical.NewInfisicalClient(context.Background(), infisical.Config{SiteUrl: siteURL})
	cred, err := client.Auth().AzureAuthLogin(identityID, "", infisical.AzureAuthLoginOptions{
		UseWorkloadIdentity: &useWI,
		IMDSClientID:        os.Getenv("INFISICAL_AZURE_AUTH_CLIENT_ID"),
	})
	if err != nil {
		t.Fatalf("AzureAuthLogin returned error: %v", err)
	}
	fmt.Println("Obtained access token via IMDS (forced):", cred.AccessToken)
}
