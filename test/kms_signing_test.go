package test

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"testing"

	infisical "github.com/infisical/go-sdk"
)

func TestKmsSigning(t *testing.T) {

	client := infisical.NewInfisicalClient(context.Background(), infisical.Config{
		SiteUrl: "http://localhost:8080",
	})
	_, err := client.Auth().UniversalAuthLogin("<>", "<>")

	if err != nil {
		panic(err)
	}

	randomBytes := make([]byte, 16)
	rand.Read(randomBytes)

	randomString := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(base64.StdEncoding.EncodeToString(randomBytes)), "=", ""), "/", ""))

	fmt.Printf("Random string: %s\n\n\n", randomString)

	newKey, err := client.Kms().Keys().Create(infisical.KmsCreateKeyOptions{
		KeyUsage:            "sign-verify",
		Description:         "test key",
		Name:                randomString,
		EncryptionAlgorithm: "rsa-4096",
		ProjectId:           "c54edf7a-f861-4131-afdc-b0ad5faec5dc",
	})

	if err != nil {
		panic(err)
	}

	signingAlgorithms, err := client.Kms().Signing().ListSigningAlgorithms(infisical.KmsListSigningAlgorithmsOptions{
		KeyId: newKey.KeyId,
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("Signing algorithms: %+v\n\n\n", signingAlgorithms)

	selectedSigningAlgorithm := signingAlgorithms[0]

	testData := "Hello, World!"

	res, err := client.Kms().Signing().SignData(infisical.KmsSignDataOptions{
		KeyId:            newKey.KeyId,
		Data:             base64.StdEncoding.EncodeToString([]byte(testData)),
		SigningAlgorithm: selectedSigningAlgorithm,
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("Signature: %s\n\n\n", res.Signature)

	verifyRes, err := client.Kms().Signing().VerifyData(infisical.KmsVerifyDataOptions{
		KeyId:            newKey.KeyId,
		Data:             base64.StdEncoding.EncodeToString([]byte(testData)),
		Signature:        res.Signature,
		SigningAlgorithm: selectedSigningAlgorithm,
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("Verification result: %+v\n\n\n", verifyRes)

	publicKey, err := client.Kms().Signing().GetPublicKey(infisical.KmsGetPublicKeyOptions{
		KeyId: newKey.KeyId,
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("Public key: %s\n\n\n", publicKey)

	_, err = client.Kms().Keys().Delete(infisical.KmsDeleteKeyOptions{
		KeyId: newKey.KeyId,
	})

	if err != nil {
		panic(err)
	}

}
