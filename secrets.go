package infisical

import (
	"os"

	api "github.com/infisical/go-sdk/packages/api/secrets"
	"github.com/infisical/go-sdk/packages/models"
	"github.com/infisical/go-sdk/packages/util"
)

type ListSecretsOptions = api.ListSecretsV3RawRequest
type RetrieveSecretOptions = api.RetrieveSecretV3RawRequest
type UpdateSecretOptions = api.UpdateSecretV3RawRequest
type CreateSecretOptions = api.CreateSecretV3RawRequest
type DeleteSecretOptions = api.DeleteSecretV3RawRequest

type SecretsInterface interface {
	List(options ListSecretsOptions) ([]models.Secret, error)
	Retrieve(options RetrieveSecretOptions) (models.Secret, error)
	Update(options UpdateSecretOptions) (models.Secret, error)
	Create(options CreateSecretOptions) (models.Secret, error)
	Delete(options DeleteSecretOptions) (models.Secret, error)
}

type Secrets struct {
	client *InfisicalClient
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

	return util.SortSecretsByKeys(secrets), nil
}

func (s *Secrets) Retrieve(options RetrieveSecretOptions) (models.Secret, error) {
	res, err := api.CallRetrieveSecretV3(s.client.httpClient, options)

	if err != nil {
		return models.Secret{}, err
	}

	return res.Secret, nil
}

func (s *Secrets) Update(options UpdateSecretOptions) (models.Secret, error) {
	res, err := api.CallUpdateSecretV3(s.client.httpClient, options)

	if err != nil {
		return models.Secret{}, err
	}

	return res.Secret, nil
}

func (s *Secrets) Create(options CreateSecretOptions) (models.Secret, error) {
	res, err := api.CallCreateSecretV3(s.client.httpClient, options)

	if err != nil {
		return models.Secret{}, err
	}

	return res.Secret, nil
}

func (s *Secrets) Delete(options DeleteSecretOptions) (models.Secret, error) {
	res, err := api.CallDeleteSecretV3(s.client.httpClient, options)

	if err != nil {
		return models.Secret{}, err
	}

	return res.Secret, nil
}

func NewSecrets(client *InfisicalClient) SecretsInterface {
	return &Secrets{client: client}
}
