package infisical

import (
	api "github.com/infisical/go-sdk/packages/api/auth"
	"github.com/infisical/go-sdk/packages/errors"
	"github.com/infisical/go-sdk/packages/models"
)

type MachineIdentityCredential = api.MachineIdentityAuthLoginResponse

type Secret = models.Secret
type SecretImport = models.SecretImport

type APIError = errors.APIError
type RequestError = errors.RequestError
