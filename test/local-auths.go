package test

import (
	"context"
	"fmt"
	"testing"

	infisical "github.com/infisical/go-sdk"
)

const AWS_AUTH_IDENTITY_ID = "99c9c780-39c5-413d-8e78-3c5d3a113a91"

const (
	UNIVERSAL_AUTH_CLIENT_ID     = "ee2a5906-d93e-42d2-8649-d9c047053271"
	UNIVERSAL_AUTH_CLIENT_SECRET = "a8769cab25eced271b27fb42755e890d7186c4cc0bd3adb6b739966a0da3ab38"
)

func TestAWSAuthLogin(t *testing.T) {

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
