package kek

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/mandacode-com/mandacode-service-hub/internal/port/out"
)

const (
	DEKSize   = 32 // 256-bit key for AES-256
	NonceSize = 12 // 96-bit nonce for AES-GCM
)

type AESKekProvider struct {
	kek []byte // Master Key Encryption Key
}

func NewAESKekProvider(kek []byte) out.KekProvider {
	return &AESKekProvider{
		kek: kek,
	}
}

func NewAESKekProviderFromHex(kekHex string) (out.KekProvider, error) {
	kek, err := hex.DecodeString(kekHex)
	if err != nil {
		return nil, fmt.Errorf("decoding KEK hex: %w", err)
	}

	if len(kek) != 32 {
		return nil, fmt.Errorf("KEK must be 32 bytes (256 bits), got %d bytes", len(kek))
	}

	return &AESKekProvider{
		kek: kek,
	}, nil
}

func (p *AESKekProvider) Encrypt(ctx context.Context, plaintext []byte) (ciphertext, nonce []byte, err error) {
	block, err := aes.NewCipher(p.kek)
	if err != nil {
		return nil, nil, fmt.Errorf("creating cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, fmt.Errorf("creating GCM: %w", err)
	}

	nonce = make([]byte, NonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return nil, nil, fmt.Errorf("generating nonce: %w", err)
	}

	ciphertext = aesGCM.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nonce, nil
}

func (p *AESKekProvider) Decrypt(ctx context.Context, ciphertext, nonce []byte) (plaintext []byte, err error) {
	block, err := aes.NewCipher(p.kek)
	if err != nil {
		return nil, fmt.Errorf("creating cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("creating GCM: %w", err)
	}

	plaintext, err = aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decrypting: %w", err)
	}

	return plaintext, nil
}

func (p *AESKekProvider) GenerateDEK(ctx context.Context) (wrappedDEK, dekNonce []byte, err error) {
	// Generate random DEK
	dek := make([]byte, DEKSize)
	if _, err := rand.Read(dek); err != nil {
		return nil, nil, fmt.Errorf("generating DEK: %w", err)
	}

	// Wrap DEK with KEK
	wrappedDEK, dekNonce, err = p.Encrypt(ctx, dek)
	if err != nil {
		return nil, nil, fmt.Errorf("wrapping DEK: %w", err)
	}

	return wrappedDEK, dekNonce, nil
}

func (p *AESKekProvider) UnwrapDEK(ctx context.Context, wrappedDEK, dekNonce []byte) (dek []byte, err error) {
	dek, err = p.Decrypt(ctx, wrappedDEK, dekNonce)
	if err != nil {
		return nil, fmt.Errorf("unwrapping DEK: %w", err)
	}

	return dek, nil
}

func (p *AESKekProvider) EncryptWithDEK(ctx context.Context, dek, plaintext []byte) (ciphertext, nonce []byte, err error) {
	block, err := aes.NewCipher(dek)
	if err != nil {
		return nil, nil, fmt.Errorf("creating cipher with DEK: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, fmt.Errorf("creating GCM with DEK: %w", err)
	}

	nonce = make([]byte, NonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return nil, nil, fmt.Errorf("generating nonce: %w", err)
	}

	ciphertext = aesGCM.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nonce, nil
}

func (p *AESKekProvider) DecryptWithDEK(ctx context.Context, dek, ciphertext, nonce []byte) (plaintext []byte, err error) {
	block, err := aes.NewCipher(dek)
	if err != nil {
		return nil, fmt.Errorf("creating cipher with DEK: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("creating GCM with DEK: %w", err)
	}

	plaintext, err = aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decrypting with DEK: %w", err)
	}

	return plaintext, nil
}

