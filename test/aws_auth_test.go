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

	accessToken, err := client.Auth().AwsIamAuthLogin("0e007fbe-7954-48f9-b888-6665eec088e7") // test ID

	if err != nil {
		panic(err)
	}

	fmt.Println(accessToken)
}
