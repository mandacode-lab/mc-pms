package main

import (
	"context"
	"database/sql"
	"fmt"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"

	_ "github.com/lib/pq"
	"github.com/mandacode-com/mandacode-ssam/configs"
	"github.com/mandacode-com/mandacode-ssam/ent"
	"github.com/mandacode-com/mandacode-ssam/internal/adapter/hasher"
	"github.com/mandacode-com/mandacode-ssam/internal/adapter/random"
	"github.com/mandacode-com/mandacode-ssam/internal/adapter/repository"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
	"github.com/mandacode-com/mandacode-ssam/internal/usecase/initialize"
)

type Adapter struct {
	// Database
	db     *sql.DB
	client *ent.Client

	// Repositories
	serviceRepo        out.ServiceRepository
	serviceQueryRepo   out.ServiceQueryRepository
	clientAppRepo      out.ClientAppRepository
	clientAppQueryRepo out.ClientAppQueryRepository

	// Services
	hasher    out.Hasher
	secretGen out.ByteRandGen
	txManager out.TransactionManager

	// Usecases
	initializeUsecase in.InitializeUsecase
}

func NewAdapter(ctx context.Context, cfg *configs.InitializeConfig) (*Adapter, error) {
	// Database setup
	db, err := sql.Open("postgres", cfg.Postgres.GetDSN())
	if err != nil {
		return nil, err
	}

	drv := entsql.OpenDB(dialect.Postgres, db)
	client := ent.NewClient(ent.Driver(drv))

	// Repositories
	serviceRepo := repository.NewEntServiceRepository(client)
	serviceQueryRepo := repository.NewEntServiceQueryRepository(client)
	clientAppRepo := repository.NewEntClientAppRepository(client)
	clientAppQueryRepo := repository.NewEntClientAppQueryRepository(client)

	// Services
	hasherSvc := hasher.NewBcryptHasher()
	secretGen := random.NewCryptoByteRandGen()
	txManager := repository.NewEntTransactionManager(client)

	// Initialize usecase
	initializeUsecase := initialize.NewUsecase(
		serviceRepo,
		serviceQueryRepo,
		clientAppRepo,
		clientAppQueryRepo,
		txManager,
		secretGen,
		hasherSvc,
	)

	return &Adapter{
		db:                db,
		client:            client,
		serviceRepo:       serviceRepo,
		serviceQueryRepo:  serviceQueryRepo,
		clientAppRepo:     clientAppRepo,
		clientAppQueryRepo: clientAppQueryRepo,
		hasher:            hasherSvc,
		secretGen:         secretGen,
		txManager:         txManager,
		initializeUsecase: initializeUsecase,
	}, nil
}

func (a *Adapter) ProvideInitializeUsecase() in.InitializeUsecase {
	return a.initializeUsecase
}

func (a *Adapter) Close() error {
	if a.client != nil {
		if err := a.client.Close(); err != nil {
			return fmt.Errorf("failed to close ent client: %w", err)
		}
	}
	if a.db != nil {
		if err := a.db.Close(); err != nil {
			return fmt.Errorf("failed to close database: %w", err)
		}
	}
	return nil
}