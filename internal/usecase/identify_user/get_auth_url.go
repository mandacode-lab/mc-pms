package identify_user

import (
	"context"

	"github.com/mandacode-com/merr"
	"github.com/mandacode-com/serengeti/internal/port/in"
)

func (u *Usecase) GetAuthURL(ctx context.Context, cmd *in.GetAuthURL) (*in.GetAuthURLView, error) {
	clientApp, err := u.clientAppQueryRepo.FindByPublicID(ctx, cmd.ClientAppID)
	if err != nil {
		return nil, merr.New(merr.ErrNotFound, ErrClientAppNotFoundMsg, err)
	}

	hashedSecret, err := u.hasher.Hash(ctx, cmd.ClientAppSecret)
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}

	if !clientApp.VerifySecretBytes(cmd.ClientAppSecret, hashedSecret) {
		return nil, merr.New(merr.ErrUnauthorized, ErrInvalidClientSecretMsg, nil)
	}

	weboauth, err := u.weboauthQueryRepo.FindByProvider(ctx, clientApp.ID(), cmd.Provider)
	if err != nil {
		return nil, merr.New(merr.ErrNotFound, ErrInvalidProviderMsg, err)
	}

	oauthProvider, exists := u.oauthProviders[cmd.Provider]
	if !exists {
		return nil, merr.New(merr.ErrBadRequest, ErrInvalidProviderMsg, nil)
	}

	state, err := u.stateService.GenerateState(ctx)
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}

	authURL := oauthProvider.GetAuthURL(ctx, weboauth.OAuthClientID(), weboauth.Scopes(), weboauth.RedirectURI(), state)

	return &in.GetAuthURLView{
		AuthURL: authURL,
		State:   state,
	}, nil
}
