package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	clientappval "github.com/mandacode-com/mandacode-service-hub/internal/domain/clientapp/value"
	"github.com/mandacode-com/mandacode-service-hub/internal/domain/shared"
	"github.com/mandacode-com/mandacode-service-hub/internal/domain/weboauth"
	weboauthval "github.com/mandacode-com/mandacode-service-hub/internal/domain/weboauth/value"
	"github.com/mandacode-com/mandacode-service-hub/internal/port/out"
)

const (
	WebOAuthByProviderKeyPrefix  = "weboauth:provider:"
	WebOAuthByClientAppKeyPrefix = "weboauth:clientapp:"
)

type CachedWebOAuthQueryRepository struct {
	repo        out.WebOAuthQueryRepository
	cache       out.CacheStore
	cacheConfig *out.CacheConfig
}

func NewCachedWebOAuthQueryRepository(repo out.WebOAuthQueryRepository, cache out.CacheStore, cacheConfig *out.CacheConfig) out.WebOAuthQueryRepository {
	return &CachedWebOAuthQueryRepository{
		repo:        repo,
		cache:       cache,
		cacheConfig: cacheConfig,
	}
}

type cachedWebOAuth struct {
	ID               int64     `json:"id"`
	ClientAppID      int64     `json:"client_app_id"`
	Provider         string    `json:"provider"`
	OAuthClientID    string    `json:"oauth_client_id"`
	OAuthSecretCT    []byte    `json:"oauth_secret_ct"`
	OAuthSecretNonce []byte    `json:"oauth_secret_nonce"`
	DEKWrapped       []byte    `json:"dek_wrapped"`
	DEKNonce         []byte    `json:"dek_nonce"`
	DEKRotatedAt     time.Time `json:"dek_rotated_at"`
	RedirectURI      string    `json:"redirect_uri"`
	Scopes           []string  `json:"scopes"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (c *CachedWebOAuthQueryRepository) FindByProvider(ctx context.Context, clientAppID clientappval.ID, provider shared.Provider) (*weboauth.WebOAuth, error) {
	cacheKey := fmt.Sprintf("%s%d:%s", WebOAuthByProviderKeyPrefix, clientAppID.Value(), string(provider))

	// Try cache first
	cached, err := c.cache.Get(ctx, cacheKey)
	if err != nil {
		// Continue to database if cache fails
	} else if cached != nil {
		var cachedData cachedWebOAuth
		if err := json.Unmarshal(cached, &cachedData); err == nil {
			return c.fromCached(&cachedData), nil
		}
	}

	// Fallback to database
	webOAuth, err := c.repo.FindByProvider(ctx, clientAppID, provider)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if cachedData := c.toCached(webOAuth); cachedData != nil {
		if data, err := json.Marshal(cachedData); err == nil {
			c.cache.Set(ctx, cacheKey, data) // Ignore cache errors
		}
	}

	return webOAuth, nil
}

func (c *CachedWebOAuthQueryRepository) FindByID(ctx context.Context, id weboauthval.ID) (*weboauth.WebOAuth, error) {
	return c.repo.FindByID(ctx, id)
}

func (c *CachedWebOAuthQueryRepository) FindByClientAppID(ctx context.Context, clientAppID clientappval.ID) ([]*weboauth.WebOAuth, error) {
	// For list operations, we don't cache as they're less predictable
	return c.repo.FindByClientAppID(ctx, clientAppID)
}

func (c *CachedWebOAuthQueryRepository) List(ctx context.Context, filter *out.WebOAuthListFilter, options *out.WebOAuthListOptions) ([]*weboauth.WebOAuth, int, error) {
	return c.repo.List(ctx, filter, options)
}

func (c *CachedWebOAuthQueryRepository) toCached(wo *weboauth.WebOAuth) *cachedWebOAuth {
	if wo == nil {
		return nil
	}

	return &cachedWebOAuth{
		ID:               wo.ID().Value(),
		ClientAppID:      wo.ClientAppID().Value(),
		Provider:         string(wo.Provider()),
		OAuthClientID:    wo.OAuthClientID(),
		OAuthSecretCT:    wo.OAuthSecretCT(),
		OAuthSecretNonce: wo.OAuthSecretNonce(),
		DEKWrapped:       wo.DEKWrapped(),
		DEKNonce:         wo.DEKNonce(),
		DEKRotatedAt:     wo.DEKRotatedAt(),
		RedirectURI:      wo.RedirectURI(),
		Scopes:           wo.Scopes(),
		CreatedAt:        wo.CreatedAt(),
		UpdatedAt:        wo.UpdatedAt(),
	}
}

func (c *CachedWebOAuthQueryRepository) fromCached(cached *cachedWebOAuth) *weboauth.WebOAuth {
	return weboauth.NewWebOAuth(
		weboauthval.NewID(cached.ID),
		clientappval.NewID(cached.ClientAppID),
		shared.Provider(cached.Provider),
		cached.OAuthClientID,
		cached.OAuthSecretCT,
		cached.OAuthSecretNonce,
		cached.DEKWrapped,
		cached.DEKNonce,
		cached.DEKRotatedAt,
		cached.RedirectURI,
		cached.Scopes,
		cached.CreatedAt,
		cached.UpdatedAt,
	)
}

