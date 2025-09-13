package weboauth_mgmt

import (
	"github.com/mandacode-com/serengeti-integrated/internal/port/in"
	"github.com/mandacode-com/serengeti-integrated/internal/port/out"
)

type Usecase struct {
	weboauthRepo       out.WebOAuthRepository
	weboauthQueryRepo  out.WebOAuthQueryRepository
	clientAppQueryRepo out.ClientAppQueryRepository
	txManager          out.TransactionManager
	kekProvider        out.KekProvider
}

func NewUsecase(
	weboauthRepo out.WebOAuthRepository,
	weboauthQueryRepo out.WebOAuthQueryRepository,
	clientAppQueryRepo out.ClientAppQueryRepository,
	txManager out.TransactionManager,
	kekProvider out.KekProvider,
) in.WebOAuthMgmtUsecase {
	return &Usecase{
		weboauthRepo:       weboauthRepo,
		weboauthQueryRepo:  weboauthQueryRepo,
		clientAppQueryRepo: clientAppQueryRepo,
		txManager:          txManager,
		kekProvider:        kekProvider,
	}
}
