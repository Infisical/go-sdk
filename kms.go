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
	ListKeys(options api.ListKmsKeysV1Request) ([]api.KmsKey, error)
	CreateKey(options api.CreateKmsKeyV1Request) (*api.KmsKey, error)
	UpdateKey(options api.UpdateKmsKeyV1Request) (*api.KmsKey, error)
	DeleteKey(options api.DeleteKmsKeyV1Request) (*api.KmsKey, error)
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

func (f *Kms) ListKeys(options api.ListKmsKeysV1Request) ([]api.KmsKey, error) {
	res, err := api.CallListKmsKeysV1(f.client.httpClient, options)
	if err != nil {
		return nil, err
	}
	return res.Keys, nil
}

func (f *Kms) CreateKey(options api.CreateKmsKeyV1Request) (*api.KmsKey, error) {
	res, err := api.CallCreateKmsKeyV1(f.client.httpClient, options)
	if err != nil {
		return nil, err
	}
	return &res.Key, nil
}

func (f *Kms) UpdateKey(options api.UpdateKmsKeyV1Request) (*api.KmsKey, error) {
	res, err := api.CallUpdateKmsKeyV1(f.client.httpClient, options)
	if err != nil {
		return nil, err
	}
	return &res.Key, nil
}

func (f *Kms) DeleteKey(options api.DeleteKmsKeyV1Request) (*api.KmsKey, error) {
	res, err := api.CallDeleteKmsKeyV1(f.client.httpClient, options)
	if err != nil {
		return nil, err
	}
	return &res.Key, nil
}

func NewKms(client *InfisicalClient) KmsInterface {
	return &Kms{client: client}
}
