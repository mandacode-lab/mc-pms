package out

import "context"

type StateService interface {
	GenerateState(ctx context.Context) (string, error)
	ValidateState(ctx context.Context, state string) bool
}
