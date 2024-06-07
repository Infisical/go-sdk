package api

// Universal auth:
type UniversalAuthLoginRequest struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

// Kubernetes auth:
type KubernetesAuthLoginRequest struct {
	IdentityID string `json:"identityId"`
	JWT        string `json:"jwt"`
}

type AzureAuthLoginRequest struct {
	IdentityID string `json:"identityId"`
	JWT        string `json:"jwt"`
}

type GCPAuthLoginRequest struct {
	IdentityID string `json:"identityId"`
	JWT        string `json:"jwt"`
}

type GenericAuthLoginResponse struct {
	AccessToken string `json:"accessToken"`
}
