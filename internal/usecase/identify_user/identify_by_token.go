package identify_user

import (
	"context"

	"github.com/mandacode-com/merr"
	"github.com/mandacode-com/serengeti-integrated/internal/port/in"
)

func (u *Usecase) IdentifyByToken(ctx context.Context, cmd *in.IdentifyByToken) (*in.UserIdentityView, error) {
	clientApp, err := u.clientAppQueryRepo.FindByPublicID(ctx, cmd.ClientApp)
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

	oauthProvider, exists := u.oauthProviders[cmd.Provider]
	if !exists {
		return nil, merr.New(merr.ErrBadRequest, ErrInvalidProviderMsg, nil)
	}

	oauthUserInfo, err := oauthProvider.GetUserInfo(ctx, cmd.Token)
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
