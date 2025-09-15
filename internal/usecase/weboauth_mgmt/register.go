package weboauth_mgmt

import (
	"context"

	"github.com/mandacode-com/merr"
	"github.com/mandacode-com/mandacode-service-hub/internal/domain/weboauth"
	"github.com/mandacode-com/mandacode-service-hub/internal/port/in"
	"github.com/mandacode-com/mandacode-service-hub/internal/port/out"
)

func (u *Usecase) RegisterWebOAuth(ctx context.Context, cmd *in.RegisterWebOAuthCommand) error {
	clientApp, err := u.clientAppQueryRepo.FindByPublicID(ctx, cmd.ClientAppID)
	if err != nil {
		return merr.New(merr.ErrNotFound, ErrClientAppNotFoundMsg, err)
	}

	_, err = u.weboauthQueryRepo.FindByProvider(ctx, clientApp.ID(), cmd.Provider)
	if err == nil {
		return merr.New(merr.ErrConflict, ErrWebOAuthAlreadyExistsMsg, nil)
	}

	wrappedDEK, dekNonce, err := u.kekProvider.GenerateDEK(ctx)
	if err != nil {
		return merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}

	dek, err := u.kekProvider.UnwrapDEK(ctx, wrappedDEK, dekNonce)
	if err != nil {
		return merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}

	secretCT, secretNonce, err := u.kekProvider.EncryptWithDEK(ctx, dek, cmd.OAuthSecret)
	if err != nil {
		return merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}

	return u.txManager.WithTx(ctx, func(tx out.Tx) error {
		newWebOAuth := weboauth.DraftWebOAuth(
			clientApp.ID(),
			cmd.Provider,
			cmd.OAuthClientID,
			secretCT,
			secretNonce,
			wrappedDEK,
			dekNonce,
			cmd.RedirectURI,
			cmd.Scopes,
		)

		_, err := u.weboauthRepo.Create(ctx, tx, newWebOAuth)
		if err != nil {
			return merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
		}

		return nil
	})
}
