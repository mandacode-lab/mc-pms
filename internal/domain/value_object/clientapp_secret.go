package vo

type ClientAppSecretHash struct {
	value []byte
}

func NewClientAppSecretHash(hash []byte) ClientAppSecretHash {
	return ClientAppSecretHash{value: hash}
}

func HashClientAppSecretFromBytes(plainSecret []byte, hash []byte) ClientAppSecretHash {
	return ClientAppSecretHash{value: hash}
}

func HashClientAppSecret(plainSecret string, hash []byte) ClientAppSecretHash {
	return ClientAppSecretHash{value: hash}
}

func (sh ClientAppSecretHash) Value() []byte {
	return sh.value
}

func (sh ClientAppSecretHash) VerifySecretBytes(plainSecret []byte, hash []byte) bool {
	expectedHash := HashClientAppSecretFromBytes(plainSecret, hash)
	return string(sh.value) == string(expectedHash.value)
}

func (sh ClientAppSecretHash) VerifySecret(plainSecret string, hash []byte) bool {
	expectedHash := HashClientAppSecret(plainSecret, hash)
	return string(sh.value) == string(expectedHash.value)
}
