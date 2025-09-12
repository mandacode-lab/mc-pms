package out

import (
	"context"

	"github.com/mandacode-com/serengeti-integrated/internal/domain/userinfo"
	userinfoval "github.com/mandacode-com/serengeti-integrated/internal/domain/userinfo/value"
	useridentityval "github.com/mandacode-com/serengeti-integrated/internal/domain/useridentity/value"
)

type UserInfoListFilter struct {
	UserIdentityID *useridentityval.ID
	Email          *string
	Nickname       *string
}

type UserInfoListOrder string

const (
	UserInfoListOrderNicknameAsc  UserInfoListOrder = "nickname_asc"
	UserInfoListOrderNicknameDesc UserInfoListOrder = "nickname_desc"
)

type UserInfoListOptions struct {
	Limit  int
	Offset int
	Order  UserInfoListOrder
}


type UserInfoRepository interface {
	Create(ctx context.Context, tx Tx, info *userinfo.UserInfo) (*userinfo.UserInfo, error)
	Update(ctx context.Context, tx Tx, info *userinfo.UserInfo) error
	FindByID(ctx context.Context, id userinfoval.ID) (*userinfo.UserInfo, error)
	Delete(ctx context.Context, tx Tx, info *userinfo.UserInfo) error
}

type UserInfoQueryRepository interface {
	FindByID(ctx context.Context, id userinfoval.ID) (*userinfo.UserInfo, error)
	FindByPublicID(ctx context.Context, publicID userinfoval.PublicID) (*userinfo.UserInfo, error)
	FindByUserIdentityID(ctx context.Context, userIdentityID useridentityval.ID) (*userinfo.UserInfo, error)
	FindByEmail(ctx context.Context, email string) (*userinfo.UserInfo, error)
	FindByNickname(ctx context.Context, nickname string) (*userinfo.UserInfo, error)
	List(ctx context.Context, filter *UserInfoListFilter, options *UserInfoListOptions) ([]*userinfo.UserInfo, int, error)
}