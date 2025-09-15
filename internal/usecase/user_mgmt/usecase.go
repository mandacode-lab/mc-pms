package user_mgmt

import (
	"github.com/mandacode-com/mandacode-service-hub/internal/port/in"
	"github.com/mandacode-com/mandacode-service-hub/internal/port/out"
)

type Usecase struct {
	userInfoQueryRepo     out.UserInfoQueryRepository
	userIdentityQueryRepo out.UserIdentityQueryRepository
	userIdentityRepo      out.UserIdentityRepository
	userInfoRepo          out.UserInfoRepository
	serviceQueryRepo      out.ServiceQueryRepository
	txManager             out.TransactionManager
}

func NewUsecase(
	userInfoQueryRepo out.UserInfoQueryRepository,
	userIdentityQueryRepo out.UserIdentityQueryRepository,
	userIdentityRepo out.UserIdentityRepository,
	userInfoRepo out.UserInfoRepository,
	serviceQueryRepo out.ServiceQueryRepository,
	txManager out.TransactionManager,
) in.UserMgmtUsecase {
	return &Usecase{
		userInfoQueryRepo:     userInfoQueryRepo,
		userIdentityQueryRepo: userIdentityQueryRepo,
		userIdentityRepo:      userIdentityRepo,
		userInfoRepo:          userInfoRepo,
		serviceQueryRepo:      serviceQueryRepo,
		txManager:             txManager,
	}
}

