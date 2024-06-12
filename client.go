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
	folders FoldersInterface
	auth    AuthInterface
}

type InfisicalClientInterface interface {
	UpdateConfiguration(config Config)
	Secrets() SecretsInterface
	Folders() FoldersInterface
	Auth() AuthInterface
}

type Config struct {
	SiteUrl   string
	UserAgent string // optional, we set this when instantiating the client in the k8s operator / cli.
}

func (c *InfisicalClient) setAccessToken(accessToken string, authMethod util.AuthMethod) {
	// We check if the accessToken starts with "Bearer ", and if it does, we remove it from the accessToken
	const bearerPrefix = "Bearer "
	if len(accessToken) >= len(bearerPrefix) && accessToken[:len(bearerPrefix)] == bearerPrefix {
		accessToken = accessToken[len(bearerPrefix):]
	}

	c.authMethod = authMethod
	c.httpClient.SetAuthScheme("Bearer")
	c.httpClient.SetAuthToken(accessToken)
}

func NewInfisicalClient(config Config) InfisicalClientInterface {
	client := &InfisicalClient{}

	client.UpdateConfiguration(config) // set httpClient and config

	// add interfaces here
	client.secrets = &Secrets{client: client}
	client.folders = &Folders{client: client}
	client.auth = &Auth{client: client}

	return client

}

func (c *InfisicalClient) UpdateConfiguration(config Config) {

	if config.UserAgent == "" {
		config.UserAgent = "infisical-go-sdk"
	}
	if config.SiteUrl == "" {
		config.SiteUrl = util.DEFAULT_INFISICAL_API_URL
	}
	config.SiteUrl = util.AppendAPIEndpoint(config.SiteUrl)

	c.config = config

	if c.httpClient == nil {
		c.httpClient = resty.New().SetHeader("User-Agent", config.UserAgent).SetBaseURL(config.SiteUrl)
	} else {
		c.httpClient.SetHeader("User-Agent", config.UserAgent).SetBaseURL(config.SiteUrl)
	}
}

func (c *InfisicalClient) Secrets() SecretsInterface {
	return c.secrets
}

func (c *InfisicalClient) Folders() FoldersInterface {
	return c.folders
}

func (c *InfisicalClient) Auth() AuthInterface {
	return c.auth
}
