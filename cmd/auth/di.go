package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"github.com/mandacode-com/serengeti-integrated/configs"
	"github.com/mandacode-com/serengeti-integrated/ent"
	"github.com/mandacode-com/serengeti-integrated/internal/adapter/cache"
	"github.com/mandacode-com/serengeti-integrated/internal/adapter/encoder"
	"github.com/mandacode-com/serengeti-integrated/internal/adapter/handler/identify_user"
	"github.com/mandacode-com/serengeti-integrated/internal/adapter/hasher"
	"github.com/mandacode-com/serengeti-integrated/internal/adapter/kek"
	"github.com/mandacode-com/serengeti-integrated/internal/adapter/oauth"
	"github.com/mandacode-com/serengeti-integrated/internal/adapter/random"
	"github.com/mandacode-com/serengeti-integrated/internal/adapter/repository"
	"github.com/mandacode-com/serengeti-integrated/internal/adapter/state"
	"github.com/mandacode-com/serengeti-integrated/internal/domain/shared"
	"github.com/mandacode-com/serengeti-integrated/internal/port/in"
	"github.com/mandacode-com/serengeti-integrated/internal/port/out"
	identify_user_usecase "github.com/mandacode-com/serengeti-integrated/internal/usecase/identify_user"
)

type Adapter struct {
	// Database
	db     *sql.DB
	client *ent.Client

	// Cache
	rdb        *redis.Client
	cacheStore out.CacheStore

	// Repositories
	clientAppQueryRepo   out.ClientAppQueryRepository
	serviceQueryRepo     out.ServiceQueryRepository
	userIdentityRepo     out.UserIdentityRepository
	userIdentityQueryRepo out.UserIdentityQueryRepository
	userInfoRepo         out.UserInfoRepository
	userInfoQueryRepo    out.UserInfoQueryRepository
	webOAuthQueryRepo    out.WebOAuthQueryRepository

	// Services
	hasher      out.Hasher
	kekProvider out.KekProvider
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
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Cache setup
	cacheConfig := &out.CacheConfig{
		DefaultExpiration: 30 * time.Minute,
	}
	cacheStore := cache.NewRedisCacheStore(rdb, cacheConfig)

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
		rdb:                   rdb,
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
	if a.client != nil {
		a.client.Close()
	}
	if a.db != nil {
		a.db.Close()
	}
	if a.rdb != nil {
		a.rdb.Close()
	}
	return nil
}