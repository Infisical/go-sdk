package util

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	credentials "cloud.google.com/go/iam/credentials/apiv1"
	"cloud.google.com/go/iam/credentials/apiv1/credentialspb"
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

func GetAzureMetadataToken(httpClient *resty.Client) (string, error) {

	type AzureMetadataResponse struct {
		AccessToken string `json:"access_token"`
	}

	metadataResponse := AzureMetadataResponse{}

	response, err := httpClient.R().
		SetResult(&metadataResponse).
		SetHeader("Metadata", "true").
		SetHeader("Accept", "application/json").
		Get(AZURE_METADATA_SERVICE_URL)

	if err != nil {
		return "", err
	}

	if response.IsError() {
		return "", fmt.Errorf("GetAzureMetadataToken: Unsuccessful response [%v %v] [status-code=%v] [Error: %s]", response.Request.Method, response.Request.URL, response.StatusCode(), TryParseErrorBody(response))
	}

	return response.String(), nil
}

func GetGCPMetadataToken(httpClient *resty.Client, identityID string) (string, error) {

	res, err := httpClient.R().
		SetHeader("Metadata-Flavor", "Google").
		Get(fmt.Sprintf("http://metadata/computeMetadata/v1/instance/service-accounts/default/identity?audience=%s&format=full", identityID))

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
		fmt.Printf("JSON unmarshal error: %v\n", err)
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

	iamCredentialsClient, err := credentials.NewIamCredentialsClient(ctx, option.WithCredentialsFile(serviceAccountKeyPath))
	if err != nil {
		return "", fmt.Errorf("failed to create IAM credentials client: %v", err)
	}

	defer iamCredentialsClient.Close()

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
