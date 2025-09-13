package out

import "context"

type ByteRandGen interface {
	Generate(ctx context.Context) ([]byte, error)
	GenerateN(ctx context.Context, n int) ([]byte, error)
}
