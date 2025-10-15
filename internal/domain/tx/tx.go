package tx

import (
	"context"
)

type Tx interface {
	Rollback() error
	Commit() error
}

type TxManager interface {
	// WithTx executes the provided function within a database transaction.
	// If the function returns an error, the transaction is rolled back.
	// Otherwise, the transaction is committed.
	WithTx(ctx context.Context, fn func(tx Tx) error) error
}
