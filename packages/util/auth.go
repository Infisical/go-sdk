package util

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	credentials "cloud.google.com/go/iam/credentials/apiv1"
	"cloud.google.com/go/iam/credentials/apiv1/credentialspb"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/go-resty/resty/v2"
	"google.golang.org/api/option"
)

func GetKubernetesServiceAccountToken(serviceAccountTokenPath string) (string, error) {

	if serviceAccountTokenPath == "" {
		serviceAccountTokenPath = DEFAULT_KUBERNETES_SERVICE_ACCOUNT_TOKEN_PATH
	}

	token, err := os.ReadFile(serviceAccountTokenPath)

	if err != nil {
		return "", err
	}

	return string(token), nil

}

// AzureIMDSIdentitySelector pins which managed identity IMDS should issue the token for
// when multiple user-assigned managed identities are attached to the host. ClientID is
// preferred; ObjectID is used only when ClientID is empty. When both are empty, IMDS uses
// the system-assigned identity (the legacy SDK behaviour).
type AzureIMDSIdentitySelector struct {
	ClientID string
	ObjectID string
}

// AzureWorkloadIdentityOptions allows callers to override the standard AZURE_* environment
// variables consumed by Azure Workload Identity. Empty fields fall back to the env vars.
type AzureWorkloadIdentityOptions struct {
	ClientID      string
	TenantID      string
	TokenFilePath string
	AuthorityHost string
}

func buildAzureMetadataServiceURL(resource string, sel AzureIMDSIdentitySelector) string {
	metadataURL := AZURE_METADATA_SERVICE_URL
	if resource != "" {
		metadataURL += url.PathEscape(resource)
	} else {
		metadataURL += AZURE_DEFAULT_RESOURCE
	}

	switch {
	case sel.ClientID != "":
		metadataURL += "&client_id=" + url.QueryEscape(sel.ClientID)
	case sel.ObjectID != "":
		metadataURL += "&object_id=" + url.QueryEscape(sel.ObjectID)
	}

	return metadataURL
}

// GetAzureMetadataToken acquires an Az Entra ID access token from the Azure Instance Metadata
// Service (IMDS). When the host has multiple user-assigned managed identities attached,
// pass a non-empty AzureIMDSIdentitySelector to disambiguate which one IMDS should use,
// otherwise IMDS will fail with HTTP 400.
func GetAzureMetadataToken(httpClient *resty.Client, customResource string, sel AzureIMDSIdentitySelector) (string, error) {

	type AzureMetadataResponse struct {
		AccessToken string `json:"access_token"`
	}

	metadataResponse := AzureMetadataResponse{}

	response, err := httpClient.R().
		SetResult(&metadataResponse).
		SetHeader("Metadata", "true").
		SetHeader("Accept", "application/json").
		Get(buildAzureMetadataServiceURL(customResource, sel))

	if err != nil {
		return "", err
	}

	if response.IsError() {
		return "", fmt.Errorf("GetAzureMetadataToken: Unsuccessful response [%v %v] [status-code=%v] [Error: %s]", response.Request.Method, response.Request.URL, response.StatusCode(), TryParseErrorBody(response))
	}

	return metadataResponse.AccessToken, nil
}

// GetAzureWorkloadIdentityToken acquires an Az Entra ID access token via Azure Workload Identity
// using the official azidentity client. The resource argument follows the same legacy
// contract used by GetAzureMetadataToken: it may be a URL-encoded value (e.g. the
// AZURE_DEFAULT_RESOURCE constant) or a plain URL such as "https://management.azure.com/".
// The function normalises it and appends the "/.default" suffix required by Az Entra ID v2 scopes.
func GetAzureWorkloadIdentityToken(ctx context.Context, customResource string, opts AzureWorkloadIdentityOptions) (string, error) {

	credOpts := &azidentity.WorkloadIdentityCredentialOptions{
		ClientID:      opts.ClientID,
		TenantID:      opts.TenantID,
		TokenFilePath: opts.TokenFilePath,
	}

	authorityHost := opts.AuthorityHost
	if authorityHost == "" {
		authorityHost = os.Getenv(AZURE_AUTHORITY_HOST_ENV_NAME)
	}
	if authorityHost != "" {
		credOpts.ClientOptions = policy.ClientOptions{
			Cloud: cloud.Configuration{ActiveDirectoryAuthorityHost: authorityHost},
		}
	}

	cred, err := azidentity.NewWorkloadIdentityCredential(credOpts)
	if err != nil {
		return "", fmt.Errorf("GetAzureWorkloadIdentityToken: failed to create credential: %w", err)
	}

	tk, err := cred.GetToken(ctx, policy.TokenRequestOptions{
		Scopes: []string{normaliseAzureScope(customResource)},
	})
	if err != nil {
		return "", fmt.Errorf("GetAzureWorkloadIdentityToken: failed to acquire token: %w", err)
	}

	return tk.Token, nil
}

// normaliseAzureScope decodes a possibly URL-encoded resource string and returns the Az Entra ID
// v2 scope form ("<resource>/.default"). Empty input falls back to the default management
// resource so behaviour matches GetAzureMetadataToken.
func normaliseAzureScope(resource string) string {
	if resource == "" {
		resource = AZURE_DEFAULT_RESOURCE
	}
	if decoded, err := url.QueryUnescape(resource); err == nil {
		resource = decoded
	}
	resource = strings.TrimSuffix(resource, "/")
	if strings.HasSuffix(resource, "/.default") {
		return resource
	}
	return resource + "/.default"
}

