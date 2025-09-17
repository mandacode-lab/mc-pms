package client_access

import (
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
)

type Usecase struct {
	clientAppQueryRepo out.ClientAppQueryRepository
	serviceQueryRepo   out.ServiceQueryRepository
	hasher             out.Hasher
}

func NewUsecase(
	clientAppQueryRepo out.ClientAppQueryRepository,
	serviceQueryRepo out.ServiceQueryRepository,
	hasher out.Hasher,
) in.ClientAccessUsecase {
	return &Usecase{
		clientAppQueryRepo: clientAppQueryRepo,
		serviceQueryRepo:   serviceQueryRepo,
		hasher:             hasher,
	}
}
