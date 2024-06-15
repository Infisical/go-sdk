package models

import "time"

type TokenType string

const (
	BEARER_TOKEN_TYPE TokenType = "Bearer"
)

type MachineIdentityCredential struct {
	AccessToken       string        `json:"accessToken"`
	ExpiresIn         time.Duration `json:"expiresIn"`
	AccessTokenMaxTTL time.Duration `json:"accessTokenMaxTTL"`
	TokenType         string        `json:"tokenType"`
}
