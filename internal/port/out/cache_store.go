package out

import (
	"context"
	"time"
)

// CacheStore provides generic caching operations with configurable expiration
type CacheStore interface {
	// Get retrieves a value by key
	Get(ctx context.Context, key string) ([]byte, error)

	// Set stores a value with key using injected expiration
	Set(ctx context.Context, key string, value []byte) error

	// Del deletes a key
	Del(ctx context.Context, key string) error

	// Exists checks if a key exists
	Exists(ctx context.Context, key string) (bool, error)
}

// CacheConfig holds configuration for cache operations including expiration
type CacheConfig struct {
	DefaultExpiration time.Duration
}