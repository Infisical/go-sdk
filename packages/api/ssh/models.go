package api

import (
	"github.com/infisical/go-sdk/packages/util"
)

type SignSshPublicKeyV1Request struct {
	CertificateTemplateID string           `json:"certificateTemplateId"`
	PublicKey             string           `json:"publicKey"`
	CertType              util.SshCertType `json:"certType,omitempty"`
	Principals            []string         `json:"principals"`
	TTL                   string           `json:"ttl,omitempty"`
	KeyID                 string           `json:"keyId,omitempty"`
}

type SignSshPublicKeyV1Response struct {
	SerialNumber string `json:"serialNumber"`
	SignedKey    string `json:"signedKey"`
}

type IssueSshCredsV1Request struct {
	CertificateTemplateID string                `json:"certificateTemplateId"`
	KeyAlgorithm          util.CertKeyAlgorithm `json:"keyAlgorithm,omitempty"`
	CertType              util.SshCertType      `json:"certType,omitempty"`
	Principals            []string              `json:"principals"`
	TTL                   string                `json:"ttl,omitempty"`
	KeyID                 string                `json:"keyId,omitempty"`
}

type IssueSshCredsV1Response struct {
	SerialNumber string                `json:"serialNumber"`
	SignedKey    string                `json:"signedKey"`
	PrivateKey   string                `json:"privateKey"`
	PublicKey    string                `json:"publicKey"`
	KeyAlgorithm util.CertKeyAlgorithm `json:"keyAlgorithm"`
}
