package api

import (
	"github.com/infisical/go-sdk/packages/models"
)

type SignSshPublicKeyV1Request struct {
	ProjectID    string                  `json:"projectId"`
	TemplateName string                  `json:"templateName"`
	PublicKey    string                  `json:"publicKey"`
	KeyAlgorithm models.CertKeyAlgorithm `json:"keyAlgorithm,omitempty"`
	CertType     models.SshCertType      `json:"certType,omitempty"`
	Principals   []string                `json:"principals"`
	TTL          string                  `json:"ttl,omitempty"`
	KeyID        string                  `json:"keyId,omitempty"`
}

type SignSshPublicKeyV1Response struct {
	SerialNumber string `json:"serialNumber"`
	SignedKey    string `json:"signedKey"`
}

type IssueSshCredsV1Request struct {
	ProjectID    string                  `json:"projectId"`
	TemplateName string                  `json:"templateName"`
	KeyAlgorithm models.CertKeyAlgorithm `json:"keyAlgorithm,omitempty"`
	CertType     models.SshCertType      `json:"certType,omitempty"`
	Principals   []string                `json:"principals"`
	TTL          string                  `json:"ttl,omitempty"`
	KeyID        string                  `json:"keyId,omitempty"`
}

type IssueSshCredsV1Response struct {
	SerialNumber string                  `json:"serialNumber"`
	SignedKey    string                  `json:"signedKey"`
	PrivateKey   string                  `json:"privateKey"`
	PublicKey    string                  `json:"publicKey"`
	KeyAlgorithm models.CertKeyAlgorithm `json:"keyAlgorithm"`
}
