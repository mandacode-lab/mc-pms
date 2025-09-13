package identify_user

import (
	"github.com/mandacode-com/serengeti-integrated/internal/domain/shared"
	"github.com/mandacode-com/serengeti-integrated/internal/port/in"
	"github.com/mandacode-com/serengeti-integrated/internal/port/out"
)

type Usecase struct {
	oauthProviders        map[shared.Provider]out.OAuthAPI
	clientAppQueryRepo    out.ClientAppQueryRepository
	weboauthQueryRepo     out.WebOAuthQueryRepository
	userIdentityRepo      out.UserIdentityRepository
	userIdentityQueryRepo out.UserIdentityQueryRepository
	userInfoRepo          out.UserInfoRepository
	userInfoQueryRepo     out.UserInfoQueryRepository
	serviceQueryRepo      out.ServiceQueryRepository
	txManager             out.TransactionManager
	hasher                out.Hasher
	kekProvider           out.KekProvider
	stateService          out.StateService
}

func NewUsecase(
	oauthProviders map[shared.Provider]out.OAuthAPI,
	clientAppQueryRepo out.ClientAppQueryRepository,
	weboauthQueryRepo out.WebOAuthQueryRepository,
	userIdentityRepo out.UserIdentityRepository,
	userIdentityQueryRepo out.UserIdentityQueryRepository,
	userInfoRepo out.UserInfoRepository,
	userInfoQueryRepo out.UserInfoQueryRepository,
	serviceQueryRepo out.ServiceQueryRepository,
	txManager out.TransactionManager,
	hasher out.Hasher,
	kekProvider out.KekProvider,
	stateService out.StateService,
) in.IdentifyUserUsecase {
	return &Usecase{
		oauthProviders:        oauthProviders,
		clientAppQueryRepo:    clientAppQueryRepo,
		weboauthQueryRepo:     weboauthQueryRepo,
		userIdentityRepo:      userIdentityRepo,
		userIdentityQueryRepo: userIdentityQueryRepo,
		userInfoRepo:          userInfoRepo,
		userInfoQueryRepo:     userInfoQueryRepo,
		serviceQueryRepo:      serviceQueryRepo,
		txManager:             txManager,
		hasher:                hasher,
		kekProvider:           kekProvider,
		stateService:          stateService,
	}
}
