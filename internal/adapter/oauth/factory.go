package oauth

import (
	"fmt"
	"net/http"

	"github.com/mandacode-com/serengeti-integrated/internal/domain/shared"
	"github.com/mandacode-com/serengeti-integrated/internal/port/out"
)

type ProviderFactory struct {
	client *http.Client
}

func NewProviderFactory() *ProviderFactory {
	return &ProviderFactory{
		client: &http.Client{},
	}
}

func NewProviderFactoryWithClient(client *http.Client) *ProviderFactory {
	return &ProviderFactory{
		client: client,
	}
}

func (f *ProviderFactory) CreateProvider(provider shared.Provider) (out.OAuthAPI, error) {
	switch provider {
	case shared.ProviderGoogle:
		return NewGoogleOAuthWithClient(f.client), nil
	case shared.ProviderKakao:
		return NewKakaoOAuthWithClient(f.client), nil
	case shared.ProviderNaver:
		return NewNaverOAuthWithClient(f.client), nil
	default:
		return nil, fmt.Errorf("unsupported OAuth provider: %s", provider)
	}
}
