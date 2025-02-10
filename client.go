package infisical

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	api "github.com/levidurfee/go-sdk/packages/api/auth"
	"github.com/levidurfee/go-sdk/packages/models"
	"github.com/levidurfee/go-sdk/packages/util"
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

	secrets        SecretsInterface
	folders        FoldersInterface
	auth           AuthInterface
	dynamicSecrets DynamicSecretsInterface
	kms            KmsInterface
	ssh            SshInterface
}

type InfisicalClientInterface interface {
	UpdateConfiguration(config Config)
	Secrets() SecretsInterface
	Folders() FoldersInterface
	Auth() AuthInterface
	DynamicSecrets() DynamicSecretsInterface
	Kms() KmsInterface
	Ssh() SshInterface
}

type Config struct {
	SiteUrl          string `default:"https://app.infisical.com"`
	CaCertificate    string
	UserAgent        string `default:"infisical-go-sdk"` // User-Agent header to be used on requests sent by the SDK. Defaults to `infisical-go-sdk`. Do not modify this unless you have a reason to do so.
	AutoTokenRefresh bool   `default:"true"`             // Wether or not to automatically refresh the auth token after using one of the .Auth() methods. Defaults to `true`.
	SilentMode       bool   `default:"false"`            // If enabled, the SDK will not print any warnings to the console.
}

