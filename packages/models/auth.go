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
}

type OIDCCredential struct {
	IdentityID string
	JWT        string
}

type JWTCredential struct {
	IdentityID string
	JWT        string
}
