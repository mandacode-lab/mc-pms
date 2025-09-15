package user_mgmt

import (
	"context"

	"github.com/mandacode-com/merr"
	"github.com/mandacode-com/mandacode-service-hub/internal/port/in"
	"github.com/mandacode-com/mandacode-service-hub/internal/port/out"
)

func (u *Usecase) FindUserInfo(ctx context.Context, query *in.FindUserInfoQuery) (*in.FindUserInfoResult, error) {
	// Build filter from query parameters
	filter := &out.UserInfoListFilter{}
	if query.Email != nil {
		filter.Email = query.Email
	}

	// If UserID is provided, we need to find the UserIdentity first to get the internal ID
	if query.UserID != nil {
		userIdentity, err := u.userIdentityQueryRepo.FindByPublicID(ctx, *query.UserID)
		if err != nil {
			return nil, merr.New(merr.ErrNotFound, ErrUserNotFoundMsg, err)
		}
		userIdentityID := userIdentity.ID()
		filter.UserIdentityID = &userIdentityID
	}

	// Query user info with pagination (default limit)
	options := &out.UserInfoListOptions{
		Limit:  100,
		Offset: 0,
		Order:  out.UserInfoListOrderNicknameAsc,
	}

	userInfos, _, err := u.userInfoQueryRepo.List(ctx, filter, options)
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}

	// Convert to result models
	users := toUserInfos(userInfos, u.userIdentityQueryRepo, u.serviceQueryRepo, ctx)

	return &in.FindUserInfoResult{
		Users: users,
	}, nil
}

