package models

type TokenType string

const (
	BEARER_TOKEN_TYPE TokenType = "Bearer"
)

type UniversalAuthCredential struct {
	ClientID     string
	ClientSecret string
}

type AccessTokenCredential struct {
	AccessToken string
}

type GCPIDTokenCredential struct {
	IdentityID string
}

type GCPIAMCredential struct {
	IdentityID                string
	ServiceAccountKeyFilePath string
}

type AWSIAMCredential struct {
	IdentityID string
}

type KubernetesCredential struct {
	IdentityID          string
	ServiceAccountToken string
}

type AzureCredential struct {
	IdentityID string
	Resource   string

	// UseWorkloadIdentity is a tri-state pointer: nil means "auto-detect from env vars",
	// while a non-nil value forces the corresponding mode on re-authentication.
	UseWorkloadIdentity *bool

	// IMDS selectors persisted so background re-auth keeps targeting the same managed
	// identity that the initial login used.
	IMDSClientID string
	IMDSObjectID string

	// Workload Identity overrides persisted alongside the credential. Empty fields
	// fall back to the AZURE_* environment variables, matching the initial login.
	WIClientID      string
	WITenantID      string
	WITokenFilePath string
	WIAuthorityHost string
}

type OIDCCredential struct {
	IdentityID string
	JWT        string
}

type JWTCredential struct {
	IdentityID string
	JWT        string
}

type LDAPCredential struct {
	IdentityID string
	Username   string
	Password   string
}

type OCICredential struct {
	IdentityID  string
	PrivateKey  string
	Fingerprint string
	UserID      string
	TenancyID   string
	Region      string
	Passphrase  *string
}
