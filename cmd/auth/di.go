package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"
	"github.com/mandacode-com/serengeti/configs"
	"github.com/mandacode-com/serengeti/ent"
	"github.com/mandacode-com/serengeti/internal/adapter/cache"
	"github.com/mandacode-com/serengeti/internal/adapter/encoder"
	"github.com/mandacode-com/serengeti/internal/adapter/handler/identify_user"
	"github.com/mandacode-com/serengeti/internal/adapter/hasher"
	"github.com/mandacode-com/serengeti/internal/adapter/kek"
	"github.com/mandacode-com/serengeti/internal/adapter/oauth"
	"github.com/mandacode-com/serengeti/internal/adapter/random"
	"github.com/mandacode-com/serengeti/internal/adapter/repository"
	"github.com/mandacode-com/serengeti/internal/adapter/state"
	"github.com/mandacode-com/serengeti/internal/domain/shared"
	redisinfra "github.com/mandacode-com/serengeti/internal/infra/redis"
	"github.com/mandacode-com/serengeti/internal/port/in"
	"github.com/mandacode-com/serengeti/internal/port/out"
	identify_user_usecase "github.com/mandacode-com/serengeti/internal/usecase/identify_user"
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
	clientAppQueryRepo    out.ClientAppQueryRepository
	serviceQueryRepo      out.ServiceQueryRepository
	userIdentityRepo      out.UserIdentityRepository
	userIdentityQueryRepo out.UserIdentityQueryRepository
	userInfoRepo          out.UserInfoRepository
	userInfoQueryRepo     out.UserInfoQueryRepository
	webOAuthQueryRepo     out.WebOAuthQueryRepository

	// Services
	hasher       out.Hasher
	kekProvider  out.KekProvider
	stateService out.StateService
	strRandGen   out.StrRandGen
	encoder      out.Encoder
	txManager    out.TransactionManager

	// OAuth Providers
	oauthProviders map[shared.Provider]out.OAuthAPI
}

func NewAdapter(ctx context.Context, cfg *configs.AuthConfig) (*Adapter, error) {
	// Database setup
	db, err := sql.Open("postgres", cfg.Postgres.DSN())
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

	// Cache setup
	cacheConfig := &out.CacheConfig{
		DefaultExpiration: 30 * time.Minute,
	}
	cacheStore := cache.NewRedisCacheStore(cacheClient, cacheConfig)

	// Repositories
	clientAppQueryRepo := repository.NewEntClientAppQueryRepository(client)
	serviceQueryRepo := repository.NewEntServiceQueryRepository(client)
	userIdentityRepo := repository.NewEntUserIdentityRepository(client)
	userIdentityQueryRepo := repository.NewEntUserIdentityQueryRepository(client)
	userInfoRepo := repository.NewEntUserInfoRepository(client)
	userInfoQueryRepo := repository.NewEntUserInfoQueryRepository(client)
	webOAuthQueryRepo := repository.NewCachedWebOAuthQueryRepository(
		repository.NewEntWebOAuthQueryRepository(client),
		cacheStore,
		cacheConfig,
	)

	// Services
	hasherSvc := hasher.NewBcryptHasher()

	// KEK Provider from config
	kekProvider, err := kek.NewAESKekProviderFromHex(cfg.KEK)
	if err != nil {
		return nil, fmt.Errorf("creating KEK provider: %w", err)
	}
	strRandGen := random.NewCryptoStrRandGen()
	stateService := state.NewCacheStateService(cacheStore, strRandGen)
	encoderSvc := encoder.NewBase64Encoder()
	txManager := repository.NewEntTransactionManager(client)

	// OAuth Providers
	oauthProviders := map[shared.Provider]out.OAuthAPI{
		shared.ProviderGoogle: oauth.NewGoogleOAuth(),
		shared.ProviderKakao:  oauth.NewKakaoOAuth(),
		shared.ProviderNaver:  oauth.NewNaverOAuth(),
	}

	return &Adapter{
		db:                    db,
		client:                client,
		cacheClient:           cacheClient,
		cacheStore:            cacheStore,
		clientAppQueryRepo:    clientAppQueryRepo,
		serviceQueryRepo:      serviceQueryRepo,
		userIdentityRepo:      userIdentityRepo,
		userIdentityQueryRepo: userIdentityQueryRepo,
		userInfoRepo:          userInfoRepo,
		userInfoQueryRepo:     userInfoQueryRepo,
		webOAuthQueryRepo:     webOAuthQueryRepo,
		hasher:                hasherSvc,
		kekProvider:           kekProvider,
		stateService:          stateService,
		strRandGen:            strRandGen,
		encoder:               encoderSvc,
		txManager:             txManager,
		oauthProviders:        oauthProviders,
	}, nil
}

func (a *Adapter) ProvideIdentifyUserUsecase() in.IdentifyUserUsecase {
	return identify_user_usecase.NewUsecase(
		a.oauthProviders,
		a.clientAppQueryRepo,
		a.webOAuthQueryRepo,
		a.userIdentityRepo,
		a.userIdentityQueryRepo,
		a.userInfoRepo,
		a.userInfoQueryRepo,
		a.serviceQueryRepo,
		a.txManager,
		a.hasher,
		a.kekProvider,
		a.stateService,
	)
}

func (a *Adapter) ProvideIdentifyUserHandler() *identify_user.Handler {
	return identify_user.NewHandler(a.ProvideIdentifyUserUsecase())
}

func (a *Adapter) Close() error {
	if err := a.client.Close(); err != nil {
		return err
	}
	if err := a.db.Close(); err != nil {
		return err
	}
	if err := a.cacheClient.Close(); err != nil {
		return err
	}
	return nil
}
