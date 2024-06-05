package infisical

import (
	"fmt"

	api "github.com/infisical/go-sdk/packages/api/auth"
	"github.com/infisical/go-sdk/packages/util"
)

func (c *Client) authenticateHttpClient() error {
	var err error
	var accessToken string

	switch c.authMethod {
	case util.ACCESS_TOKEN:
		accessToken = c.config.Auth.accessToken
	case util.UNIVERSAL_AUTH:
		accessToken, err = api.CallUniversalAuthLogin(c.httpClient, api.UniversalAuthLoginRequest{
			ClientId:     c.config.Auth.universalAuth.clientId,
			ClientSecret: c.config.Auth.universalAuth.clientSecret,
		})
	}

	if err != nil {
		return err
	}

	if accessToken == "" {
		return fmt.Errorf("no access token obtained")
	}

	c.httpClient.SetAuthScheme("Bearer") // For now all our auth methods are Bearer based, but this could potentially change in the future.
	c.httpClient.SetAuthToken(accessToken)

	return nil
}
