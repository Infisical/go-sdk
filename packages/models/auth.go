package models

type TokenType string

const (
	TokenTypeBearer TokenType = "Bearer"
)

type UniversalAuthCredential struct {
	AccessToken       string    `json:"accessToken"`
	ExpiresIn         int64     `json:"expiresIn"`
	AccessTokenMaxTTL int64     `json:"accessTokenMaxTTL"`
	TokenType         TokenType `json:"tokenType"`
}
