package out

import "context"

// KekProvider provides Key Encryption Key services
type KekProvider interface {
	// Encrypt encrypts data using the Key Encryption Key
	Encrypt(ctx context.Context, plaintext []byte) (ciphertext, nonce []byte, err error)
	
	// Decrypt decrypts data using the Key Encryption Key
	Decrypt(ctx context.Context, ciphertext, nonce []byte) (plaintext []byte, err error)
	
	// GenerateDEK generates a new Data Encryption Key and wraps it with the KEK
	GenerateDEK(ctx context.Context) (wrappedDEK, dekNonce []byte, err error)
	
	// UnwrapDEK unwraps a Data Encryption Key using the KEK
	UnwrapDEK(ctx context.Context, wrappedDEK, dekNonce []byte) (dek []byte, err error)
	
	// EncryptWithDEK encrypts data using a Data Encryption Key
	EncryptWithDEK(ctx context.Context, dek, plaintext []byte) (ciphertext, nonce []byte, err error)
	
	// DecryptWithDEK decrypts data using a Data Encryption Key
	DecryptWithDEK(ctx context.Context, dek, ciphertext, nonce []byte) (plaintext []byte, err error)
}