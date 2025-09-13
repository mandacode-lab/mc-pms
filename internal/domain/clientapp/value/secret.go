package clientappval

import (
	"errors"
)

var ErrInvalidSecretHash = errors.New("invalid secret hash")

type SecretHash struct {
	value []byte
}

func NewSecretHash(hash []byte) SecretHash {
	return SecretHash{value: hash}
}



func HashSecretFromBytes(plainSecret []byte, hash []byte) SecretHash {
	return SecretHash{value: hash}
}

func HashSecret(plainSecret string, hash []byte) SecretHash {
	return SecretHash{value: hash}
}



func (sh SecretHash) Value() []byte {
	return sh.value
}

func (sh SecretHash) VerifySecretBytes(plainSecret []byte, hash []byte) bool {
	expectedHash := HashSecretFromBytes(plainSecret, hash)
	return string(sh.value) == string(expectedHash.value)
}

func (sh SecretHash) VerifySecret(plainSecret string, hash []byte) bool {
	expectedHash := HashSecret(plainSecret, hash)
	return string(sh.value) == string(expectedHash.value)
}
