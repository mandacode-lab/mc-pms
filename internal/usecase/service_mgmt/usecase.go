package servicemgmt

import (
	"github.com/mandacode-com/mandacode-ssam/internal/domain/service"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/tx"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
)

type Usecase struct {
	serviceRepo      service.ServiceRepository
	serviceQueryRepo service.ServiceQueryRepository
	serviceIDGen     out.IDGenerator
	txManager        tx.TxManager
}

func NewUsecase(
	serviceRepo service.ServiceRepository,
	serviceQueryRepo service.ServiceQueryRepository,
	serviceIDGen out.IDGenerator,
	txManager tx.TxManager,
) in.ServiceMgmtUsecase {
	return &Usecase{
		serviceRepo:      serviceRepo,
		serviceQueryRepo: serviceQueryRepo,
		serviceIDGen:     serviceIDGen,
		txManager:        txManager,
	}
}
