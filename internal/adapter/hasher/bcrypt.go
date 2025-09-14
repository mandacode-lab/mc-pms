package hasher

import (
	"context"

	"github.com/mandacode-com/serengeti-integrated/internal/port/out"
	"golang.org/x/crypto/bcrypt"
)

const DefaultCost = 12

type BcryptHasher struct {
	cost int
}

func NewBcryptHasher() out.Hasher {
	return &BcryptHasher{
		cost: DefaultCost,
	}
}

func NewBcryptHasherWithCost(cost int) out.Hasher {
	return &BcryptHasher{
		cost: cost,
	}
}

func (h *BcryptHasher) Hash(ctx context.Context, data []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(data, h.cost)
}

func (h *BcryptHasher) Compare(ctx context.Context, hash, data []byte) error {
	return bcrypt.CompareHashAndPassword(hash, data)
}
