package infisical

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	api "github.com/infisical/go-sdk/packages/api/auth"
	"github.com/infisical/go-sdk/packages/util"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
)

type KubernetesAuthLoginOptions struct {
	IdentityID              string
	ServiceAccountTokenPath string
}

// func epochTime() time.Time { return time.Unix(0, 0) }

type AuthInterface interface {
	SetAccessToken(accessToken string)
	UniversalAuthLogin(clientID string, clientSecret string) (credential MachineIdentityCredential, err error)
	KubernetesAuthLogin(identityID string, serviceAccountTokenPath string) (credential MachineIdentityCredential, err error)
	KubernetesRawServiceAccountTokenLogin(identityID string, serviceAccountToken string) (credential MachineIdentityCredential, err error)
	AzureAuthLogin(identityID string) (credential MachineIdentityCredential, err error)
	GcpIdTokenAuthLogin(identityID string) (credential MachineIdentityCredential, err error)
	GcpIamAuthLogin(identityID string, serviceAccountKeyFilePath string) (credential MachineIdentityCredential, err error)
	AwsIamAuthLogin(identityId string) (credential MachineIdentityCredential, err error)
	// Renew access token. This is useful when you want to refresh the access token before it expires.
	TokenRenew(accessToken string) (credential MachineIdentityCredential, err error)
}

type Auth struct {
	client *InfisicalClient
}

func (a *Auth) SetAccessToken(accessToken string) {
	a.client.setAccessToken(accessToken, util.ACCESS_TOKEN)
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

	a.client.setAccessToken(credential.AccessToken, util.UNIVERSAL_AUTH)
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

	a.client.setAccessToken(credential.AccessToken, util.KUBERNETES)
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

	a.client.setAccessToken(credential.AccessToken, util.KUBERNETES)
	return credential, nil
}

func (a *Auth) AzureAuthLogin(identityID string) (credential MachineIdentityCredential, err error) {
	if identityID == "" {
		identityID = os.Getenv(util.INFISICAL_AZURE_AUTH_IDENTITY_ID_ENV_NAME)
	}

	jwt, jwtError := util.GetAzureMetadataToken(a.client.httpClient)

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

	a.client.setAccessToken(credential.AccessToken, util.AZURE)
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

	a.client.setAccessToken(credential.AccessToken, util.GCP_ID_TOKEN)
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

	a.client.setAccessToken(credential.AccessToken, util.GCP_IAM)
	return credential, nil
}

func (a *Auth) AwsIamAuthLogin(identityId string) (credential MachineIdentityCredential, err error) {

	if identityId == "" {
		identityId = os.Getenv(util.INFISICAL_AWS_IAM_AUTH_IDENTITY_ID_ENV_NAME)
	}

	awsRegion, err := util.GetAwsRegion()
	if err != nil {
		return MachineIdentityCredential{}, err
	}

	awsCfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if err != nil {
		return MachineIdentityCredential{}, fmt.Errorf("unable to load SDK config, %v", err)
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

	credentials, err := awsCfg.Credentials.Retrieve(context.TODO())
	if err != nil {
		return MachineIdentityCredential{}, fmt.Errorf("error retrieving credentials: %v", err)
	}

	hashGenerator := sha256.New()
	hashGenerator.Write([]byte(iamRequestBody))
	payloadHash := fmt.Sprintf("%x", hashGenerator.Sum(nil))

	signer := v4.NewSigner()
	err = signer.SignHTTP(context.TODO(), credentials, req, payloadHash, "sts", awsRegion, time.Now())

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

	a.client.setAccessToken(credential.AccessToken, util.AWS_IAM)
	return credential, nil
}

func (a *Auth) TokenRenew(accessToken string) (credential MachineIdentityCredential, err error) {
	credential, err = api.CallTokenRenew(a.client.httpClient, api.TokenRenewRequest{
		AccessToken: accessToken,
	})

	if err != nil {
		return MachineIdentityCredential{}, err
	}

	a.client.setAccessToken(credential.AccessToken, util.UNIVERSAL_AUTH)
	return credential, nil

}

func NewAuth(client *InfisicalClient) AuthInterface {
	return &Auth{client: client}
}
