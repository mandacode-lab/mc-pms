package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/mandacode-com/mandacode-ssam/ent"
	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
)

func asEntTx(tx out.Tx) (*ent.Tx, error) {
	if tx == nil {
		return nil, errors.New("transaction is nil")
	}

	if entTx, ok := tx.(*ent.Tx); ok {
		return entTx, nil
	}
	return nil, fmt.Errorf("transaction type must be *ent.Tx")
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
			if tx != nil {
				_ = tx.Rollback() // Best effort rollback
			}
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		if tx != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return fmt.Errorf("transaction rollback error: %w (original error: %v)", rbErr, err)
			}
		}
		return err
	}

	if tx == nil {
		return errors.New("transaction is nil")
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
