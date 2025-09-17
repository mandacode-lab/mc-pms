package initialize

import (
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
)

type Usecase struct {
	serviceRepo        out.ServiceRepository
	serviceQueryRepo   out.ServiceQueryRepository
	clientAppRepo      out.ClientAppRepository
	clientAppQueryRepo out.ClientAppQueryRepository
	txManager          out.TransactionManager
	secretGen          out.ByteRandGen
	hasher             out.Hasher
}

func NewUsecase(
	serviceRepo out.ServiceRepository,
	serviceQueryRepo out.ServiceQueryRepository,
	clientAppRepo out.ClientAppRepository,
	clientAppQueryRepo out.ClientAppQueryRepository,
	txManager out.TransactionManager,
	secretGen out.ByteRandGen,
	hasher out.Hasher,
) in.InitializeUsecase {
	return &Usecase{
		serviceRepo:        serviceRepo,
		serviceQueryRepo:   serviceQueryRepo,
		clientAppRepo:      clientAppRepo,
		clientAppQueryRepo: clientAppQueryRepo,
		txManager:          txManager,
		secretGen:          secretGen,
		hasher:             hasher,
	}
}