package cache

import (
	"context"
	"fmt"

	clientappval "github.com/mandacode-com/serengeti-integrated/internal/domain/clientapp/value"
	"github.com/mandacode-com/serengeti-integrated/internal/domain/shared"
	"github.com/mandacode-com/serengeti-integrated/internal/port/out"
)

// CacheInvalidator handles cache invalidation for WebOAuth configurations
type CacheInvalidator struct {
	cache out.CacheStore
}

func NewCacheInvalidator(cache out.CacheStore) *CacheInvalidator {
	return &CacheInvalidator{
		cache: cache,
	}
}

// InvalidateWebOAuth removes cached web oauth data
func (c *CacheInvalidator) InvalidateWebOAuth(ctx context.Context, clientAppID clientappval.ID, provider shared.Provider) error {
	key := fmt.Sprintf("weboauth:provider:%d:%s", clientAppID.Value(), string(provider))
	return c.cache.Del(ctx, key)
}