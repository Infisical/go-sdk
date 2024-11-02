package models

type Secret struct {
	ID            string `json:"id"`
	Workspace     string `json:"workspace"`
	Environment   string `json:"environment"`
	Version       int    `json:"version"`
	Type          string `json:"type"`
	SecretKey     string `json:"secretKey"`
	SecretValue   string `json:"secretValue"`
	SecretComment string `json:"secretComment"`
	SecretPath    string `json:"secretPath,omitempty"`
}

type ListSecretsWithETagResult struct {
	Secrets    []Secret
	ETag       string
	IsModified bool
}

type SecretImport struct {
	SecretPath  string   `json:"secretPath"`
	Environment string   `json:"environment"`
	FolderID    string   `json:"folderId"`
	Secrets     []Secret `json:"secrets"`
}
