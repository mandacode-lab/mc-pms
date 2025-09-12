package in

import (
	serviceval "github.com/mandacode-com/serengeti-integrated/internal/domain/service/value"
	"github.com/mandacode-com/serengeti-integrated/internal/domain/shared"
	useridentityval "github.com/mandacode-com/serengeti-integrated/internal/domain/useridentity/value"
)

type FindUserInfoQuery struct {
	ServiceID *serviceval.PublicID
	UserID    *useridentityval.PublicID
	Provider  *shared.Provider
	Email     *string
	Nickname  *string
}

type FindUserInfoView struct {
	Users []UserInfo
}

type DeleteUserCommand struct {
	UserID useridentityval.PublicID
}

type UserMgmtUsecase interface {
	FindUserInfo(cmd FindUserInfoQuery) (FindUserInfoView, error)
	DeleteUser(cmd DeleteUserCommand) error
}
