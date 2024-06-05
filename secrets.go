package infisical

import api "github.com/infisical/go-sdk/packages/api/secrets"

type ListSecretOptions = api.ListSecretsOptions

type SecretsInterface interface {
	List(options ListSecretOptions) string
}

type Secrets struct {
	client *Client
}

func (s *Secrets) List(options ListSecretOptions) string {
	api.CallListSecretsV3(s.client.httpClient, options)

	return "ListSecrets"
}

func NewSecrets(client *Client) SecretsInterface {
	return &Secrets{client: client}
}
