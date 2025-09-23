package repository

import (
	"context"
	"fmt"

	"github.com/mandacode-com/mandacode-ssam/ent"
	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
)

func asEntTx(tx out.Tx) (*ent.Tx, error) {
	entTx, ok := tx.(*ent.Tx)
	if !ok {
		return nil, fmt.Errorf("transaction type must be *ent.Tx")
	}
	return entTx, nil
}

type EntTransactionManager struct {
	client *ent.Client
}

func NewEntTransactionManager(client *ent.Client) out.TransactionManager {
	return &EntTransactionManager{
		client: client,
	}
}

func (tm *EntTransactionManager) WithTx(ctx context.Context, fn func(tx out.Tx) error) error {
	tx, err := tm.client.Tx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback() // Best effort rollback
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("failed to rollback transaction: %w", rollbackErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

