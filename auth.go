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
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type UniversalAuthLoginOptions = api.UniversalAuthLoginRequest

type KubernetesAuthLoginOptions struct {
	IdentityID              string
	ServiceAccountTokenPath string
}

// func epochTime() time.Time { return time.Unix(0, 0) }

type AuthInterface interface {
	SetAccessToken(accessToken string)
	UniversalAuthLogin(clientID string, clientSecret string) (accessToken string, err error)
	KubernetesAuthLogin(identityID string, serviceAccountTokenPath string) (accessToken string, err error)
	AzureAuthLogin(identityID string) (accessToken string, err error)
	GcpIdTokenAuthLogin(identityID string) (accessToken string, err error)
	GcpIamAuthLogin(identityID string, serviceAccountKeyFilePath string) (accessToken string, err error)
	AwsIamAuthLogin(identityId string) (accessToken string, err error)
}

type Auth struct {
	client *InfisicalClient
}

func (a *Auth) SetAccessToken(accessToken string) {
	a.client.setAccessToken(accessToken, util.ACCESS_TOKEN)
}

func (a *Auth) UniversalAuthLogin(clientID string, clientSecret string) (token string, err error) {

	if clientID == "" {
		clientID = os.Getenv(util.INFISICAL_UNIVERSAL_AUTH_CLIENT_ID_ENV_NAME)
	}
	if clientSecret == "" {
		clientSecret = os.Getenv(util.INFISICAL_UNIVERSAL_AUTH_CLIENT_SECRET_ENV_NAME)
	}

	token, err = api.CallUniversalAuthLogin(a.client.httpClient, api.UniversalAuthLoginRequest{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	})

	if err != nil {
		return "", err
	}

	a.client.setAccessToken(token, util.UNIVERSAL_AUTH)
	return token, nil

}

func (a *Auth) KubernetesAuthLogin(identityID string, serviceAccountTokenPath string) (accessToken string, err error) {

	if serviceAccountTokenPath == "" {
		serviceAccountTokenPath = os.Getenv(util.DEFAULT_KUBERNETES_SERVICE_ACCOUNT_TOKEN_PATH)
	}
	if identityID == "" {
		identityID = os.Getenv(util.INFISICAL_KUBERNETES_IDENTITY_ID_ENV_NAME)
	}

	serviceAccountToken, serviceAccountTokenErr := util.GetKubernetesServiceAccountToken(serviceAccountTokenPath)

	if serviceAccountTokenErr != nil {
		return "", serviceAccountTokenErr
	}

	accessToken, err = api.CallKubernetesAuthLogin(a.client.httpClient, api.KubernetesAuthLoginRequest{
		IdentityID: identityID,
		JWT:        serviceAccountToken,
	})

	if err != nil {
		return "", err
	}

	a.client.setAccessToken(accessToken, util.KUBERNETES)
	return accessToken, nil

}

func (a *Auth) AzureAuthLogin(identityID string) (accessToken string, err error) {
	if identityID == "" {
		identityID = os.Getenv(util.INFISICAL_AZURE_AUTH_IDENTITY_ID_ENV_NAME)
	}

	jwt, jwtError := util.GetAzureMetadataToken(a.client.httpClient)

	if jwtError != nil {
		return "", jwtError
	}

	accessToken, err = api.CallAzureAuthLogin(a.client.httpClient, api.AzureAuthLoginRequest{
		IdentityID: identityID,
		JWT:        jwt,
	})

	if err != nil {
		return "", err
	}

	a.client.setAccessToken(accessToken, util.AZURE)
	return accessToken, nil
}

func (a *Auth) GcpIdTokenAuthLogin(identityID string) (accessToken string, err error) {
	if identityID == "" {
		identityID = os.Getenv(util.INFISICAL_GCP_AUTH_IDENTITY_ID_ENV_NAME)
	}

	jwt, jwtError := util.GetGCPMetadataToken(a.client.httpClient, identityID)

	if jwtError != nil {
		return "", jwtError
	}

	accessToken, err = api.CallGCPAuthLogin(a.client.httpClient, api.GCPAuthLoginRequest{
		IdentityID: identityID,
		JWT:        jwt,
	})

	if err != nil {
		return "", err
	}

	a.client.setAccessToken(accessToken, util.GCP_ID_TOKEN)
	return accessToken, nil
}

type AwsIamRequestData struct {
	HTTPRequestMethod string `json:"iamHttpRequestMethod"`
	IamRequestBody    string `json:"iamRequestBody"`
	IamRequestHeaders string `json:"iamRequestHeaders"`
	IdentityId        string `json:"identityId"`
}

