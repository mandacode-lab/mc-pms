package weboauth_mgmt

import (
	"context"

	"github.com/mandacode-com/merr"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
)

func (u *Usecase) UpdateWebOAuth(ctx context.Context, cmd *in.UpdateWebOAuthCommand) error {
	clientApp, err := u.clientAppQueryRepo.FindByPublicID(ctx, cmd.ClientAppID)
	if err != nil {
		return merr.New(merr.ErrNotFound, ErrClientAppNotFoundMsg, err)
	}

	webOAuth, err := u.weboauthQueryRepo.FindByProvider(ctx, clientApp.ID(), cmd.Provider)
	if err != nil {
		return merr.New(merr.ErrNotFound, ErrWebOAuthNotFoundMsg, err)
	}

	return u.txManager.WithTx(ctx, func(tx out.Tx) error {
		if cmd.OAuthClientID != nil {
			webOAuth.UpdateClientID(*cmd.OAuthClientID)
		}

		if cmd.RedirectURI != nil {
			webOAuth.UpdateRedirectURI(*cmd.RedirectURI)
		}

		if cmd.Scopes != nil {
			webOAuth.UpdateScopes(*cmd.Scopes)
		}

		if cmd.OAuthSecret != nil {
			dek, err := u.kekProvider.UnwrapDEK(ctx, webOAuth.DEKWrapped(), webOAuth.DEKNonce())
			if err != nil {
				return merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
			}

			secretCT, secretNonce, err := u.kekProvider.EncryptWithDEK(ctx, dek, *cmd.OAuthSecret)
			if err != nil {
				return merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
			}

			webOAuth.UpdateSecret(secretCT, secretNonce)
		}

		err := u.weboauthRepo.Update(ctx, tx, webOAuth)
		if err != nil {
			return merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
		}

		return nil
	})
}
