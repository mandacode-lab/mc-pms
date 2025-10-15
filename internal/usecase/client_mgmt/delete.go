package clientmgmt

import (
	"context"

	"github.com/mandacode-com/mandacode-ssam/internal/domain/tx"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/merr"
)

func (u *Usecase) DeleteClient(ctx context.Context, req *in.DeleteClientInput) error {
	// Find the svcClient
	svcClient, err := u.clientQueryRepo.FindByPublicID(ctx, req.ClientID)
	if err != nil {
		return merr.New(merr.ErrNotFound, ErrClientNotFoundMsg, err)
	}

	// Delete within transaction
	err = u.txManager.WithTx(ctx, func(tx tx.Tx) error {
		return u.clientRepo.Delete(ctx, tx, svcClient)
	})
	if err != nil {
		return merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}

	return nil
}
