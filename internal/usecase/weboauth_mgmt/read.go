package weboauth_mgmt

import (
	"context"

	"github.com/mandacode-com/merr"
	"github.com/mandacode-com/serengeti/internal/domain/weboauth"
	"github.com/mandacode-com/serengeti/internal/port/in"
)

func (u *Usecase) ReadWebOAuth(ctx context.Context, query *in.ReadWebOAuthQuery) (*in.ReadWebOAuthView, error) {
	clientApp, err := u.clientAppQueryRepo.FindByPublicID(ctx, query.ClientAppID)
	if err != nil {
		return nil, merr.New(merr.ErrNotFound, ErrClientAppNotFoundMsg, err)
	}

	var webOAuths []*weboauth.WebOAuth
	if query.Provider != nil {
		webOAuth, err := u.weboauthQueryRepo.FindByProvider(ctx, clientApp.ID(), *query.Provider)
		if err != nil {
			return nil, merr.New(merr.ErrNotFound, ErrWebOAuthNotFoundMsg, err)
		}
		webOAuths = []*weboauth.WebOAuth{webOAuth}
	} else {
		webOAuths, err = u.weboauthQueryRepo.FindByClientAppID(ctx, clientApp.ID())
		if err != nil {
			return nil, merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
		}
	}

	webOAuthInfos := make([]in.WebOAuthInfo, 0, len(webOAuths))
	for _, wo := range webOAuths {
		webOAuthInfos = append(webOAuthInfos, toWebOAuthInfo(wo, clientApp.PublicID()))
	}

	return &in.ReadWebOAuthView{
		WebOAuths: webOAuthInfos,
	}, nil
}

func (u *Usecase) ReadWebOAuthSecret(ctx context.Context, query *in.ReadWebOAuthSecretQuery) (*in.ReadWebOAuthSecretView, error) {
	clientApp, err := u.clientAppQueryRepo.FindByPublicID(ctx, query.ClientAppID)
	if err != nil {
		return nil, merr.New(merr.ErrNotFound, ErrClientAppNotFoundMsg, err)
	}

	webOAuth, err := u.weboauthQueryRepo.FindByProvider(ctx, clientApp.ID(), query.Provider)
	if err != nil {
		return nil, merr.New(merr.ErrNotFound, ErrWebOAuthNotFoundMsg, err)
	}

	dek, err := u.kekProvider.UnwrapDEK(ctx, webOAuth.DEKWrapped(), webOAuth.DEKNonce())
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}

	oauthSecret, err := u.kekProvider.DecryptWithDEK(ctx, dek, webOAuth.OAuthSecretCT(), webOAuth.OAuthSecretNonce())
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}

	return &in.ReadWebOAuthSecretView{
		OAuthSecret: oauthSecret,
	}, nil
}
