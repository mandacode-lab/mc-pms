package cache

import (
	"context"
	"errors"

	"github.com/mandacode-com/serengeti-integrated/internal/port/out"
	"github.com/redis/go-redis/v9"
)

type RedisCacheStore struct {
	client redis.UniversalClient
	config *out.CacheConfig
}

func NewRedisCacheStore(client redis.UniversalClient, config *out.CacheConfig) out.CacheStore {
	return &RedisCacheStore{
		client: client,
		config: config,
	}
}

func (r *RedisCacheStore) Get(ctx context.Context, key string) ([]byte, error) {
	result, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil // Key doesn't exist
		}
		return nil, err
	}
	return []byte(result), nil
}

func (r *RedisCacheStore) Set(ctx context.Context, key string, value []byte) error {
	return r.client.Set(ctx, key, value, r.config.DefaultExpiration).Err()
}

func (r *RedisCacheStore) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisCacheStore) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
