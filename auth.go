package infisical

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	api "github.com/infisical/go-sdk/packages/api/auth"
	"github.com/infisical/go-sdk/packages/models"
	"github.com/infisical/go-sdk/packages/util"
	"github.com/oracle/oci-go-sdk/v65/common"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
)

type KubernetesAuthLoginOptions struct {
	IdentityID              string
	ServiceAccountTokenPath string
}

type AuthInterface interface {
	SetAccessToken(accessToken string)
	GetAccessToken() string
	UniversalAuthLogin(clientID string, clientSecret string) (credential MachineIdentityCredential, err error)
	JwtAuthLogin(identityID string, jwt string) (credential MachineIdentityCredential, err error)
	KubernetesAuthLogin(identityID string, serviceAccountTokenPath string) (credential MachineIdentityCredential, err error)
	KubernetesRawServiceAccountTokenLogin(identityID string, serviceAccountToken string) (credential MachineIdentityCredential, err error)
	AzureAuthLogin(identityID string, resource string) (credential MachineIdentityCredential, err error)
	GcpIdTokenAuthLogin(identityID string) (credential MachineIdentityCredential, err error)
	GcpIamAuthLogin(identityID string, serviceAccountKeyFilePath string) (credential MachineIdentityCredential, err error)
	AwsIamAuthLogin(identityId string) (credential MachineIdentityCredential, err error)
	OidcAuthLogin(identityId string, jwt string) (credential MachineIdentityCredential, err error)
	OciAuthLogin(options OciAuthLoginOptions) (credential MachineIdentityCredential, err error)
	LDAPAuthLogin(identityID string, username string, password string) (credential MachineIdentityCredential, err error)
	RevokeAccessToken() error
}

type Auth struct {
	client *InfisicalClient
}

func (a *Auth) SetAccessToken(accessToken string) {
	a.client.setPlainAccessToken(accessToken)
}

func (a *Auth) GetAccessToken() string {
	// case: user has set an access token manually, so we get it directly from the credential
	if a.client.authMethod == util.ACCESS_TOKEN {
		if parsedCreds, ok := a.client.credential.(models.AccessTokenCredential); ok {
			return parsedCreds.AccessToken
		}
		return ""
	}
	return a.client.tokenDetails.AccessToken
}

func (a *Auth) RevokeAccessToken() error {
	if a.client.tokenDetails.AccessToken == "" {
		return errors.New("sdk client is not authenticated, cannot revoke access token")
	}

	_, err := api.CallRevokeAccessToken(a.client.httpClient, api.RevokeAccessTokenRequest{
		AccessToken: a.client.tokenDetails.AccessToken,
	})

	if err != nil {
		return err
	}
	a.client.clearAccessToken()

	return nil
}
func (a *Auth) UniversalAuthLogin(clientID string, clientSecret string) (credential MachineIdentityCredential, err error) {

	if clientID == "" {
		clientID = os.Getenv(util.INFISICAL_UNIVERSAL_AUTH_CLIENT_ID_ENV_NAME)
	}
	if clientSecret == "" {
		clientSecret = os.Getenv(util.INFISICAL_UNIVERSAL_AUTH_CLIENT_SECRET_ENV_NAME)
	}

	credential, err = api.CallUniversalAuthLogin(a.client.httpClient, api.UniversalAuthLoginRequest{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	})

	if err != nil {
		return MachineIdentityCredential{}, err
	}

	a.client.setAccessToken(
		credential,
		models.UniversalAuthCredential{ClientID: clientID, ClientSecret: clientSecret},
		util.UNIVERSAL_AUTH,
	)
	return credential, nil

}

func (a *Auth) KubernetesAuthLogin(identityID string, serviceAccountTokenPath string) (credential MachineIdentityCredential, err error) {

	if serviceAccountTokenPath == "" {
		serviceAccountTokenPath = os.Getenv(util.DEFAULT_KUBERNETES_SERVICE_ACCOUNT_TOKEN_PATH)
	}
	if identityID == "" {
		identityID = os.Getenv(util.INFISICAL_KUBERNETES_IDENTITY_ID_ENV_NAME)
	}

	serviceAccountToken, serviceAccountTokenErr := util.GetKubernetesServiceAccountToken(serviceAccountTokenPath)

	if serviceAccountTokenErr != nil {
		return MachineIdentityCredential{}, serviceAccountTokenErr
	}

	credential, err = api.CallKubernetesAuthLogin(a.client.httpClient, api.KubernetesAuthLoginRequest{
		IdentityID: identityID,
		JWT:        serviceAccountToken,
	})

	if err != nil {
		return MachineIdentityCredential{}, err
	}

	a.client.setAccessToken(
		credential,
		models.KubernetesCredential{IdentityID: identityID, ServiceAccountToken: serviceAccountToken},
		util.KUBERNETES,
	)

	return credential, nil

}

