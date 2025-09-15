package in

import (
	"context"
	
	serviceval "github.com/mandacode-com/mandacode-ssam/internal/domain/service/value"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/shared"
	useridentityval "github.com/mandacode-com/mandacode-ssam/internal/domain/useridentity/value"
)

type FindUserInfoQuery struct {
	ServiceID *serviceval.PublicID
	UserID    *useridentityval.PublicID
	Provider  *shared.Provider
	Email     *string
	Nickname  *string
}

type FindUserInfoResult struct {
	Users []UserInfo
}

type DeleteUserCommand struct {
	UserID useridentityval.PublicID
}

type UserMgmtUsecase interface {
	FindUserInfo(ctx context.Context, query *FindUserInfoQuery) (*FindUserInfoResult, error)
	DeleteUser(ctx context.Context, cmd *DeleteUserCommand) error
}
