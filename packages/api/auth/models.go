package api

import (
	"time"

	"github.com/infisical/go-sdk/packages/models"
)

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

type MachineIdentityAuthLoginResponse struct {
	AccessToken       string `json:"accessToken"`
	ExpiresIn         int64  `json:"expiresIn"`
	AccessTokenMaxTTL int64  `json:"accessTokenMaxTTL"`
	TokenType         string `json:"tokenType"`
}

func (m MachineIdentityAuthLoginResponse) ToMachineIdentity() models.MachineIdentityCredential {
	return models.MachineIdentityCredential{
		AccessToken:       m.AccessToken,
		ExpiresIn:         time.Duration(m.ExpiresIn * int64(time.Second)),
		AccessTokenMaxTTL: time.Duration(m.AccessTokenMaxTTL * int64(time.Second)),
		TokenType:         m.TokenType,
	}
}

type TokenRenewRequest struct {
	AccessToken string `json:"accessToken"`
}
