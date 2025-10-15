package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"

	_ "github.com/lib/pq"
	"github.com/mandacode-com/mandacode-ssam/configs"
	"github.com/mandacode-com/mandacode-ssam/ent"
	"github.com/mandacode-com/mandacode-ssam/internal/adapter/cache"
	"github.com/mandacode-com/mandacode-ssam/internal/adapter/encoder"
	"github.com/mandacode-com/mandacode-ssam/internal/adapter/handler/client_mgmt"
	"github.com/mandacode-com/mandacode-ssam/internal/adapter/handler/service_mgmt"
	"github.com/mandacode-com/mandacode-ssam/internal/adapter/hasher"
	"github.com/mandacode-com/mandacode-ssam/internal/adapter/iam"
	"github.com/mandacode-com/mandacode-ssam/internal/adapter/idgen"
	"github.com/mandacode-com/mandacode-ssam/internal/adapter/random"
	"github.com/mandacode-com/mandacode-ssam/internal/adapter/repository"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/client"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/service"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/tx"
	redisinfra "github.com/mandacode-com/mandacode-ssam/internal/infra/redis"
	"github.com/mandacode-com/mandacode-ssam/internal/middleware"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
	client_usecase "github.com/mandacode-com/mandacode-ssam/internal/usecase/client_mgmt"
	service_usecase "github.com/mandacode-com/mandacode-ssam/internal/usecase/service_mgmt"
	"github.com/redis/go-redis/v9"
)

type Adapter struct {
	// Database
	db     *sql.DB
	client *ent.Client

	// Cache
	cacheClient redis.UniversalClient
	cacheStore  out.CacheStore

	// Repositories
	serviceRepo        service.ServiceRepository
	serviceQueryRepo   service.ServiceQueryRepository
	clientAppRepo      client.ClientRepository
	clientAppQueryRepo client.ClientQueryRepository
	txManager          tx.TxManager

	// Services
	hasher             out.Hasher
	secretGen          out.ByteRandGen
	encoder            out.Encoder
	iamService         out.IAMService
	serviceIDGenerator out.IDGenerator
	clientIDGenerator  out.IDGenerator

	// Middleware
	permissionMiddleware *middleware.PermissionMiddleware
}

func NewAdapter(ctx context.Context, cfg *configs.ManagementConfig) (*Adapter, error) {
	// Database setup
	db, err := sql.Open("postgres", cfg.Postgres.GetDSN())
	if err != nil {
		return nil, err
	}

	drv := entsql.OpenDB(dialect.Postgres, db)
	client := ent.NewClient(ent.Driver(drv))

	// Redis setup
	cacheClient, err := redisinfra.NewClient(ctx, redisinfra.Config{
		Addr:             cfg.Redis.Addr,
		Password:         cfg.Redis.Password,
		DB:               cfg.Redis.DB,
		Mode:             cfg.Redis.Mode,
		SentinelMaster:   cfg.Redis.SentinelMaster,
		SentinelPassword: cfg.Redis.SentinelPassword,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create redis client: %w", err)
	}

	// Cache setup
	cacheConfig := &out.CacheConfig{
		DefaultExpiration: 30 * time.Minute,
	}
	cacheStore := cache.NewRedisCacheStore(cacheClient, cacheConfig)

	// Repositories
	serviceRepo := repository.NewEntServiceRepository(client)
	serviceQueryRepo := repository.NewEntServiceQueryRepository(client)
	clientAppRepo := repository.NewEntClientAppRepository(client)
	clientAppQueryRepo := repository.NewEntClientAppQueryRepository(client)

	// Services
	hasherSvc := hasher.NewBcryptHasher()
	secretGen := random.NewCryptoByteRandGen()
	encoderSvc := encoder.NewBase64Encoder()
	txManager := repository.NewEntTransactionManager(client)
	serviceIDGenerator := idgen.NewPrefixedULIDGenerator(cfg.IDPrefix.Service)
	clientIDGenerator := idgen.NewUUIDGenerator()

	// IAM Service - choose based on configuration
	var iamSvc out.IAMService
	if cfg.IAM.Enabled {
		iamSvc, err = iam.NewHTTPIAMService(cfg.IAM.ServiceURL)
		if err != nil {
			return nil, fmt.Errorf("failed to create IAM service: %w", err)
		}
	} else {
		iamSvc = iam.NewNoOpIAMService()
	}

	// Middleware
	permissionMiddleware := middleware.NewPermissionMiddleware(iamSvc, cfg.IAM.Enabled)

	return &Adapter{
		// Database
		db:     db,
		client: client,

		// Cache
		cacheClient: cacheClient,
		cacheStore:  cacheStore,

		// Repositories
		serviceRepo:        serviceRepo,
		serviceQueryRepo:   serviceQueryRepo,
		clientAppRepo:      clientAppRepo,
		clientAppQueryRepo: clientAppQueryRepo,
		txManager:          txManager,

		// Services
		hasher:               hasherSvc,
		secretGen:            secretGen,
		encoder:              encoderSvc,
		iamService:           iamSvc,
		permissionMiddleware: permissionMiddleware,
		serviceIDGenerator:   serviceIDGenerator,
		clientIDGenerator:    clientIDGenerator,
	}, nil
}

func (a *Adapter) ProvideServiceMgmtUsecase() in.ServiceMgmtUsecase {
	return service_usecase.NewUsecase(
		a.serviceRepo,
		a.serviceQueryRepo,
		a.serviceIDGenerator,
		a.txManager,
	)
}

func (a *Adapter) ProvideServiceMgmtHandler() *servicemgmt.Handler {
	handler := servicemgmt.NewHandler(a.ProvideServiceMgmtUsecase())
	handler.SetPermissionMiddleware(a.permissionMiddleware)
	return handler
}

func (a *Adapter) ProvideClientAppMgmtUsecase() in.ClientMgmtUsecase {
	return client_usecase.NewUsecase(
		a.clientAppRepo,
		a.clientAppQueryRepo,
		a.serviceQueryRepo,
		a.clientIDGenerator,
		a.txManager,
		a.secretGen,
		a.hasher,
		a.encoder,
	)
}

func (a *Adapter) ProvideClientAppMgmtHandler() *clientmgmt.Handler {
	handler := clientmgmt.NewHandler(a.ProvideClientAppMgmtUsecase())
	handler.SetPermissionMiddleware(a.permissionMiddleware)
	return handler
}

func (a *Adapter) Close() error {
	var err error
	if a.client != nil {
		if closeErr := a.client.Close(); closeErr != nil {
			err = closeErr
		}
	}
	if a.db != nil {
		if closeErr := a.db.Close(); closeErr != nil {
			err = closeErr
		}
	}
	if a.cacheClient != nil {
		if closeErr := a.cacheClient.Close(); closeErr != nil {
			err = closeErr
		}
	}
	return err
}
