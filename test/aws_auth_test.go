package test

import (
	"fmt"
	"testing"

	infisical "github.com/infisical/go-sdk"
)

func TestAWSAuthLogin(t *testing.T) {

	client, err := infisical.NewInfisicalClient(infisical.Config{
		SiteUrl: "https://c61b724baab4.ngrok.app",
	})

	if err != nil {
		t.Fatalf("Failed to create Infisical client: %v", err)
	}

	accessToken, err := client.Auth().AwsIamAuthLogin("e2cddb75-a0e0-4c89-bfc0-4d536599f725") // test ID

	if err != nil {
		panic(err)
	}

	fmt.Println(accessToken)
}
