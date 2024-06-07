package infisical

import (
	"github.com/go-resty/resty/v2"
	"github.com/infisical/go-sdk/packages/util"
)

type InfisicalClient struct {
	authMethod util.AuthMethod
	httpClient *resty.Client
	config     Config

	secrets SecretsInterface
	auth    AuthInterface
}

type InfisicalClientInterface interface {
	Secrets() SecretsInterface
	Auth() AuthInterface
}

type Config struct {
	SiteUrl   string
	UserAgent string // optional, we set this when instantiating the client in the k8s operator / cli.
}

func (c *InfisicalClient) setAccessToken(accessToken string, authMethod util.AuthMethod) {
	c.authMethod = authMethod
	c.httpClient.SetAuthScheme("Bearer")
	c.httpClient.SetAuthToken(accessToken)
}

func NewInfisicalClient(config Config) (InfisicalClientInterface, error) {

	if config.UserAgent == "" {
		config.UserAgent = "infisical-go-sdk"
	}
	if config.SiteUrl == "" {
		config.SiteUrl = util.DEFAULT_INFISICAL_API_URL
	}

	config.SiteUrl = util.AppendAPIEndpoint(config.SiteUrl)

	client := &InfisicalClient{
		config:     config,
		httpClient: resty.New().SetHeader("User-Agent", config.UserAgent).SetBaseURL(config.SiteUrl),
	}

	// add interfaces here
	client.secrets = &Secrets{client: client}

	return client, nil

}

func (c *InfisicalClient) Secrets() SecretsInterface {
	return c.secrets
}

func (c *InfisicalClient) Auth() AuthInterface {
	return c.auth
}