func GetGCPMetadataToken(httpClient *resty.Client, identityID string) (string, error) {

	res, err := httpClient.R().
		SetHeader("Metadata-Flavor", "Google").
		Get(fmt.Sprintf("http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/identity?audience=%s&format=full", identityID))

	if err != nil {
		return "", err
	}

	if res.IsError() {
		return "", fmt.Errorf("GetGCPMetadataToken: Unsuccessful response [%v %v] [status-code=%v] [Error: %s]", res.Request.Method, res.Request.URL, res.StatusCode(), TryParseErrorBody(res))
	}

	return res.String(), nil

}

func GetAwsEC2IdentityDocumentRegion(timeout int) (string, error) {

	type AwsIdentityDocument struct {
		Region string `json:"region"`
	}

	httpClient := resty.New().SetTimeout(time.Duration(timeout) * time.Millisecond)

	res, err := httpClient.R().
		SetHeader("X-aws-ec2-metadata-token-ttl-seconds", "21600").
		Put(AWS_EC2_METADATA_TOKEN_URL)

	if err != nil {
		return "", err
	}

	if res.IsError() {
		return "", fmt.Errorf("GetAwsEC2IdentityDocumentRegion: Unsuccessful response [%v %v] [status-code=%v] [Error: %s]", res.Request.Method, res.Request.URL, res.StatusCode(), TryParseErrorBody(res))
	}

	metadataToken := res.String()

	res, err = httpClient.R().
		SetHeader("X-aws-ec2-metadata-token", metadataToken).
		SetHeader("Accept", "application/json").
		Get(AWS_EC2_INSTANCE_IDENTITY_DOCUMENT_URL)

	if err != nil {
		return "", err
	}

	if res.IsError() {
		return "", fmt.Errorf("GetAwsEC2IdentityDocumentRegion: Unsuccessful response [%v %v] [status-code=%v] [Error: %s]", res.Request.Method, res.Request.URL, res.StatusCode(), TryParseErrorBody(res))
	}

	// For some reason using .SetResult(&AwsIdentityDocument{}) doesn't work and just results in an empty object. This works though..
	var identityDocument AwsIdentityDocument
	err = json.Unmarshal(res.Body(), &identityDocument)
	if err != nil {
		return "", err
	}

	return identityDocument.Region, nil

}

func GetGCPIamServiceAccountToken(identityID string, serviceAccountKeyPath string) (string, error) {

	type JwtPayload struct {
		Sub string `json:"sub"`
		Aud string `json:"aud"`
	}

	ctx := context.Background()

	serviceAccountKey, err := os.ReadFile(serviceAccountKeyPath)
	if err != nil {
		return "", err
	}

	var creds map[string]string
	if err := json.Unmarshal(serviceAccountKey, &creds); err != nil {
		return "", fmt.Errorf("failed to unmarshal service account key: %v", err)
	}

	clientEmail := creds["client_email"]
	if clientEmail == "" {
		return "", fmt.Errorf("client email not found in service account key")
	}

	payload := JwtPayload{
		Sub: clientEmail,
		Aud: identityID,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JWT payload: %v", err)
	}

	iamCredentialsClient, err := credentials.NewIamCredentialsClient(ctx, option.WithCredentialsFile(serviceAccountKeyPath)) //nolint:staticcheck // deprecated but no drop-in replacement available yet
	if err != nil {
		return "", fmt.Errorf("failed to create IAM credentials client: %v", err)
	}

	defer iamCredentialsClient.Close() //nolint:errcheck

	signJwtRequest := &credentialspb.SignJwtRequest{
		Name:    fmt.Sprintf("projects/-/serviceAccounts/%s", clientEmail),
		Payload: string(payloadJSON),
	}

	resp, err := iamCredentialsClient.SignJwt(ctx, signJwtRequest)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %v. Ensure the IAM Service Account Credentials API is enabled", err)
	}

	signedJwt := resp.SignedJwt
	if signedJwt == "" {
		return "", fmt.Errorf("failed to sign JWT: signedJwt is empty")
	}

	return signedJwt, nil

}

func GetAwsRegion() (string, error) {
	// in Lambda environments, the region is available in the AWS_REGION environment variable
	region := os.Getenv("AWS_REGION")

	if region != "" {
		return region, nil
	}

	// in EC2 environments, the region is available in the identity doc

	region, err := GetAwsEC2IdentityDocumentRegion(5000)

	if err != nil {
		return "", err
	}

	return region, nil

}

func RetrieveAwsCredentials() (credentials aws.Credentials, region string, err error) {
	presetAwsCfg, err := config.LoadDefaultConfig(context.TODO())

	if err == nil && presetAwsCfg.Region != "" {
		creds, err := presetAwsCfg.Credentials.Retrieve(context.TODO())
		if err == nil {
			return creds, presetAwsCfg.Region, nil
		}
	}

	awsRegion, err := GetAwsRegion()
	if err != nil {
		return aws.Credentials{}, "", err
	}

	awsCfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if err != nil {
		return aws.Credentials{}, "", fmt.Errorf("unable to load SDK config, %v", err)
	}

	creds, err := awsCfg.Credentials.Retrieve(context.TODO())
	if err != nil {
		return aws.Credentials{}, "", fmt.Errorf("error retrieving credentials: %v", err)
	}

	return creds, awsRegion, nil
}
