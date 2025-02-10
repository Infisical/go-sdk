package infisical

import (
	"encoding/base64"

	api "github.com/levidurfee/go-sdk/packages/api/kms"
)

type KmsEncryptDataOptions = api.KmsEncryptDataV1Request
type KmsDecryptDataOptions = api.KmsDecryptDataV1Request
type KmsListKeysOptions = api.KmsListKeysV1Request
type KmsCreateKeyOptions = api.KmsCreateKeyV1Request
type KmsUpdateKeyOptions = api.KmsUpdateKeyV1Request
type KmsDeleteKeyOptions = api.KmsDeleteKeyV1Request
type KmsKey = api.KmsKey

type KmsInterface interface {
	EncryptData(options KmsEncryptDataOptions) (string, error)
	DecryptData(options KmsDecryptDataOptions) (string, error)
	ListKeys(options KmsListKeysOptions) ([]KmsKey, error)
	CreateKey(options KmsCreateKeyOptions) (*KmsKey, error)
	UpdateKey(options KmsUpdateKeyOptions) (*KmsKey, error)
	DeleteKey(options KmsDeleteKeyOptions) (*KmsKey, error)
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

func (f *Kms) ListKeys(options KmsListKeysOptions) ([]KmsKey, error) {
	res, err := api.CallListKmsKeysV1(f.client.httpClient, options)
	if err != nil {
		return nil, err
	}
	return res.Keys, nil
}

func (f *Kms) CreateKey(options KmsCreateKeyOptions) (*KmsKey, error) {
	res, err := api.CallCreateKmsKeyV1(f.client.httpClient, options)
	if err != nil {
		return nil, err
	}
	return &res.Key, nil
}

func (f *Kms) UpdateKey(options KmsUpdateKeyOptions) (*KmsKey, error) {
	res, err := api.CallUpdateKmsKeyV1(f.client.httpClient, options)
	if err != nil {
		return nil, err
	}
	return &res.Key, nil
}

func (f *Kms) DeleteKey(options KmsDeleteKeyOptions) (*KmsKey, error) {
	res, err := api.CallDeleteKmsKeyV1(f.client.httpClient, options)
	if err != nil {
		return nil, err
	}
	return &res.Key, nil
}

func NewKms(client *InfisicalClient) KmsInterface {
	return &Kms{client: client}
}
