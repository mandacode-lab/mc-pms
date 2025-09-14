package identify_user

import (
	"context"
	"encoding/json"

	"github.com/mandacode-com/merr"
	serviceval "github.com/mandacode-com/serengeti/internal/domain/service/value"
	"github.com/mandacode-com/serengeti/internal/domain/shared"
	"github.com/mandacode-com/serengeti/internal/domain/useridentity"
	"github.com/mandacode-com/serengeti/internal/domain/userinfo"
	"github.com/mandacode-com/serengeti/internal/port/in"
	"github.com/mandacode-com/serengeti/internal/port/out"
)

func (u *Usecase) findOrCreateUser(ctx context.Context, serviceID serviceval.ID, provider shared.Provider, oauthUserInfo *out.OAuthUserInfo) (*useridentity.UserIdentity, error) {
	userIdentity, err := u.userIdentityQueryRepo.FindByProviderID(ctx, serviceID, oauthUserInfo.ProviderID)
	if err != nil {
		err = u.txManager.WithTx(ctx, func(tx out.Tx) error {
			newUserIdentity := useridentity.DraftUserIdentity(
				serviceID,
				oauthUserInfo.ProviderID,
				provider,
			)

			userIdentity, err = u.userIdentityRepo.Create(ctx, tx, newUserIdentity)
			if err != nil {
				return merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
			}

			nickname := oauthUserInfo.ProviderID
			if oauthUserInfo.Nickname != "" {
				nickname = oauthUserInfo.Nickname
			} else if oauthUserInfo.Email != "" {
				nickname = oauthUserInfo.Email
			}

			newUserInfo := userinfo.DraftUserInfo(
				userIdentity.ID(),
				nickname,
				oauthUserInfo.Email,
				oauthUserInfo.RawData,
			)

			_, err = u.userInfoRepo.Create(ctx, tx, newUserInfo)
			if err != nil {
				return merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
			}

			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return userIdentity, nil
}

func (u *Usecase) buildUserIdentityView(ctx context.Context, userIdentity *useridentity.UserIdentity, serviceID serviceval.ID, rawData map[string]any) (*in.UserIdentityView, error) {
	userInfo, err := u.userInfoQueryRepo.FindByUserIdentityID(ctx, userIdentity.ID())
	if err != nil {
		return nil, merr.New(merr.ErrNotFound, ErrInternalServerMsg, err)
	}

	service, err := u.serviceQueryRepo.FindByID(ctx, serviceID)
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}

	return toUserIdentityView(userIdentity, userInfo, service.PublicID(), rawData), nil
}

func (u *Usecase) extractRawData(oauthUserInfo *out.OAuthUserInfo) map[string]any {
	var rawData map[string]any
	if len(oauthUserInfo.RawData) > 0 {
		if err := json.Unmarshal(oauthUserInfo.RawData, &rawData); err != nil {
			rawData = make(map[string]any)
		}
	} else {
		rawData = make(map[string]any)
	}
	return rawData
}
