package infisical

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	api "github.com/infisical/go-sdk/packages/api/auth"
	"github.com/infisical/go-sdk/packages/models"
	"github.com/infisical/go-sdk/packages/util"
)

type InfisicalClient struct {
	authMethod       util.AuthMethod
	credential       interface{}
	tokenDetails     MachineIdentityCredential
	lastFetchedTime  time.Time
	firstFetchedTime time.Time

	mu sync.RWMutex

	httpClient *resty.Client
	config     Config

	secrets SecretsInterface
	folders FoldersInterface
	auth    AuthInterface
}

type InfisicalClientInterface interface {
	UpdateConfiguration(config *Config)
	Secrets() SecretsInterface
	Folders() FoldersInterface
	Auth() AuthInterface
}

type Config struct {
	SiteUrl          string
	UserAgent        string // optional, we set this when instantiating the client in the k8s operator / cli.
	AutoTokenRefresh bool   // defaults to trues
	SilentMode       bool   // defaults to false
}

func NewInfisicalClientConfig(options ...func(*Config)) *Config {
	cfg := &Config{
		SiteUrl:          util.DEFAULT_INFISICAL_API_URL,
		UserAgent:        "infisical-go-sdk",
		AutoTokenRefresh: true,
		SilentMode:       false,
	}

	for _, opt := range options {
		opt(cfg)
	}
	return cfg
}

func WithSiteUrl(siteUrl string) func(*Config) {
	return func(s *Config) {
		s.SiteUrl = siteUrl
	}
}

func WithUserAgent(userAgent string) func(*Config) {
	return func(s *Config) {
		s.UserAgent = userAgent
	}
}

func WithAutoTokenRefresh(autoTokenRefresh bool) func(*Config) {
	return func(s *Config) {
		s.AutoTokenRefresh = autoTokenRefresh
	}
}

func WithSilentMode(silentMode bool) func(*Config) {
	return func(s *Config) {
		s.SilentMode = silentMode
	}
}

func (c *InfisicalClient) setAccessToken(tokenDetails MachineIdentityCredential, credential interface{}, authMethod util.AuthMethod) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.tokenDetails = tokenDetails
	c.lastFetchedTime = time.Now()

	if c.authMethod != authMethod || c.firstFetchedTime.IsZero() {
		c.firstFetchedTime = time.Now()
		c.authMethod = authMethod
	}

	c.credential = credential
	c.httpClient.SetAuthScheme("Bearer")
	c.httpClient.SetAuthToken(c.tokenDetails.AccessToken)
}

func NewInfisicalClient(config *Config) InfisicalClientInterface {
	client := &InfisicalClient{}

	if config == nil {
		config = NewInfisicalClientConfig()
	}

	client.UpdateConfiguration(config) // set httpClient and config

	// add interfaces here
	client.secrets = &Secrets{client: client}
	client.folders = &Folders{client: client}
	client.auth = &Auth{client: client}

	if config.AutoTokenRefresh {
		go client.handleTokenLifeCycle()
	}

	return client
}

