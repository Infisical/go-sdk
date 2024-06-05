package infisical

type UniversalAuth struct {
	clientId     string
	clientSecret string
}

type GcpIdTokenAuth struct {
	identityId string
}

type GcpIamAuth struct {
	identityId                string
	serviceAccountKeyFilePath string
}

type AwsIamAuth struct {
	identityId string
}

type AzureAuth struct {
	identityId string
}

type KubernetesAuth struct {
	identityId              string
	serviceAccountTokenPath string
}

type Authentication struct {
	universalAuth  UniversalAuth
	gcpIdTokenAuth GcpIdTokenAuth
	gcpIamAuth     GcpIamAuth
	awsIamAuth     AwsIamAuth
	azureAuth      AzureAuth
	kubernetesAuth KubernetesAuth

	accessToken string
}
