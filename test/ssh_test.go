package test

// import (
// 	"context"
// 	"fmt"
// 	"os"
// 	"testing"

// 	infisical "github.com/infisical/go-sdk"
// )

// func TestSshIssueCreds(t *testing.T) {
// 	client := infisical.NewInfisicalClient(context.Background(), infisical.Config{
//         SiteUrl:          "http://localhost:8080",
//         AutoTokenRefresh: true,
//     })

// 	// Authenticate using Universal Auth
//     _, err := client.Auth().UniversalAuthLogin(os.Getenv("GO_SDK_TEST_UNIVERSAL_AUTH_CLIENT_ID"), os.Getenv("GO_SDK_TEST_UNIVERSAL_AUTH_CLIENT_SECRET"))
//     if err != nil {
//         fmt.Printf("Authentication failed: %v\n", err)
//         os.Exit(1)
//     }

// 	// Test issuing SSH credentials
// 	creds, err := client.Ssh().IssueCredentials(infisical.IssueSshCredsOptions{
// 		ProjectID: os.Getenv("GO_SDK_TEST_PROJECT_ID"),
// 		TemplateName: "template-name",
// 		Principals: []string{"ec2-user"},
// 	})

//     if err != nil {
// 		t.Fatalf("Failed to issue SSH credentials: %v", err)
//     }

// 	// Test signing SSH public key
// 	creds2, err := client.Ssh().SignKey(infisical.SignSshPublicKeyOptions{
// 		ProjectID: os.Getenv("GO_SDK_TEST_PROJECT_ID"),
// 		TemplateName: "template-name",
// 		Principals: []string{"ec2-user"},
// 		PublicKey: "ssh-rsa ...",
// 	})

//     if err != nil {
// 		t.Fatalf("Failed to sign SSH public key: %v", err)
//     }

// 	fmt.Print("Newly-issued SSH credentials: ", creds)
// 	fmt.Print("Signed SSH credential: ", creds2)
// }
