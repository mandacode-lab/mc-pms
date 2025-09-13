package identify_user

import (
	"context"

	"github.com/mandacode-com/merr"
	"github.com/mandacode-com/serengeti-integrated/internal/port/in"
)

func (u *Usecase) IdentifyByCode(ctx context.Context, cmd *in.IdentifyByCode) (*in.UserIdentityView, error) {
	clientApp, err := u.clientAppQueryRepo.FindByPublicID(ctx, cmd.ClientAppID)
	if err != nil {
		return nil, merr.New(merr.ErrNotFound, ErrClientAppNotFoundMsg, err)
	}

	hashedSecret, err := u.hasher.Hash(cmd.ClientAppSecret)
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

	dek, err := u.kekProvider.UnwrapDEK(ctx, weboauth.DEKWrapped(), weboauth.DEKNonce())
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}

	oauthSecret, err := u.kekProvider.DecryptWithDEK(ctx, dek, weboauth.OAuthSecretCT(), weboauth.OAuthSecretNonce())
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}

	accessToken, err := oauthProvider.GetAccessToken(cmd.OAuthCode, weboauth.OAuthClientID(), oauthSecret)
	if err != nil {
		return nil, merr.New(merr.ErrUnauthorized, ErrInvalidOAuthCodeMsg, err)
	}

	oauthUserInfo, err := oauthProvider.GetUserInfo(accessToken)
	if err != nil {
		return nil, merr.New(merr.ErrUnauthorized, ErrInvalidTokenMsg, err)
	}

	userIdentity, err := u.findOrCreateUser(ctx, clientApp.ServiceID(), cmd.Provider, oauthUserInfo)
	if err != nil {
		return nil, err
	}

	rawData := u.extractRawData(oauthUserInfo)

	return u.buildUserIdentityView(ctx, userIdentity, clientApp.ServiceID(), rawData)
}