func (a *Auth) KubernetesRawServiceAccountTokenLogin(identityID string, serviceAccountToken string) (credential MachineIdentityCredential, err error) {

	if identityID == "" {
		identityID = os.Getenv(util.INFISICAL_KUBERNETES_IDENTITY_ID_ENV_NAME)
	}

	credential, err = api.CallKubernetesAuthLogin(a.client.httpClient, api.KubernetesAuthLoginRequest{
		IdentityID: identityID,
		JWT:        serviceAccountToken,
	})

	if err != nil {
		return MachineIdentityCredential{}, err
	}

	a.client.setAccessToken(
		credential,
		models.KubernetesCredential{IdentityID: identityID, ServiceAccountToken: serviceAccountToken},
		util.KUBERNETES,
	)
	return credential, nil
}

func (a *Auth) AzureAuthLogin(identityID string, resource string) (credential MachineIdentityCredential, err error) {
	if identityID == "" {
		identityID = os.Getenv(util.INFISICAL_AZURE_AUTH_IDENTITY_ID_ENV_NAME)
	}

	jwt, jwtError := util.GetAzureMetadataToken(a.client.httpClient, resource)

	if jwtError != nil {
		return MachineIdentityCredential{}, jwtError
	}

	credential, err = api.CallAzureAuthLogin(a.client.httpClient, api.AzureAuthLoginRequest{
		IdentityID: identityID,
		JWT:        jwt,
	})

	if err != nil {
		return MachineIdentityCredential{}, err
	}

	a.client.setAccessToken(
		credential,
		models.AzureCredential{IdentityID: identityID, Resource: resource},
		util.AZURE,
	)
	return credential, nil
}

func (a *Auth) GcpIdTokenAuthLogin(identityID string) (credential MachineIdentityCredential, err error) {
	if identityID == "" {
		identityID = os.Getenv(util.INFISICAL_GCP_AUTH_IDENTITY_ID_ENV_NAME)
	}

	jwt, jwtError := util.GetGCPMetadataToken(a.client.httpClient, identityID)

	if jwtError != nil {
		return MachineIdentityCredential{}, jwtError
	}

	credential, err = api.CallGCPAuthLogin(a.client.httpClient, api.GCPAuthLoginRequest{
		IdentityID: identityID,
		JWT:        jwt,
	})

	if err != nil {
		return MachineIdentityCredential{}, err
	}

	a.client.setAccessToken(
		credential,
		models.GCPIDTokenCredential{IdentityID: identityID},
		util.GCP_ID_TOKEN,
	)
	return credential, nil
}

func (a *Auth) GcpIamAuthLogin(identityID string, serviceAccountKeyFilePath string) (credential MachineIdentityCredential, err error) {
	if identityID == "" {
		identityID = os.Getenv(util.INFISICAL_GCP_AUTH_IDENTITY_ID_ENV_NAME)
	}
	if serviceAccountKeyFilePath == "" {
		serviceAccountKeyFilePath = os.Getenv(util.INFISICAL_GCP_IAM_SERVICE_ACCOUNT_KEY_FILE_PATH_ENV_NAME)
	}

	jwt, jwtError := util.GetGCPIamServiceAccountToken(identityID, serviceAccountKeyFilePath)

	if jwtError != nil {
		return MachineIdentityCredential{}, jwtError
	}

	credential, err = api.CallGCPAuthLogin(a.client.httpClient, api.GCPAuthLoginRequest{
		IdentityID: identityID,
		JWT:        jwt,
	})

	if err != nil {
		return MachineIdentityCredential{}, err
	}

	a.client.setAccessToken(
		credential,
		models.GCPIAMCredential{IdentityID: identityID, ServiceAccountKeyFilePath: serviceAccountKeyFilePath},
		util.GCP_IAM,
	)
	return credential, nil
}

