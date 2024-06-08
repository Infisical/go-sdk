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

type AwsIamAuthLoginRequest struct {
	HTTPRequestMethod string `json:"iamHttpRequestMethod"`
	IamRequestBody    string `json:"iamRequestBody"`
	IamRequestHeaders string `json:"iamRequestHeaders"`
	IdentityId        string `json:"identityId"`
}

type GCPAuthLoginRequest struct {
	IdentityID string `json:"identityId"`
	JWT        string `json:"jwt"`
}

type GenericAuthLoginResponse struct {
	AccessToken string `json:"accessToken"`
}
