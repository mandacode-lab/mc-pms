package out

import (
	"context"

	"github.com/mandacode-com/serengeti-integrated/internal/domain/weboauth"
	weboauthval "github.com/mandacode-com/serengeti-integrated/internal/domain/weboauth/value"
	clientappval "github.com/mandacode-com/serengeti-integrated/internal/domain/clientapp/value"
	"github.com/mandacode-com/serengeti-integrated/internal/domain/shared"
)

type WebOAuthListFilter struct {
	ClientAppID *clientappval.ID
	Provider    *shared.Provider
}

type WebOAuthListOrder string

const (
	WebOAuthListOrderProviderAsc  WebOAuthListOrder = "provider_asc"
	WebOAuthListOrderProviderDesc WebOAuthListOrder = "provider_desc"
)

type WebOAuthListOptions struct {
	Limit  int
	Offset int
	Order  WebOAuthListOrder
}


type WebOAuthRepository interface {
	Create(ctx context.Context, tx Tx, oauth *weboauth.WebOAuth) (*weboauth.WebOAuth, error)
	Update(ctx context.Context, tx Tx, oauth *weboauth.WebOAuth) error
	FindByID(ctx context.Context, id weboauthval.ID) (*weboauth.WebOAuth, error)
	Delete(ctx context.Context, tx Tx, oauth *weboauth.WebOAuth) error
}

type WebOAuthQueryRepository interface {
	FindByID(ctx context.Context, id weboauthval.ID) (*weboauth.WebOAuth, error)
	FindByProvider(ctx context.Context, clientAppID clientappval.ID, provider shared.Provider) (*weboauth.WebOAuth, error)
	FindByClientAppID(ctx context.Context, clientAppID clientappval.ID) ([]*weboauth.WebOAuth, error)
	List(ctx context.Context, filter *WebOAuthListFilter, options *WebOAuthListOptions) ([]*weboauth.WebOAuth, int, error)
}