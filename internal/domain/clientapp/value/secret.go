package clientappval

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

var ErrInvalidSecretHash = errors.New("invalid secret hash")

type SecretHash struct {
	value []byte
}

func NewSecretHash(hash []byte) SecretHash {
	return SecretHash{value: hash}
}

func HashSecret(plainSecret string) SecretHash {
	hash := sha256.Sum256([]byte(plainSecret))
	return SecretHash{value: hash[:]}
}

func GenerateSecret() (string, SecretHash, error) {
	secretBytes := make([]byte, 32)
	_, err := rand.Read(secretBytes)
	if err != nil {
		return "", SecretHash{}, err
	}

	plainSecret := base64.URLEncoding.EncodeToString(secretBytes)
	hash := HashSecret(plainSecret)

	return plainSecret, hash, nil
}

func (sh SecretHash) Value() []byte {
	return sh.value
}

func (sh SecretHash) VerifySecret(plainSecret string) bool {
	hash := HashSecret(plainSecret)
	return string(sh.value) == string(hash.value)
}
