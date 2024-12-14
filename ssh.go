package infisical

import (
	api "github.com/infisical/go-sdk/packages/api/ssh"
)

type SignSshPublicKeyOptions = api.SignSshPublicKeyV1Request
type IssueSshCredsOptions = api.IssueSshCredsV1Request

type SshInterface interface {
	SignKey(options SignSshPublicKeyOptions) (api.SignSshPublicKeyV1Response, error)
	IssueCredentials(options IssueSshCredsOptions) (api.IssueSshCredsV1Response, error)
}

type Ssh struct {
	client *InfisicalClient
}

func (f *Ssh) SignKey(options SignSshPublicKeyOptions) (api.SignSshPublicKeyV1Response, error) {
	res, err := api.CallSignSshPublicKeyV1(f.client.httpClient, options)

	if err != nil {
		return api.SignSshPublicKeyV1Response{}, err
	}

	return res, nil
}

func (f *Ssh) IssueCredentials(options IssueSshCredsOptions) (api.IssueSshCredsV1Response, error) {
	res, err := api.CallIssueSshCredsV1(f.client.httpClient, options)

	if err != nil {
		return api.IssueSshCredsV1Response{}, err
	}

	return res, nil
}

func NewSsh(client *InfisicalClient) SshInterface {
	return &Ssh{client: client}
}
