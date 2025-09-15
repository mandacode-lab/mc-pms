package clientapp_mgmt

import (
	"context"

	clientappval "github.com/mandacode-com/mandacode-ssam/internal/domain/clientapp/value"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
	"github.com/mandacode-com/merr"
)

func (u *Usecase) DeleteClientApp(ctx context.Context, req *in.DeleteClientAppRequest) error {
	// Parse client app ID
	clientAppID, err := clientappval.ParsePublicID(req.ClientAppID)
	if err != nil {
		return merr.New(merr.ErrBadRequest, ErrInvalidClientAppIDMsg, err)
	}

	// Find the client app by public ID
	clientApp, err := u.clientAppQueryRepo.FindByPublicID(ctx, clientAppID)
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
