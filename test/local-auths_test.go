package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	infisical "github.com/infisical/go-sdk"
)

const AWS_AUTH_IDENTITY_ID = "99c9c780-39c5-413d-8e78-3c5d3a113a91"

const (
	UNIVERSAL_AUTH_CLIENT_ID     = "ee2a5906-d93e-42d2-8649-d9c047053271"
	UNIVERSAL_AUTH_CLIENT_SECRET = "a8769cab25eced271b27fb42755e890d7186c4cc0bd3adb6b739966a0da3ab38"
)

const TOKEN_AUTH = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZGVudGl0eUlkIjoiMWMyZGQyY2QtMGUzMi00ZjI3LWJhMjUtY2ZjNDM2NzU4MGZjIiwiaWRlbnRpdHlBY2Nlc3NUb2tlbklkIjoiNTNmY2M5NjEtNWY2OS00M2JkLTljOGYtMjgyYWRjZjQ0N2VlIiwiYXV0aFRva2VuVHlwZSI6ImlkZW50aXR5QWNjZXNzVG9rZW4iLCJpYXQiOjE3MzA1NDc5MzMsImV4cCI6MTczMzEzOTkzM30.Ogd57j3m5UeUNY2fVMKpnZ_L8XLSCx_aw7G6Lyu57VQ"

func CallListSecrets(client infisical.InfisicalClientInterface) error {
	_, err := client.Secrets().List(infisical.ListSecretsOptions{
		ProjectID:   "437f23bb-86ce-4766-9861-e8ba49fd7e95",
		Environment: "dev",
	})
	if err != nil {
		return err
	}
	return nil
}

func AwsIAmLogin() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := infisical.NewInfisicalClient(ctx, infisical.Config{})

	_, err := client.Auth().AwsIamAuthLogin(AWS_AUTH_IDENTITY_ID)
	if err != nil {
		fmt.Printf("AWS Auth Error: %v\n", err)
	}

	err = CallListSecrets(client)

	if err != nil {
		return err
	}

	return nil
}

func UniversalAuthLogin() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := infisical.NewInfisicalClient(ctx, infisical.Config{})

	_, err := client.Auth().UniversalAuthLogin(UNIVERSAL_AUTH_CLIENT_ID, UNIVERSAL_AUTH_CLIENT_SECRET)
	if err != nil {
		fmt.Printf("Universal Auth Error: %v\n", err)
	}

	err = CallListSecrets(client)

	if err != nil {
		return err
	}

	return nil
}

func AccessTokenLogin() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := infisical.NewInfisicalClient(ctx, infisical.Config{})

	client.Auth().SetAccessToken(TOKEN_AUTH)

	err := CallListSecrets(client)

	if err != nil {
		return err
	}

	return nil
}

func TestAWSAuthLogin(t *testing.T) {

	for {
		fmt.Printf("Testing AWS Auth\n")
		err := AwsIAmLogin()
		if err != nil {
			fmt.Printf("AWS Auth Error: %v\n", err)
		}

		fmt.Printf("Testing Universal Auth\n")
		err = UniversalAuthLogin()
		if err != nil {
			fmt.Printf("Universal Auth Error: %v\n", err)
		}

		fmt.Printf("Testing Access Token Auth\n")
		err = AccessTokenLogin()
		if err != nil {
			fmt.Printf("Access Token Auth Error: %v\n", err)
		}

		time.Sleep(5 * time.Second)

	}

}
