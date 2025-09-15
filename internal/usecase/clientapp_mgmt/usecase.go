package clientapp_mgmt

import (
	"github.com/mandacode-com/mandacode-service-hub/internal/port/in"
	"github.com/mandacode-com/mandacode-service-hub/internal/port/out"
)

type Usecase struct {
	clientAppRepo      out.ClientAppRepository
	clientAppQueryRepo out.ClientAppQueryRepository
	serviceQueryRepo   out.ServiceQueryRepository
	txManager          out.TransactionManager
	secretGen          out.ByteRandGen
	hasher             out.Hasher
	encoder            out.Encoder
}

func NewUsecase(
	clientAppRepo out.ClientAppRepository,
	clientAppQueryRepo out.ClientAppQueryRepository,
	serviceQueryRepo out.ServiceQueryRepository,
	txManager out.TransactionManager,
	secretGen out.ByteRandGen,
	hasher out.Hasher,
	encoder out.Encoder,
) in.ClientAppMgmtUsecase {
	return &Usecase{
		clientAppRepo:      clientAppRepo,
		clientAppQueryRepo: clientAppQueryRepo,
		serviceQueryRepo:   serviceQueryRepo,
		txManager:          txManager,
		secretGen:          secretGen,
		hasher:             hasher,
		encoder:            encoder,
	}
}
