package user_mgmt

import (
	"context"

	serviceval "github.com/mandacode-com/serengeti/internal/domain/service/value"
	"github.com/mandacode-com/serengeti/internal/domain/useridentity"
	"github.com/mandacode-com/serengeti/internal/domain/userinfo"
	"github.com/mandacode-com/serengeti/internal/port/in"
	"github.com/mandacode-com/serengeti/internal/port/out"
)

func toUserInfo(ui *userinfo.UserInfo, serviceID serviceval.PublicID, userIdentity *useridentity.UserIdentity) in.UserInfo {
	return in.UserInfo{
		ServiceID: serviceID,
		UserID:    userIdentity.PublicID(),
		Nickname:  ui.Nickname(),
		Email:     ui.Email(),
		Provider:  userIdentity.Provider(),
		CreatedAt: ui.CreatedAt(),
		UpdatedAt: ui.UpdatedAt(),
	}
}

func toUserInfos(userInfos []*userinfo.UserInfo, userIdentityQueryRepo out.UserIdentityQueryRepository, serviceQueryRepo out.ServiceQueryRepository, ctx context.Context) []in.UserInfo {
	users := make([]in.UserInfo, 0, len(userInfos))

	for _, ui := range userInfos {
		userIdentity, err := userIdentityQueryRepo.FindByID(ctx, ui.UserIdentityID())
		if err != nil {
			continue
		}

		service, err := serviceQueryRepo.FindByID(ctx, userIdentity.ServiceID())
		if err != nil {
			continue
		}

		users = append(users, toUserInfo(ui, service.PublicID(), userIdentity))
	}

	return users
}

