package test

import (
	"fmt"
	"testing"

	infisical "github.com/infisical/go-sdk"
)

func TestAWSAuthLogin(t *testing.T) {

	client, err := infisical.NewInfisicalClient(infisical.Config{
		SiteUrl: "http://localhost:8080",
	})

	if err != nil {
		t.Fatalf("Failed to create Infisical client: %v", err)
	}

	accessToken, err := client.Auth().AwsIamAuthLogin("TEST")

	if err != nil {
		panic(err)
	}

	fmt.Println(accessToken)
}
