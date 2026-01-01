package idgen

import (
	"crypto/rand"
	"fmt"

	"github.com/oklog/ulid/v2"
)

// ULIDGenerator generates ULID-based unique identifiers
type ULIDGenerator struct {
	entropy *ulid.MonotonicEntropy
}

func NewULIDGenerator() *ULIDGenerator {
	entropy := ulid.Monotonic(rand.Reader, 0)
	return &ULIDGenerator{
		entropy: entropy,
	}
}

// Generate creates a new ULID
func (g *ULIDGenerator) Generate() (string, error) {
	id, err := ulid.New(ulid.Now(), g.entropy)
	if err != nil {
		return "", fmt.Errorf("failed to generate ULID: %w", err)
	}
	return id.String(), nil
}

// GenerateWithPrefix creates a new ULID with a prefix
func (g *ULIDGenerator) GenerateWithPrefix(prefix string) (string, error) {
	id, err := g.Generate()
	if err != nil {
		return "", err
	}
	return prefix + "_" + id, nil
}
