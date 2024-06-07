package util

import (
	"os"
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
