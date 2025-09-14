package random

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/mandacode-com/serengeti/internal/port/out"
)

const (
	DefaultByteSize   = 32
	DefaultStringSize = 32
)

type CryptoByteRandGen struct {
	defaultSize int
}

func NewCryptoByteRandGen() out.ByteRandGen {
	return &CryptoByteRandGen{
		defaultSize: DefaultByteSize,
	}
}

func (g *CryptoByteRandGen) Generate(ctx context.Context) ([]byte, error) {
	return g.GenerateN(ctx, g.defaultSize)
}

func (g *CryptoByteRandGen) GenerateN(ctx context.Context, n int) ([]byte, error) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, fmt.Errorf("generating random bytes: %w", err)
	}
	return bytes, nil
}

type CryptoStrRandGen struct {
	defaultSize int
}

func NewCryptoStrRandGen() out.StrRandGen {
	return &CryptoStrRandGen{
		defaultSize: DefaultStringSize,
	}
}

func (g *CryptoStrRandGen) Generate(ctx context.Context) (string, error) {
	return g.GenerateN(ctx, g.defaultSize)
}

func (g *CryptoStrRandGen) GenerateN(ctx context.Context, n int) (string, error) {
	// Generate enough bytes to get the desired string length
	// Base64 encoding expands by ~33%, so we need fewer bytes
	byteSize := (n*3 + 3) / 4 // Calculate required bytes for base64

	bytes := make([]byte, byteSize)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("generating random bytes: %w", err)
	}

	// Use URL-safe base64 and trim to exact length
	result := base64.URLEncoding.EncodeToString(bytes)
	if len(result) > n {
		result = result[:n]
	}

	return result, nil
}
