package clientmgmt

import (
	"github.com/mandacode-com/mandacode-ssam/internal/domain/client"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/service"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/tx"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
)

type Usecase struct {
	clientRepo       client.ClientRepository
	clientQueryRepo  client.ClientQueryRepository
	serviceQueryRepo service.ServiceQueryRepository
	clientIDGen      out.IDGenerator
	txManager        tx.TxManager
	secretGen        out.ByteRandGen
	hasher           out.Hasher
	encoder          out.Encoder
}

func NewUsecase(
	clientAppRepo client.ClientRepository,
	clientAppQueryRepo client.ClientQueryRepository,
	serviceQueryRepo service.ServiceQueryRepository,
	clientIDGen out.IDGenerator,
	txManager tx.TxManager,
	secretGen out.ByteRandGen,
	hasher out.Hasher,
	encoder out.Encoder,
) in.ClientMgmtUsecase {
	return &Usecase{
		clientRepo:       clientAppRepo,
		clientQueryRepo:  clientAppQueryRepo,
		serviceQueryRepo: serviceQueryRepo,
		clientIDGen:      clientIDGen,
		txManager:        txManager,
		secretGen:        secretGen,
		hasher:           hasher,
		encoder:          encoder,
	}
}
