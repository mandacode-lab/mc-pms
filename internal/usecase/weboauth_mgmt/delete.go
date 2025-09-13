package weboauth_mgmt

import (
	"context"

	"github.com/mandacode-com/merr"
	"github.com/mandacode-com/serengeti-integrated/internal/port/in"
	"github.com/mandacode-com/serengeti-integrated/internal/port/out"
)

func (u *Usecase) DeleteWebOAuth(ctx context.Context, cmd *in.DeleteWebOAuthCommand) error {
	clientApp, err := u.clientAppQueryRepo.FindByPublicID(ctx, cmd.ClientAppID)
	if err != nil {
		return merr.New(merr.ErrNotFound, ErrClientAppNotFoundMsg, err)
	}

	webOAuth, err := u.weboauthQueryRepo.FindByProvider(ctx, clientApp.ID(), cmd.Provider)
	if err != nil {
		return merr.New(merr.ErrNotFound, ErrWebOAuthNotFoundMsg, err)
	}

	return u.txManager.WithTx(ctx, func(tx out.Tx) error {
		err := u.weboauthRepo.Delete(ctx, tx, webOAuth)
		if err != nil {
			return merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
		}
		return nil
	})
}
