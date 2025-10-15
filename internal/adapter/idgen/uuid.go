package idgen

import (
	"github.com/google/uuid"
	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
)

type UUIDGenerator struct{}

// Generate implements out.IDGenerator.
func (p *UUIDGenerator) Generate() (string, error) {
	id := uuid.New()
	return id.String(), nil
}

func NewUUIDGenerator() out.IDGenerator {
	return &UUIDGenerator{}
}
