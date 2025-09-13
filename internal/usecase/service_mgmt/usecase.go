package service_mgmt

import (
	"github.com/mandacode-com/serengeti-integrated/internal/port/in"
	"github.com/mandacode-com/serengeti-integrated/internal/port/out"
)

type Usecase struct {
	serviceRepo      out.ServiceRepository
	serviceQueryRepo out.ServiceQueryRepository
	txManager        out.TransactionManager
}

func NewUsecase(
	serviceRepo out.ServiceRepository,
	serviceQueryRepo out.ServiceQueryRepository,
	txManager out.TransactionManager,
) in.ServiceMgmtUsecase {
	return &Usecase{
		serviceRepo:      serviceRepo,
		serviceQueryRepo: serviceQueryRepo,
		txManager:        txManager,
	}
}
