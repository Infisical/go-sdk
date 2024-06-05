package api

import "github.com/infisical/go-sdk/packages/models"

type ListSecretsRequest struct {
	AttachToProcessEnv bool

	// ProjectId and ProjectSlug are used to fetch secrets from the project. Only one of them is required.
	ProjectId   string `json:"workspaceId"`
	ProjectSlug string `json:"workspaceSlug"`

	Environment            string `json:"environment"`
	ExpandSecretReferences bool   `json:"expandSecretReferences"`
	IncludeImports         bool   `json:"include_imports"`
	Recursive              bool   `json:"recursive"`
	SecretPath             string `json:"secretPath"`
}

type ListSecretsResponse struct {
	Secrets []models.Secret       `json:"secrets"`
	Imports []models.SecretImport `json:"imports"`
}
