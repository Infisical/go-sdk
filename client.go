package infisical

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/infisical/go-sdk/packages/util"
)

type Client struct {
	authMethod util.AuthMethod
	httpClient *resty.Client
	config     Config

	Secrets SecretsInterface
}

type ClientInterface interface {
	Secrets() SecretsInterface
}

type Config struct {
	CacheTtl  int
	SiteUrl   string
	Auth      Authentication
	UserAgent string // optional, we set this when instantiating the client in the k8s operator / cli.
}

func ValidateAuth(ao *Authentication) (util.AuthMethod, error) {
	if ao.UniversalAuth.ClientID != "" && ao.UniversalAuth.ClientSecret != "" {
		return util.UNIVERSAL_AUTH, nil
	} else if ao.GCPIdToken.IdentityID != "" {
		return util.GCP_ID_TOKEN, nil
	} else if ao.GCPIam.IdentityID != "" && ao.GCPIam.ServiceAccountKeyFilePath != "" {
		return util.GCP_IAM, nil
	} else if ao.AWSIam.IdentityID != "" {
		return util.AWS_IAM, nil
	} else if ao.Azure.IdentityID != "" {
		return util.AZURE, nil
	} else if ao.Kubernetes.IdentityID != "" && ao.Kubernetes.ServiceAccountTokenPath != "" {
		return util.KUBERNETES, nil
	} else if ao.AccessToken != "" {
		return util.ACCESS_TOKEN, nil
	}

	// Read all potential environment variables

	// Universal auth:
	universalAuthClientIdEnv := os.Getenv(util.INFISICAL_UNIVERSAL_AUTH_CLIENT_ID_ENV_NAME)
	universalAuthClientSecretEnv := os.Getenv(util.INFISICAL_UNIVERSAL_AUTH_CLIENT_SECRET_ENV_NAME)

	// GCP ID token auth:
	gcpIdTokenIdentityIdEnv := os.Getenv(util.INFISICAL_GCP_AUTH_IDENTITY_ID_ENV_NAME)

	// GCP IAM auth:
	gcpIamIdentityIdEnv := os.Getenv(util.INFISICAL_GCP_IAM_SERVICE_ACCOUNT_KEY_FILE_PATH_ENV_NAME)

	// AWS IAM auth:
	awsIamIdentityIdEnv := os.Getenv(util.INFISICAL_AWS_IAM_AUTH_IDENTITY_ID_ENV_NAME)

	// Azure auth:
	azureIdentityIdEnv := os.Getenv(util.INFISICAL_AZURE_AUTH_IDENTITY_ID_ENV_NAME)

	// Kubernetes auth:
	kubernetesIdentityIdEnv := os.Getenv(util.INFISICAL_KUBERNETES_IDENTITY_ID_ENV_NAME)
	kubernetesServiceAccountTokenPathEnv := os.Getenv(util.INFISICAL_KUBERNETES_SERVICE_ACCOUNT_TOKEN_PATH_ENV_NAME)

	// Access token:
	accessTokenEnv := os.Getenv(util.INFISICAL_ACCESS_TOKEN_ENV_NAME)

	if universalAuthClientIdEnv != "" && universalAuthClientSecretEnv != "" {
		ao.UniversalAuth.ClientID = universalAuthClientIdEnv
		ao.UniversalAuth.ClientSecret = universalAuthClientSecretEnv
		return util.UNIVERSAL_AUTH, nil
	} else if gcpIdTokenIdentityIdEnv != "" {
		ao.GCPIdToken.IdentityID = gcpIdTokenIdentityIdEnv
		return util.GCP_ID_TOKEN, nil
	} else if gcpIamIdentityIdEnv != "" {
		ao.GCPIam.IdentityID = gcpIamIdentityIdEnv
		return util.GCP_IAM, nil
	} else if awsIamIdentityIdEnv != "" {
		ao.AWSIam.IdentityID = awsIamIdentityIdEnv
		return util.AWS_IAM, nil
	} else if azureIdentityIdEnv != "" {
		ao.Azure.IdentityID = azureIdentityIdEnv
		return util.AZURE, nil
	} else if kubernetesIdentityIdEnv != "" && kubernetesServiceAccountTokenPathEnv != "" {
		ao.Kubernetes.IdentityID = kubernetesIdentityIdEnv
		ao.Kubernetes.ServiceAccountTokenPath = kubernetesServiceAccountTokenPathEnv
		return util.KUBERNETES, nil
	} else if accessTokenEnv != "" {
		ao.AccessToken = accessTokenEnv
		return util.ACCESS_TOKEN, nil
	}

	return "", fmt.Errorf("no authentication method is set")
}

func NewInfisicalClient(config Config) (*Client, error) {

	if config.UserAgent == "" {
		config.UserAgent = "infisical-go-sdk"
	}

	if config.SiteUrl == "" {
		config.SiteUrl = util.DEFAULT_INFISICAL_API_URL
	}

	config.SiteUrl = util.AppendAPIEndpoint(config.SiteUrl)

	// Auth method validation
	authMethod, err := ValidateAuth(&config.Auth)

	if err != nil {
		return nil, fmt.Errorf("error while instantiating Infisical client: %v", err)
	}

	client := &Client{
		authMethod: authMethod,
		config:     config,
		httpClient: resty.New().SetHeader("User-Agent", config.UserAgent).SetBaseURL(config.SiteUrl),
	}
	err = client.authenticateHttpClient()

	if err != nil {
		return nil, fmt.Errorf("error while instantiating Infisical client: %v", err)
	}

	client.Secrets = &Secrets{client: client}
	// add other interfaces here

	return client, nil

}
