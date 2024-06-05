package infisical

type UniversalAuth struct {
	ClientID     string
	ClientSecret string
}

type GcpIdTokenAuth struct {
	IdentityID string
}

type GcpIamAuth struct {
	IdentityID                string
	ServiceAccountKeyFilePath string
}

type AwsIamAuth struct {
	IdentityID string
}

type AzureAuth struct {
	IdentityID string
}

type KubernetesAuth struct {
	IdentityID              string
	ServiceAccountTokenPath string
}

type Authentication struct {
	UniversalAuth UniversalAuth
	GCPIdToken    GcpIdTokenAuth
	GCPIam        GcpIamAuth
	AWSIam        AwsIamAuth
	Azure         AzureAuth
	Kubernetes    KubernetesAuth

	AccessToken string
}

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

type SecretImport struct {
	SecretPath  string   `json:"secretPath"`
	Environment string   `json:"environment"`
	FolderID    string   `json:"folderId"`
	Secrets     []Secret `json:"secrets"`
}
