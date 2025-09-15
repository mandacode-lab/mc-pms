package state

import (
	"context"
	"fmt"

	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
)

const (
	StateKeyPrefix = "oauth_state:"
)

type CacheStateService struct {
	cache      out.CacheStore
	strRandGen out.StrRandGen
}

func NewCacheStateService(cache out.CacheStore, strRandGen out.StrRandGen) out.StateService {
	return &CacheStateService{
		cache:      cache,
		strRandGen: strRandGen,
	}
}

func (s *CacheStateService) GenerateState(ctx context.Context) (string, error) {
	// Generate random state using string generator
	state, err := s.strRandGen.Generate(ctx)
	if err != nil {
		return "", fmt.Errorf("generating random state: %w", err)
	}

	// Store in cache with configured expiration
	key := StateKeyPrefix + state
	err = s.cache.Set(ctx, key, []byte("1"))
	if err != nil {
		return "", fmt.Errorf("storing state in cache: %w", err)
	}

	return state, nil
}

func (s *CacheStateService) ValidateState(ctx context.Context, state string) bool {
	key := StateKeyPrefix + state

	// Check if state exists in cache
	exists, err := s.cache.Exists(ctx, key)
	if err != nil {
		return false
	}

	if !exists {
		return false
	}

	// Delete the state after validation (one-time use)
	s.cache.Del(ctx, key)

	return true
}
