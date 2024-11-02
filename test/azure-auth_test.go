package test

import (
	"context"
	"fmt"
	"testing"

	infisical "github.com/infisical/go-sdk"
)

const AZURE_AUTH_IDENTITY_ID = "ed99db9d-5793-476c-8702-ee040669ae0d"

func CallListSecretsAzure(client infisical.InfisicalClientInterface) error {
	_, err := client.Secrets().List(infisical.ListSecretsOptions{
		ProjectID:   "437f23bb-86ce-4766-9861-e8ba49fd7e95",
		Environment: "dev",
	})
	if err != nil {
		return err
	}
	return nil
}

// ubuntu@74.243.217.7 /// normal password uppercase

func TestAzureLogin(t *testing.T) {

	client := infisical.NewInfisicalClient(context.Background(), infisical.Config{})

	_, err := client.Auth().AzureAuthLogin(AZURE_AUTH_IDENTITY_ID, "")
	if err != nil {
		fmt.Printf("Azure Auth Error: %v\n", err)
	}

	err = CallListSecretsAzure(client)

	if err != nil {
		fmt.Printf("List Secrets Error: %v\n", err)
	}

}
