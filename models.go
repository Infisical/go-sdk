package infisical

import (
	"github.com/infisical/go-sdk/packages/errors"
	"github.com/infisical/go-sdk/packages/models"
)

type Secret = models.Secret
type SecretImport = models.SecretImport

type APIError = errors.APIError
type RequestError = errors.RequestError

func IsAPIError(err error) bool {
	return errors.IsAPIError(err)
}

func IsRequestError(err error) bool {
	return errors.IsRequestError(err)
}
