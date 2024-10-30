package api

type KmsEncryptDataV1Request struct {
	KeyId     string
	Plaintext string `json:"plaintext"`
}

type KmsEncryptDataV1Response struct {
	Ciphertext string `json:"ciphertext"`
}

type KmsDecryptDataV1Request struct {
	KeyId      string
	Ciphertext string `json:"ciphertext"`
}

type KmsDecryptDataV1Response struct {
	Plaintext string `json:"plaintext"`
}
