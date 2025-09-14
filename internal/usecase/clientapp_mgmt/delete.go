package clientapp_mgmt

import (
	"context"

	"github.com/mandacode-com/merr"
	"github.com/mandacode-com/serengeti/internal/port/in"
	"github.com/mandacode-com/serengeti/internal/port/out"
)

func (u *Usecase) DeleteClientApp(ctx context.Context, cmd *in.DeleteClientAppCommand) error {
	// Find the client app by public ID
	clientApp, err := u.clientAppQueryRepo.FindByPublicID(ctx, cmd.ClientAppID)
	if err != nil {
		return merr.New(merr.ErrNotFound, ErrClientAppNotFoundMsg, err)
	}

	// Delete within transaction
	err = u.txManager.WithTx(ctx, func(tx out.Tx) error {
		return u.clientAppRepo.Delete(ctx, tx, clientApp)
	})
	if err != nil {
		return merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}

	return nil
}

