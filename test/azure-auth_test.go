package test

import (
	"context"
	"fmt"
	"testing"

	infisical "github.com/infisical/go-sdk"
)

const AZURE_AUTH_IDENTITY_ID = "ed99db9d-5793-476c-8702-ee040669ae0d"

// ubuntu@74.243.217.7 /// normal password uppercase

func TestAzureLogin(t *testing.T) {

	client := infisical.NewInfisicalClient(context.Background(), infisical.Config{})

	_, err := client.Auth().AwsIamAuthLogin(AWS_AUTH_IDENTITY_ID)
	if err != nil {
		fmt.Printf("AWS Auth Error: %v\n", err)
	}

	_, err = client.Auth().UniversalAuthLogin(UNIVERSAL_AUTH_CLIENT_ID, UNIVERSAL_AUTH_CLIENT_SECRET)
	if err != nil {
		fmt.Printf("Universal Auth Error: %v\n", err)
	}

}
