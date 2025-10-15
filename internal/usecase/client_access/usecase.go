package clientaccess

import (
	"github.com/mandacode-com/mandacode-ssam/internal/domain/client"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/service"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
)

type Usecase struct {
	clientQueryRepo  client.ClientQueryRepository
	serviceQueryRepo service.ServiceQueryRepository
	hasher           out.Hasher
	encoder          out.Encoder
}

func NewUsecase(
	clientQueryRepo client.ClientQueryRepository,
	serviceQueryRepo service.ServiceQueryRepository,
	hasher out.Hasher,
	encoder out.Encoder,
) in.ClientAccessUsecase {
	return &Usecase{
		clientQueryRepo:  clientQueryRepo,
		serviceQueryRepo: serviceQueryRepo,
		hasher:           hasher,
		encoder:          encoder,
	}
}
