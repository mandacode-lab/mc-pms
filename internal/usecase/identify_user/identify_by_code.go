package identify_user

import (
	"context"

	"github.com/mandacode-com/merr"
	"github.com/mandacode-com/mandacode-service-hub/internal/port/in"
)

func (u *Usecase) IdentifyByCode(ctx context.Context, cmd *in.IdentifyByCode) (*in.UserIdentityView, error) {
	// Validate state parameter to prevent CSRF attacks
	if cmd.State == "" {
		return nil, merr.New(merr.ErrBadRequest, ErrMissingStateMsg, nil)
	}

	if !u.stateService.ValidateState(ctx, cmd.State) {
		return nil, merr.New(merr.ErrUnauthorized, ErrInvalidStateMsg, nil)
	}

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

	dek, err := u.kekProvider.UnwrapDEK(ctx, weboauth.DEKWrapped(), weboauth.DEKNonce())
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}

	oauthSecret, err := u.kekProvider.DecryptWithDEK(ctx, dek, weboauth.OAuthSecretCT(), weboauth.OAuthSecretNonce())
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}

	accessToken, err := oauthProvider.GetAccessToken(ctx, cmd.OAuthCode, weboauth.OAuthClientID(), oauthSecret, weboauth.RedirectURI())
	if err != nil {
		return nil, merr.New(merr.ErrUnauthorized, ErrInvalidOAuthCodeMsg, err)
	}

	oauthUserInfo, err := oauthProvider.GetUserInfo(ctx, accessToken)
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
