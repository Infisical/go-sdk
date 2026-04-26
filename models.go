package infisical

import (
	api "github.com/infisical/go-sdk/packages/api/auth"
	"github.com/infisical/go-sdk/packages/errors"
	"github.com/infisical/go-sdk/packages/models"
)

type OciAuthLoginOptions struct {
	IdentityID  string
	PrivateKey  string
	Fingerprint string
	UserID      string
	TenancyID   string
	Region      string
	Passphrase  *string
}

// AzureAuthLoginOptions configures how the SDK obtains the Azure AD access token that is
// exchanged with Infisical's Azure auth method. All fields are optional; when omitted the
// SDK falls back to the matching environment variables, preserving the legacy behaviour.
type AzureAuthLoginOptions struct {
	// IMDSClientID pins the IMDS request to a specific user-assigned managed identity by
	// its AAD client (application) ID. Required when multiple managed identities are
	// attached to the host. Falls back to INFISICAL_AZURE_AUTH_CLIENT_ID.
	IMDSClientID string
	// IMDSObjectID pins the IMDS request by the managed identity's principal (object) ID.
	// Used only when IMDSClientID is empty. Falls back to INFISICAL_AZURE_AUTH_OBJECT_ID.
	IMDSObjectID string

	// UseWorkloadIdentity is a tri-state mode override:
	//   nil   -> auto-detect: use Workload Identity when AZURE_CLIENT_ID, AZURE_TENANT_ID
	//           and AZURE_FEDERATED_TOKEN_FILE are set and the token file exists,
	//           otherwise fall back to IMDS (preserves legacy behaviour).
	//   true  -> force Workload Identity.
	//   false -> force IMDS (explicit opt-out for pods that use Workload Identity for
	//           other clients but want IMDS for Infisical).
	// The INFISICAL_AZURE_AUTH_USE_WORKLOAD_IDENTITY env var has the same semantics and
	// is consulted only when this field is nil.
	UseWorkloadIdentity *bool

	// WIClientID overrides AZURE_CLIENT_ID for Workload Identity. Empty -> env var.
	WIClientID string
	// WITenantID overrides AZURE_TENANT_ID for Workload Identity. Empty -> env var.
	WITenantID string
	// WITokenFilePath overrides AZURE_FEDERATED_TOKEN_FILE for Workload Identity.
	WITokenFilePath string
	// WIAuthorityHost overrides AZURE_AUTHORITY_HOST for Workload Identity.
	WIAuthorityHost string
}

type MachineIdentityCredential = api.MachineIdentityAuthLoginResponse

type Secret = models.Secret
type SecretImport = models.SecretImport

type APIError = errors.APIError
type RequestError = errors.RequestError
type NotModifiedError = errors.NotModifiedError
