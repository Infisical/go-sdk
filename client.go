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
	if ao.universalAuth.clientId != "" && ao.universalAuth.clientSecret != "" {
		return util.UNIVERSAL_AUTH, nil
	} else if ao.gcpIdTokenAuth.identityId != "" {
		return util.GCP_ID_TOKEN, nil
	} else if ao.gcpIamAuth.identityId != "" && ao.gcpIamAuth.serviceAccountKeyFilePath != "" {
		return util.GCP_IAM, nil
	} else if ao.awsIamAuth.identityId != "" {
		return util.AWS_IAM, nil
	} else if ao.azureAuth.identityId != "" {
		return util.AZURE, nil
	} else if ao.kubernetesAuth.identityId != "" && ao.kubernetesAuth.serviceAccountTokenPath != "" {
		return util.KUBERNETES, nil
	} else if ao.accessToken != "" {
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
		ao.universalAuth.clientId = universalAuthClientIdEnv
		ao.universalAuth.clientSecret = universalAuthClientSecretEnv
		return util.UNIVERSAL_AUTH, nil
	} else if gcpIdTokenIdentityIdEnv != "" {
		ao.gcpIdTokenAuth.identityId = gcpIdTokenIdentityIdEnv
		return util.GCP_ID_TOKEN, nil
	} else if gcpIamIdentityIdEnv != "" {
		ao.gcpIamAuth.identityId = gcpIamIdentityIdEnv
		return util.GCP_IAM, nil
	} else if awsIamIdentityIdEnv != "" {
		ao.awsIamAuth.identityId = awsIamIdentityIdEnv
		return util.AWS_IAM, nil
	} else if azureIdentityIdEnv != "" {
		ao.azureAuth.identityId = azureIdentityIdEnv
		return util.AZURE, nil
	} else if kubernetesIdentityIdEnv != "" && kubernetesServiceAccountTokenPathEnv != "" {
		ao.kubernetesAuth.identityId = kubernetesIdentityIdEnv
		ao.kubernetesAuth.serviceAccountTokenPath = kubernetesServiceAccountTokenPathEnv
		return util.KUBERNETES, nil
	} else if accessTokenEnv != "" {
		ao.accessToken = accessTokenEnv
		return util.ACCESS_TOKEN, nil
	}

	return "", fmt.Errorf("no authentication method is set")
}

func NewInfisicalClient(config Config) (*Client, error) {

	if config.UserAgent == "" {
		config.UserAgent = "infisical-go-sdk"
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