func (a *Auth) AwsIamAuthLogin(identityId string) (credential MachineIdentityCredential, err error) {

	if identityId == "" {
		identityId = os.Getenv(util.INFISICAL_AWS_IAM_AUTH_IDENTITY_ID_ENV_NAME)
	}

	awsCredentials, awsRegion, err := util.RetrieveAwsCredentials()
	if err != nil {
		return MachineIdentityCredential{}, err
	}

	// Prepare request for signing
	iamRequestURL := fmt.Sprintf("https://sts.%s.amazonaws.com/", awsRegion)
	iamRequestBody := "Action=GetCallerIdentity&Version=2011-06-15"

	req, err := http.NewRequest(http.MethodPost, iamRequestURL, strings.NewReader(iamRequestBody))
	if err != nil {
		return MachineIdentityCredential{}, fmt.Errorf("error creating HTTP request: %v", err)
	}

	currentTime := time.Now().UTC()
	req.Header.Add("X-Amz-Date", currentTime.Format("20060102T150405Z"))

	hashGenerator := sha256.New()
	hashGenerator.Write([]byte(iamRequestBody))
	payloadHash := fmt.Sprintf("%x", hashGenerator.Sum(nil))

	signer := v4.NewSigner()
	err = signer.SignHTTP(context.TODO(), awsCredentials, req, payloadHash, "sts", awsRegion, time.Now())

	if err != nil {
		return MachineIdentityCredential{}, fmt.Errorf("error signing request: %v", err)
	}

	var realHeaders map[string]string = make(map[string]string)
	for name, values := range req.Header {
		if strings.ToLower(name) == "content-length" {
			continue
		}
		realHeaders[name] = values[0]
	}
	realHeaders["Host"] = fmt.Sprintf("sts.%s.amazonaws.com", awsRegion)
	realHeaders["Content-Type"] = "application/x-www-form-urlencoded; charset=utf-8"
	realHeaders["Content-Length"] = fmt.Sprintf("%d", len(iamRequestBody))

	// convert the headers to a json marshalled string
	jsonStringHeaders, err := json.Marshal(realHeaders)

	if err != nil {
		return MachineIdentityCredential{}, fmt.Errorf("error marshalling headers: %v", err)
	}

	credential, tokenErr := api.CallAWSIamAuthLogin(a.client.httpClient, api.AwsIamAuthLoginRequest{
		HTTPRequestMethod: req.Method,
		// Encoding is intended, we decode it on severside, and I know everything happening on the server is being done correctly. So it's something broken in this code somewhere.
		IamRequestBody:    base64.StdEncoding.EncodeToString([]byte(iamRequestBody)),
		IamRequestHeaders: base64.StdEncoding.EncodeToString(jsonStringHeaders),
		IdentityId:        identityId,
	})

	if tokenErr != nil {
		return MachineIdentityCredential{}, tokenErr
	}

	a.client.setAccessToken(
		credential,
		models.AWSIAMCredential{IdentityID: identityId},
		util.AWS_IAM,
	)
	return credential, nil
}

func (a *Auth) OidcAuthLogin(identityId string, jwt string) (credential MachineIdentityCredential, err error) {
	if identityId == "" {
		identityId = os.Getenv(util.INFISICAL_OIDC_AUTH_IDENTITY_ID_ENV_NAME)
	}

	credential, err = api.CallOidcAuthLogin(a.client.httpClient, api.OidcAuthLoginRequest{
		IdentityID: identityId,
		JWT:        jwt,
	})

	if err != nil {
		return MachineIdentityCredential{}, err
	}

	a.client.setAccessToken(
		credential,
		models.OIDCCredential{IdentityID: identityId},
		util.OIDC_AUTH,
	)
	return credential, nil

}

func (a *Auth) JwtAuthLogin(identityID string, jwt string) (credential MachineIdentityCredential, err error) {
	credential, err = api.CallJwtAuthLogin(a.client.httpClient, api.JwtAuthLoginRequest{
		IdentityID: identityID,
		JWT:        jwt,
	})

	if err != nil {
		return MachineIdentityCredential{}, err
	}

	a.client.setAccessToken(
		credential,
		models.JWTCredential{IdentityID: identityID, JWT: jwt},
		util.JWT_AUTH,
	)
	return credential, nil
}

func (a *Auth) OciAuthLogin(options OciAuthLoginOptions) (credential MachineIdentityCredential, err error) {

	if options.IdentityID == "" {
		options.IdentityID = os.Getenv(util.INFISICAL_OCI_AUTH_IDENTITY_ID_ENV_NAME)
	}

	provider := common.NewRawConfigurationProvider(
		options.TenancyID,
		options.UserID,
		options.Region,
		options.Fingerprint,
		options.PrivateKey,
		options.Passphrase,
	)

	requestURL := fmt.Sprintf("https://identity.%s.oraclecloud.com/20160918/users/%s", options.Region, options.UserID)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return MachineIdentityCredential{}, fmt.Errorf("OciAuthLogin: failed to create request: %w", err)
	}

	req.Header.Set("host", fmt.Sprintf("identity.%s.oraclecloud.com", options.Region))
	req.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))

	signer := common.DefaultRequestSigner(provider)
	err = signer.Sign(req)
	if err != nil {
		return MachineIdentityCredential{}, fmt.Errorf("OciAuthLogin: failed to sign request: %w", err)
	}

	headersMap := make(map[string]string)
	for name, values := range req.Header {
		if len(values) > 0 {
			// Convert header names to lowercase to match OCI signature expectations
			lowerName := strings.ToLower(name)
			headersMap[lowerName] = values[0]
		}
	}

	return api.CallOciAuthLogin(a.client.httpClient, api.OciAuthLoginRequest{
		IdentityID: options.IdentityID,
		UserOcid:   options.UserID,
		Headers:    headersMap,
	})
}

func (a *Auth) LDAPAuthLogin(identityID string, username string, password string) (credential MachineIdentityCredential, err error) {
	if identityID == "" {
		identityID = os.Getenv(util.INFISICAL_LDAP_AUTH_IDENTITY_ID_ENV_NAME)
	}

	return api.CallLdapAuthLogin(a.client.httpClient, api.LdapAuthLoginRequest{
		IdentityID: identityID,
		Username:   username,
		Password:   password,
	})
}

func NewAuth(client *InfisicalClient) AuthInterface {
	return &Auth{client: client}
}
