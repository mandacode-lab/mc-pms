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
	"github.com/mandacode-com/mandacode-ssam/internal/adapter/handler/clientapp_mgmt"
	"github.com/mandacode-com/mandacode-ssam/internal/adapter/handler/service_mgmt"
	"github.com/mandacode-com/mandacode-ssam/internal/adapter/handler/user_mgmt"
	"github.com/mandacode-com/mandacode-ssam/internal/adapter/handler/weboauth_mgmt"
	"github.com/mandacode-com/mandacode-ssam/internal/adapter/hasher"
	"github.com/mandacode-com/mandacode-ssam/internal/adapter/kek"
	"github.com/mandacode-com/mandacode-ssam/internal/adapter/random"
	"github.com/mandacode-com/mandacode-ssam/internal/adapter/repository"
	redisinfra "github.com/mandacode-com/mandacode-ssam/internal/infra/redis"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
	clientapp_mgmt_usecase "github.com/mandacode-com/mandacode-ssam/internal/usecase/clientapp_mgmt"
	service_mgmt_usecase "github.com/mandacode-com/mandacode-ssam/internal/usecase/service_mgmt"
	user_mgmt_usecase "github.com/mandacode-com/mandacode-ssam/internal/usecase/user_mgmt"
	weboauth_mgmt_usecase "github.com/mandacode-com/mandacode-ssam/internal/usecase/weboauth_mgmt"
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
	serviceRepo           out.ServiceRepository
	serviceQueryRepo      out.ServiceQueryRepository
	clientAppRepo         out.ClientAppRepository
	clientAppQueryRepo    out.ClientAppQueryRepository
	webOAuthRepo          out.WebOAuthRepository
	webOAuthQueryRepo     out.WebOAuthQueryRepository
	userIdentityRepo      out.UserIdentityRepository
	userIdentityQueryRepo out.UserIdentityQueryRepository
	userInfoRepo          out.UserInfoRepository
	userInfoQueryRepo     out.UserInfoQueryRepository

	// Services
	hasher      out.Hasher
	kekProvider out.KekProvider
	secretGen   out.ByteRandGen
	encoder     out.Encoder
	txManager   out.TransactionManager
}

func NewAdapter(ctx context.Context, cfg *configs.CoreConfig) (*Adapter, error) {
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
	webOAuthRepo := repository.NewEntWebOAuthRepository(client)
	webOAuthQueryRepo := repository.NewCachedWebOAuthQueryRepository(
		repository.NewEntWebOAuthQueryRepository(client),
		cacheStore,
		cacheConfig,
	)
	userIdentityRepo := repository.NewEntUserIdentityRepository(client)
	userIdentityQueryRepo := repository.NewEntUserIdentityQueryRepository(client)
	userInfoRepo := repository.NewEntUserInfoRepository(client)
	userInfoQueryRepo := repository.NewEntUserInfoQueryRepository(client)

	// Services
	hasherSvc := hasher.NewBcryptHasher()

	// KEK Provider from config
	kekProvider, err := kek.NewAESKekProviderFromHex(cfg.KEK)
	if err != nil {
		return nil, fmt.Errorf("creating KEK provider: %w", err)
	}
	secretGen := random.NewCryptoByteRandGen()
	encoderSvc := encoder.NewBase64Encoder()
	txManager := repository.NewEntTransactionManager(client)

	return &Adapter{
		db:                    db,
		client:                client,
		cacheClient:           cacheClient,
		cacheStore:            cacheStore,
		serviceRepo:           serviceRepo,
		serviceQueryRepo:      serviceQueryRepo,
		clientAppRepo:         clientAppRepo,
		clientAppQueryRepo:    clientAppQueryRepo,
		webOAuthRepo:          webOAuthRepo,
		webOAuthQueryRepo:     webOAuthQueryRepo,
		userIdentityRepo:      userIdentityRepo,
		userIdentityQueryRepo: userIdentityQueryRepo,
		userInfoRepo:          userInfoRepo,
		userInfoQueryRepo:     userInfoQueryRepo,
		hasher:                hasherSvc,
		kekProvider:           kekProvider,
		secretGen:             secretGen,
		encoder:               encoderSvc,
		txManager:             txManager,
	}, nil
}

func (a *Adapter) ProvideServiceMgmtUsecase() in.ServiceMgmtUsecase {
	return service_mgmt_usecase.NewUsecase(
		a.serviceRepo,
		a.serviceQueryRepo,
		a.txManager,
	)
}

func (a *Adapter) ProvideServiceMgmtHandler() *service_mgmt.Handler {
	return service_mgmt.NewHandler(a.ProvideServiceMgmtUsecase())
}

func (a *Adapter) ProvideClientAppMgmtUsecase() in.ClientAppMgmtUsecase {
	return clientapp_mgmt_usecase.NewUsecase(
		a.clientAppRepo,
		a.clientAppQueryRepo,
		a.serviceQueryRepo,
		a.txManager,
		a.secretGen,
		a.hasher,
		a.encoder,
	)
}

func (a *Adapter) ProvideClientAppMgmtHandler() *clientapp_mgmt.Handler {
	return clientapp_mgmt.NewHandler(a.ProvideClientAppMgmtUsecase())
}

func (a *Adapter) ProvideWebOAuthMgmtUsecase() in.WebOAuthMgmtUsecase {
	return weboauth_mgmt_usecase.NewUsecase(
		a.webOAuthRepo,
		a.webOAuthQueryRepo,
		a.clientAppQueryRepo,
		a.txManager,
		a.kekProvider,
	)
}

func (a *Adapter) ProvideWebOAuthMgmtHandler() *weboauth_mgmt.Handler {
	return weboauth_mgmt.NewHandler(a.ProvideWebOAuthMgmtUsecase())
}

func (a *Adapter) ProvideUserMgmtUsecase() in.UserMgmtUsecase {
	return user_mgmt_usecase.NewUsecase(
		a.userInfoQueryRepo,
		a.userIdentityQueryRepo,
		a.userIdentityRepo,
		a.userInfoRepo,
		a.serviceQueryRepo,
		a.txManager,
	)
}

func (a *Adapter) ProvideUserMgmtHandler() *user_mgmt.Handler {
	return user_mgmt.NewHandler(a.ProvideUserMgmtUsecase())
}

func (a *Adapter) Close() error {
	if a.client != nil {
		a.client.Close()
	}
	if a.db != nil {
		a.db.Close()
	}
	if a.cacheClient != nil {
		a.cacheClient.Close()
	}
	return nil
}