func (c *InfisicalClient) UpdateConfiguration(config *Config) {
	c.mu.Lock()
	defer c.mu.Unlock()

	config.SiteUrl = util.AppendAPIEndpoint(config.SiteUrl)
	c.config = *config

	if c.httpClient == nil {
		c.httpClient = resty.New().
			SetHeader("User-Agent", config.UserAgent).
			SetBaseURL(config.SiteUrl)
	} else {
		c.httpClient.
			SetHeader("User-Agent", config.UserAgent).
			SetBaseURL(config.SiteUrl)
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

func (c *InfisicalClient) handleTokenLifeCycle() {
	var warningPrinted = false
	authStrategies := map[util.AuthMethod]func(cred interface{}) (credential MachineIdentityCredential, err error){
		util.UNIVERSAL_AUTH: func(cred interface{}) (credential MachineIdentityCredential, err error) {
			if parsedCreds, ok := cred.(models.UniversalAuthCredential); ok {
				return c.auth.UniversalAuthLogin(parsedCreds.ClientID, parsedCreds.ClientSecret)
			}
			return MachineIdentityCredential{}, fmt.Errorf("failed to parse UniversalAuthCredential")
		},
		util.KUBERNETES: func(cred interface{}) (credential MachineIdentityCredential, err error) {
			if parsedCreds, ok := cred.(models.KubernetesCredential); ok {
				return c.auth.KubernetesRawServiceAccountTokenLogin(parsedCreds.IdentityID, parsedCreds.ServiceAccountToken)
			}
			return MachineIdentityCredential{}, fmt.Errorf("failed to parse KubernetesAuthCredential")
		},
		util.AZURE: func(cred interface{}) (credential MachineIdentityCredential, err error) {
			if parsedCreds, ok := cred.(models.AzureCredential); ok {
				return c.auth.AzureAuthLogin(parsedCreds.IdentityID, parsedCreds.Resource)
			}
			return MachineIdentityCredential{}, fmt.Errorf("failed to parse AzureAuthCredential")
		},
		util.GCP_ID_TOKEN: func(cred interface{}) (credential MachineIdentityCredential, err error) {
			if parsedCreds, ok := cred.(models.GCPIDTokenCredential); ok {
				return c.auth.GcpIdTokenAuthLogin(parsedCreds.IdentityID)
			}

			return MachineIdentityCredential{}, fmt.Errorf("failed to parse GCPIDTokenCredential")
		},
		util.GCP_IAM: func(cred interface{}) (credential MachineIdentityCredential, err error) {
			if parsedCreds, ok := cred.(models.GCPIAMCredential); ok {
				return c.auth.GcpIamAuthLogin(parsedCreds.IdentityID, parsedCreds.ServiceAccountKeyFilePath)
			}
			return MachineIdentityCredential{}, fmt.Errorf("failed to parse GCPIAMCredential")
		},
		util.AWS_IAM: func(cred interface{}) (credential MachineIdentityCredential, err error) {
			if parsedCreds, ok := cred.(models.AWSIAMCredential); ok {
				return c.auth.AwsIamAuthLogin(parsedCreds.IdentityID)
			}

			return MachineIdentityCredential{}, fmt.Errorf("failed to parse AWSIAMCredential")
		},
	}

	for {

		c.mu.RLock()
		config := c.config
		authMethod := c.authMethod
		tokenDetails := c.tokenDetails
		clientCredential := c.credential
		c.mu.RUnlock()

		if config.AutoTokenRefresh && authMethod != "" {

			if !config.SilentMode && !warningPrinted && tokenDetails.AccessTokenMaxTTL != 0 && tokenDetails.ExpiresIn != 0 {
				if tokenDetails.AccessTokenMaxTTL < 60 || tokenDetails.ExpiresIn < 60 {
					util.PrintWarning("Machine Identity access token TTL or max TTL is less than 60 seconds. This may cause excessive API calls, and you may be subject to rate-limits.")
				}
				warningPrinted = true
			}

			c.mu.RLock()

			timeNow := time.Now()
			timeSinceLastFetchSeconds := timeNow.Sub(c.lastFetchedTime).Seconds()
			timeSinceFirstFetchSeconds := timeNow.Sub(c.firstFetchedTime).Seconds()
			c.mu.RUnlock()

			if timeSinceFirstFetchSeconds >= float64(tokenDetails.AccessTokenMaxTTL-10) {
				newToken, err := authStrategies[c.authMethod](clientCredential)

				if err != nil && !config.SilentMode {
					util.PrintWarning(fmt.Sprintf("Failed to re-authenticate: %s\n", err.Error()))
				} else {
					c.setAccessToken(newToken, c.credential, c.authMethod)
					// fmt.Println("Access token successfully re-authenticated\n")
					c.mu.Lock()
					c.firstFetchedTime = time.Now()
					c.mu.Unlock()
				}

			} else if timeSinceLastFetchSeconds >= float64(tokenDetails.ExpiresIn-5) {
				// fmt.Printf("Access token expired, renewing...\n")

				renewedCredential, err := api.CallRenewAccessToken(c.httpClient, api.RenewAccessTokenRequest{AccessToken: tokenDetails.AccessToken})

				if err != nil {
					if !config.SilentMode {
						util.PrintWarning(fmt.Sprintf("Failed to renew access token: %s", err.Error()))
					}
				} else {
					// fmt.Println("Access token successfully renewed\n")
					c.setAccessToken(renewedCredential, clientCredential, authMethod)
				}
			}

			c.mu.RLock()
			nextAccessTokenExpiresInTime := c.lastFetchedTime.Add(time.Duration(tokenDetails.ExpiresIn*int64(time.Second)) - (5 * time.Second))
			accessTokenMaxTTLExpiresInTime := c.firstFetchedTime.Add(time.Duration(tokenDetails.AccessTokenMaxTTL*int64(time.Second)) - (5 * time.Second))
			expiresIn := time.Duration(c.tokenDetails.ExpiresIn * int64(time.Second))
			c.mu.RUnlock()

			if nextAccessTokenExpiresInTime.After(accessTokenMaxTTLExpiresInTime) {
				time.Sleep(expiresIn - nextAccessTokenExpiresInTime.Sub(accessTokenMaxTTLExpiresInTime))
			} else {
				time.Sleep(expiresIn - (5 * time.Second))
			}
		} else {
			time.Sleep(5 * time.Second)
		}

	}
}
