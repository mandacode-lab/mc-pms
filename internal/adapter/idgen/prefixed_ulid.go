package idgen

import (
	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
	"github.com/oklog/ulid/v2"
)

type PrefixedULIDGenerator struct {
	prefix string
}

// Generate implements out.IDGenerator.
func (p *PrefixedULIDGenerator) Generate() (string, error) {
	id, err := ulid.New(ulid.Now(), ulid.Monotonic(nil, 0))
	if err != nil {
		return "", err
	}
	return p.prefix + id.String(), nil
}

func NewPrefixedULIDGenerator(prefix string) out.IDGenerator {
	return &PrefixedULIDGenerator{prefix: prefix}
}
