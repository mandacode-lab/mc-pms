package out

import "context"

type StrRandGen interface {
	Generate(ctx context.Context) (string, error)
	GenerateN(ctx context.Context, n int) (string, error)
}