func setDefaults(cfg *Config) {
	t := reflect.TypeOf(*cfg) // we need to dereference the pointer to get the struct type
	v := reflect.ValueOf(cfg).Elem()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		defaultVal := field.Tag.Get("default")
		if defaultVal == "" {
			continue
		}

		switch field.Type.Kind() {
		case reflect.Int:
			if v.Field(i).Int() == 0 {
				val, _ := strconv.Atoi(defaultVal)
				v.Field(i).SetInt(int64(val))
			}
		case reflect.String:
			if v.Field(i).String() == "" {
				v.Field(i).SetString(defaultVal)
			}
		case reflect.Bool:
			if !v.Field(i).Bool() {
				val, _ := strconv.ParseBool(defaultVal)
				v.Field(i).SetBool(val)
			}
		}
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

func (c *InfisicalClient) setPlainAccessToken(accessToken string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.authMethod = util.ACCESS_TOKEN
	c.httpClient.SetAuthScheme("Bearer")
	c.httpClient.SetAuthToken(accessToken)
}

func NewInfisicalClient(context context.Context, config Config) InfisicalClientInterface {
	client := &InfisicalClient{}
	setDefaults(&config)
	client.UpdateConfiguration(config) // set httpClient and config

	// add interfaces here
	client.secrets = NewSecrets(client)
	client.folders = NewFolders(client)
	client.auth = NewAuth(client)
	client.dynamicSecrets = NewDynamicSecrets(client)
	client.kms = NewKms(client)
	client.ssh = NewSsh(client)

	if config.AutoTokenRefresh {
		go client.handleTokenLifeCycle(context)
	}

	return client
}

func (c *InfisicalClient) UpdateConfiguration(config Config) {
	c.mu.Lock()
	defer c.mu.Unlock()

	setDefaults(&config)
	config.SiteUrl = util.AppendAPIEndpoint(config.SiteUrl)
	c.config = config

	if c.httpClient == nil {
		c.httpClient = resty.New().
			SetHeader("User-Agent", config.UserAgent).
			SetBaseURL(config.SiteUrl)

	} else {
		c.httpClient.
			SetHeader("User-Agent", config.UserAgent).
			SetBaseURL(config.SiteUrl)
	}

	if config.CaCertificate != "" {
		caCertPool, err := x509.SystemCertPool()
		if err != nil && !config.SilentMode {
			util.PrintWarning(fmt.Sprintf("failed to load system root CA pool: %v", err))
		}

		if ok := caCertPool.AppendCertsFromPEM([]byte(config.CaCertificate)); !ok && !config.SilentMode {
			util.PrintWarning("failed to append CA certificate")
		}

		tlsConfig := &tls.Config{
			RootCAs: caCertPool,
		}

		c.httpClient.SetTLSClientConfig(tlsConfig)
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

func (c *InfisicalClient) DynamicSecrets() DynamicSecretsInterface {
	return c.dynamicSecrets
}

func (c *InfisicalClient) Kms() KmsInterface {
	return c.kms
}

func (c *InfisicalClient) Ssh() SshInterface {
	return c.ssh
}

func (c *InfisicalClient) handleTokenLifeCycle(context context.Context) {
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

	const RE_AUTHENTICATION_INTERVAL_BUFFER = 2
	const RENEWAL_INTERVAL_BUFFER = 5

	for {
		select {
		case <-context.Done():
			return // The context has been cancelled, clean up and return from the loop to stop the goroutine
		default:
			{

				c.mu.RLock()
				config := c.config
				authMethod := c.authMethod
				tokenDetails := c.tokenDetails
				clientCredential := c.credential
				c.mu.RUnlock()

				if config.AutoTokenRefresh && authMethod != "" && authMethod != util.ACCESS_TOKEN {

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

					if timeSinceFirstFetchSeconds >= float64(tokenDetails.AccessTokenMaxTTL-RE_AUTHENTICATION_INTERVAL_BUFFER) {
						newToken, err := authStrategies[c.authMethod](clientCredential)

						if err != nil && !config.SilentMode {
							util.PrintWarning(fmt.Sprintf("Failed to re-authenticate: %s\n", err.Error()))
						} else {
							c.setAccessToken(newToken, c.credential, c.authMethod)
							c.mu.Lock()
							c.firstFetchedTime = time.Now()
							c.mu.Unlock()
						}

					} else if timeSinceLastFetchSeconds >= float64(tokenDetails.ExpiresIn-RENEWAL_INTERVAL_BUFFER) {
						timeUntilMaxTTL := float64(tokenDetails.AccessTokenMaxTTL) - timeSinceFirstFetchSeconds

						// Case 1: The time until the max TTL is less than the time until the next access token expiry
						if timeUntilMaxTTL < float64(tokenDetails.ExpiresIn) {
							// If renewing would exceed max TTL, directly re-authenticate
							newToken, err := authStrategies[c.authMethod](clientCredential)
							if err != nil && !config.SilentMode {
								util.PrintWarning(fmt.Sprintf("Failed to re-authenticate: %s\n", err.Error()))
							} else {
								c.setAccessToken(newToken, c.credential, c.authMethod)
								c.mu.Lock()
								c.firstFetchedTime = time.Now()
								c.mu.Unlock()
							}
							// Case 2: The time until the max TTL is greater than the time until the next access token expiry
						} else {
							renewedCredential, err := api.CallRenewAccessToken(c.httpClient, api.RenewAccessTokenRequest{AccessToken: tokenDetails.AccessToken})

							if err != nil {
								if !config.SilentMode {
									util.PrintWarning(fmt.Sprintf("Failed to renew access token: %s\n\nAttempting to re-authenticate.", err.Error()))
								}

								newToken, err := authStrategies[c.authMethod](clientCredential)
								if err != nil && !config.SilentMode {
									util.PrintWarning(fmt.Sprintf("Failed to re-authenticate: %s\n", err.Error()))
								} else {
									c.setAccessToken(newToken, c.credential, c.authMethod)
									c.mu.Lock()
									c.firstFetchedTime = time.Now()
									c.mu.Unlock()
								}
							} else {
								c.setAccessToken(renewedCredential, clientCredential, authMethod)
							}
						}
					}

					c.mu.RLock()
					nextAccessTokenExpiresInTime := c.lastFetchedTime.Add(time.Duration(tokenDetails.ExpiresIn*int64(time.Second)) - (5 * time.Second))
					accessTokenMaxTTLExpiresInTime := c.firstFetchedTime.Add(time.Duration(tokenDetails.AccessTokenMaxTTL*int64(time.Second)) - (5 * time.Second))
					expiresIn := time.Duration(c.tokenDetails.ExpiresIn * int64(time.Second))
					c.mu.RUnlock()

					if nextAccessTokenExpiresInTime.After(accessTokenMaxTTLExpiresInTime) {
						// Calculate the sleep time
						sleepTime := expiresIn - nextAccessTokenExpiresInTime.Sub(accessTokenMaxTTLExpiresInTime)

						// Ensure we sleep for at least 1 second
						if sleepTime < 1*time.Second {
							sleepTime = time.Second * 1
						}

						if err := util.SleepWithContext(context, sleepTime); err != nil && (err == util.ErrContextCanceled || errors.Is(err, util.ErrContextDeadlineExceeded)) {
							return
						}

					} else {
						sleepTime := expiresIn - (5 * time.Second)

						if sleepTime < time.Second {
							sleepTime = time.Millisecond * 500
						}

						if err := util.SleepWithContext(context, sleepTime); err != nil && (err == util.ErrContextCanceled || errors.Is(err, util.ErrContextDeadlineExceeded)) {
							return
						}
					}
				} else {
					if err := util.SleepWithContext(context, 1*time.Second); err != nil && (err == util.ErrContextCanceled || errors.Is(err, util.ErrContextDeadlineExceeded)) {
						return
					}
				}
			}
		}
	}
}
