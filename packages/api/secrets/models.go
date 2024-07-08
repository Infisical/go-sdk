package api

import "github.com/infisical/go-sdk/packages/models"

// List secrets
type ListSecretsV3RawRequest struct {
	AttachToProcessEnv bool `json:"-"`

	// ProjectId and ProjectSlug are used to fetch secrets from the project. Only one of them is required.
	ProjectID              string `json:"workspaceId,omitempty"`
	ProjectSlug            string `json:"workspaceSlug,omitempty"`
	Environment            string `json:"environment"`
	ExpandSecretReferences bool   `json:"expandSecretReferences"`
	IncludeImports         bool   `json:"include_imports"`
	Recursive              bool   `json:"recursive"`
	SecretPath             string `json:"secretPath,omitempty"`
}

type ListSecretsV3RawResponse struct {
	Secrets []models.Secret       `json:"secrets"`
	Imports []models.SecretImport `json:"imports"`
}

// Retrieve secret

type RetrieveSecretV3RawRequest struct {
	SecretKey string `json:"-"`

	ProjectID      string `json:"workspaceId"`
	Environment    string `json:"environment"`
	SecretPath     string `json:"secretPath,omitempty"`
	Type           string `json:"type,omitempty"`
	IncludeImports bool   `json:"include_imports"`
	Version        int    `json:"version"`
}

type RetrieveSecretV3RawResponse struct {
	Secret models.Secret `json:"secret"`
}

// Update secret
type UpdateSecretV3RawRequest struct {
	SecretKey string `json:"-"`

	ProjectID   string `json:"workspaceId"`
	Environment string `json:"environment"`
	SecretPath  string `json:"secretPath,omitempty"`
	Type        string `json:"type,omitempty"`

	NewSecretValue           string `json:"secretValue,omitempty"`
	NewSkipMultilineEncoding bool   `json:"skipMultilineEncoding,omitempty"`
}

type UpdateSecretV3RawResponse struct {
	Secret models.Secret `json:"secret"`
}

// Create secret
type CreateSecretV3RawRequest struct {
	SecretKey string `json:"-"`

	ProjectID             string `json:"workspaceId"`
	Environment           string `json:"environment"`
	SecretPath            string `json:"secretPath,omitempty"`
	Type                  string `json:"type,omitempty"`
	SecretComment         string `json:"secretComment,omitempty"`
	SkipMultiLineEncoding bool   `json:"skipMultilineEncoding"`
	SecretValue           string `json:"secretValue"`
}

type CreateSecretV3RawResponse struct {
	Secret models.Secret `json:"secret"`
}

// Delete secret
type DeleteSecretV3RawRequest struct {
	SecretKey string `json:"-"`

	ProjectID   string `json:"workspaceId"`
	Environment string `json:"environment"`
	SecretPath  string `json:"secretPath,omitempty"`
	Type        string `json:"type,omitempty"`
}

type DeleteSecretV3RawResponse struct {
	Secret models.Secret `json:"secret"`
}
