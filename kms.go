package infisical

import (
	"encoding/base64"

	api "github.com/infisical/go-sdk/packages/api/kms"
)

type KmsEncryptDataOptions = api.KmsEncryptDataV1Request
type KmsDecryptDataOptions = api.KmsDecryptDataV1Request

type KmsInterface interface {
	EncryptData(options KmsEncryptDataOptions) (string, error)
	DecryptData(options KmsDecryptDataOptions) (string, error)
}

type Kms struct {
	client *InfisicalClient
}

func (f *Kms) EncryptData(options KmsEncryptDataOptions) (string, error) {
	options.Plaintext = base64.StdEncoding.EncodeToString([]byte(options.Plaintext))
	res, err := api.CallKmsEncryptDataV1(f.client.httpClient, options)

	if err != nil {
		return "", err
	}

	return res.Ciphertext, nil
}

func (f *Kms) DecryptData(options KmsDecryptDataOptions) (string, error) {
	res, err := api.CallKmsDecryptDataV1(f.client.httpClient, options)

	if err != nil {
		return "", err
	}

	decodedPlaintext, err := base64.StdEncoding.DecodeString(res.Plaintext)
	if err != nil {
		return "", err
	}

	return string(decodedPlaintext), nil
}

func NewKms(client *InfisicalClient) KmsInterface {
	return &Kms{client: client}
}
