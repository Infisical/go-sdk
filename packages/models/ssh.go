package models

type CertKeyAlgorithm string

const (
	RSA2048   CertKeyAlgorithm = "RSA_2048"
	RSA4096   CertKeyAlgorithm = "RSA_4096"
	ECDSAP256 CertKeyAlgorithm = "EC_prime256v1"
	ECDSAP384 CertKeyAlgorithm = "EC_secp384r1"
)

type SshCertType string

const (
	UserCert SshCertType = "user"
	HostCert SshCertType = "host"
)

// type SignedSshKeyWithKeyPair struct {
// 	SerialNumber  string           `json:"serialNumber"`  // Unique serial number for the SSH key
// 	SignedKey     string           `json:"signedKey"`     // Signed SSH public key (certificate)
// 	PrivateKey    string           `json:"privateKey"`    // Corresponding private key
// 	PublicKey     string           `json:"publicKey"`     // Corresponding public key
// 	KeyAlgorithm  CertKeyAlgorithm `json:"keyAlgorithm"`  // Algorithm used for the key
// }