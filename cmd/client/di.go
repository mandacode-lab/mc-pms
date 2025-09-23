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
	"github.com/mandacode-com/mandacode-ssam/internal/adapter/handler/client_access"
	"github.com/mandacode-com/mandacode-ssam/internal/adapter/hasher"
	"github.com/mandacode-com/mandacode-ssam/internal/adapter/repository"
	redisinfra "github.com/mandacode-com/mandacode-ssam/internal/infra/redis"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
	access_usecase "github.com/mandacode-com/mandacode-ssam/internal/usecase/client_access"
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
	serviceQueryRepo   out.ServiceQueryRepository
	clientAppQueryRepo out.ClientAppQueryRepository

	// Services
	hasher out.Hasher
}

func NewAdapter(ctx context.Context, cfg *configs.ClientConfig) (*Adapter, error) {
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

	// Repositories (only read repositories needed)
	serviceQueryRepo := repository.NewEntServiceQueryRepository(client)
	clientAppQueryRepo := repository.NewEntClientAppQueryRepository(client)

	// Services
	hasherSvc := hasher.NewBcryptHasher()

	return &Adapter{
		db:                 db,
		client:             client,
		cacheClient:        cacheClient,
		cacheStore:         cacheStore,
		serviceQueryRepo:   serviceQueryRepo,
		clientAppQueryRepo: clientAppQueryRepo,
		hasher:             hasherSvc,
	}, nil
}

func (a *Adapter) ProvideClientAccessUsecase() in.ClientAccessUsecase {
	return access_usecase.NewUsecase(
		a.clientAppQueryRepo,
		a.serviceQueryRepo,
		a.hasher,
	)
}

func (a *Adapter) ProvideClientAccessHandler() *clientaccess.Handler {
	return clientaccess.NewHandler(a.ProvideClientAccessUsecase())
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
