package in

import (
	"context"
	
	clientappval "github.com/mandacode-com/mandacode-service-hub/internal/domain/clientapp/value"
	"github.com/mandacode-com/mandacode-service-hub/internal/domain/shared"
)

type RegisterWebOAuthCommand struct {
	ClientAppID   clientappval.PublicID
	RedirectURI   string
	Provider      shared.Provider
	Scopes        []string
	OAuthClientID string
	OAuthSecret   []byte
}

type ReadWebOAuthQuery struct {
	ClientAppID clientappval.PublicID
	Provider    *shared.Provider
}

type ReadWebOAuthView struct {
	WebOAuths []WebOAuthInfo
}

type UpdateWebOAuthCommand struct {
	ClientAppID   clientappval.PublicID
	Provider      shared.Provider
	RedirectURI   *string
	Scopes        *[]string
	OAuthClientID *string
	OAuthSecret   *[]byte
}

type UpdateWebOAuthView struct {
	WebOAuthInfo
}

type DeleteWebOAuthCommand struct {
	ClientAppID clientappval.PublicID
	Provider    shared.Provider
}

type ReadWebOAuthSecretQuery struct {
	ClientAppID clientappval.PublicID
	Provider    shared.Provider
}

type ReadWebOAuthSecretView struct {
	OAuthSecret []byte
}

type WebOAuthMgmtUsecase interface {
	RegisterWebOAuth(ctx context.Context, cmd *RegisterWebOAuthCommand) error
	ReadWebOAuth(ctx context.Context, query *ReadWebOAuthQuery) (*ReadWebOAuthView, error)
	ReadWebOAuthSecret(ctx context.Context, query *ReadWebOAuthSecretQuery) (*ReadWebOAuthSecretView, error)
	UpdateWebOAuth(ctx context.Context, cmd *UpdateWebOAuthCommand) error
	DeleteWebOAuth(ctx context.Context, cmd *DeleteWebOAuthCommand) error
}
