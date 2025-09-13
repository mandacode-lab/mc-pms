package out

import "context"

type Hasher interface {
	Hash(ctx context.Context, data []byte) ([]byte, error)
	Compare(ctx context.Context, hash, data []byte) error
}