func (a *Auth) GcpIamAuthLogin(identityID string, serviceAccountKeyFilePath string) (accessToken string, err error) {
	if identityID == "" {
		identityID = os.Getenv(util.INFISICAL_GCP_AUTH_IDENTITY_ID_ENV_NAME)
	}
	if serviceAccountKeyFilePath == "" {
		serviceAccountKeyFilePath = os.Getenv(util.INFISICAL_GCP_IAM_SERVICE_ACCOUNT_KEY_FILE_PATH_ENV_NAME)
	}

	jwt, jwtError := util.GetGCPIamServiceAccountToken(identityID, serviceAccountKeyFilePath)

	if jwtError != nil {
		return "", jwtError
	}

	accessToken, err = api.CallGCPAuthLogin(a.client.httpClient, api.GCPAuthLoginRequest{
		IdentityID: identityID,
		JWT:        jwt,
	})

	if err != nil {
		return "", err
	}

	a.client.setAccessToken(accessToken, util.GCP_IAM)
	return accessToken, nil
}

func (a *Auth) AwsIamAuthLogin(identityId string) (accessToken string, err error) {

	if identityId == "" {
		identityId = os.Getenv(util.INFISICAL_AWS_IAM_AUTH_IDENTITY_ID_ENV_NAME)
	}

	awsRegion, regionErr := util.GetAwsRegion()

	if regionErr != nil {
		return "", regionErr
	}

	awsCfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if err != nil {
		return "", fmt.Errorf("unable to load SDK config, %v", err)
	}

	stsClient := sts.NewFromConfig(awsCfg)

	// You can use the stsClient to perform operations if needed
	// For example, calling GetCallerIdentity to validate credentials
	_, err = stsClient.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})

	if err != nil {
		return "", fmt.Errorf("failed to retrieve credentials, %v", err)
	}

	// Prepare request for signing
	iamRequestURL := fmt.Sprintf("https://sts.%s.amazonaws.com/", awsRegion)
	iamRequestBody := "Action=GetCallerIdentity&Version=2011-06-15"

	req, err := http.NewRequest(http.MethodPost, iamRequestURL, strings.NewReader(iamRequestBody))
	if err != nil {
		return "", fmt.Errorf("error creating HTTP request: %v", err)
	}

	currentTime := time.Now().UTC()

	req.Header.Add("X-Amz-Date", currentTime.Format("20060102T150405Z"))
	// req.Header.Add("Host", fmt.Sprintf("sts.%s.amazonaws.com", awsRegion))

	// creds, err := awsCfg.Credentials.Retrieve(context.Background())

	// credentials := credentials.NewCredentials(&credentials.StaticProvider{Value: credentials.Value{
	// 	AccessKeyID:     creds.AccessKeyID,
	// 	SecretAccessKey: creds.SecretAccessKey,
	// 	SessionToken:    creds.SessionToken,
	// }})

	credentials, err := awsCfg.Credentials.Retrieve(context.TODO())

	if err != nil {
		return "", fmt.Errorf("error retrieving credentials: %v", err)
	}

	hasher := sha256.New()
	hasher.Write([]byte(iamRequestBody))
	payloadHash := fmt.Sprintf("%x", hasher.Sum(nil))

	signer := v4.NewSigner()

	err = signer.SignHTTP(context.TODO(), credentials, req, payloadHash, "sts", awsRegion, time.Now())

	if err != nil {
		return "", fmt.Errorf("error signing request: %v", err)
	}

	var realHeaders map[string]string = make(map[string]string)

	for name, values := range req.Header {
		if strings.ToLower(name) == "content-length" {
			continue
		}
		fmt.Printf("Header: %v has value: %v\n\n", name, values)
		realHeaders[name] = values[0]
	}

	realHeaders["Host"] = fmt.Sprintf("sts.%s.amazonaws.com", awsRegion)
	realHeaders["Content-Type"] = "application/x-www-form-urlencoded; charset=utf-8"
	realHeaders["Content-Length"] = fmt.Sprintf("%d", len(iamRequestBody))

	// convert the headers to a json marshalled string
	jsonStringHeaders, err := json.Marshal(realHeaders)

	fmt.Printf("Test: 11\n")

	if err != nil {
		return "", fmt.Errorf("error marshalling headers: %v", err)
	}

	fmt.Printf("Test: 12\n")

	iamRequestData := AwsIamRequestData{
		HTTPRequestMethod: req.Method,
		// Encoding is intended, we decode it on severside, and I know everything happening on the server is being done correctly. So it's something broken in this code somewhere.
		IamRequestBody:    base64.StdEncoding.EncodeToString([]byte(iamRequestBody)),
		IamRequestHeaders: base64.StdEncoding.EncodeToString(jsonStringHeaders),
		IdentityId:        identityId,
	}

	res, _ := a.client.httpClient.R().
		SetBody(iamRequestData).
		Post("/v1/auth/aws-auth/login")

	fmt.Printf("Response status code: %v\n", res.StatusCode())
	fmt.Printf("Response body: %v\n", res.String())

	fmt.Printf("Test: 15\n")

	fmt.Printf("iamRequestData: %v\n", iamRequestData)

	return "", nil
}

func NewAuth(client *InfisicalClient) AuthInterface {
	return &Auth{client: client}
}
