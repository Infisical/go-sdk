package infisical

import (
	"os"

	api "github.com/infisical/go-sdk/packages/api/secrets"
	"github.com/infisical/go-sdk/packages/models"
	"github.com/infisical/go-sdk/packages/util"
)

type ListSecretsOptions = api.ListSecretsRequest

type SecretsInterface interface {
	List(options ListSecretsOptions) ([]models.Secret, error)
}

type Secrets struct {
	client *Client
}

func (s *Secrets) List(options ListSecretsOptions) ([]models.Secret, error) {
	res, err := api.CallListSecretsV3(s.client.httpClient, options)

	if err != nil {
		return nil, err
	}

	if options.Recursive {
		util.EnsureUniqueSecretsByKey(&res.Secrets)
	}

	secrets := append([]models.Secret(nil), res.Secrets...) // Clone main secrets slice, we will modify this if imports are enabled
	if options.IncludeImports {

		// Append secrets from imports
		for _, importBlock := range res.Imports {
			for _, importSecret := range importBlock.Secrets {
				// Only append the secret if it is not already in the list, imports take precedence
				if !util.ContainsSecret(secrets, importSecret.SecretKey) {
					secrets = append(secrets, importSecret)
				}
			}
		}
	}

	if options.AttachToProcessEnv {
		for _, secret := range secrets {
			// Only set the environment variable if it is not already set
			if os.Getenv(secret.SecretKey) == "" {
				os.Setenv(secret.SecretKey, secret.SecretValue)
			}
		}

	}

	return secrets, nil
}

func NewSecrets(client *Client) SecretsInterface {
	return &Secrets{client: client}
}
