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
		accessToken = c.config.Auth.AccessToken
	case util.UNIVERSAL_AUTH:
		accessToken, err = api.CallUniversalAuthLogin(c.httpClient, api.UniversalAuthLoginRequest{
			ClientID:     c.config.Auth.UniversalAuth.ClientID,
			ClientSecret: c.config.Auth.UniversalAuth.ClientSecret,
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
