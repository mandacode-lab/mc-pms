package in

import (
	clientappval "github.com/mandacode-com/serengeti-integrated/internal/domain/clientapp/value"
	"github.com/mandacode-com/serengeti-integrated/internal/domain/shared"
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
	RegisterWebOAuth(cmd RegisterWebOAuthCommand) error
	ReadWebOAuth(query ReadWebOAuthQuery) (ReadWebOAuthView, error)
	ReadWebOAuthSecret(query ReadWebOAuthSecretQuery) (ReadWebOAuthSecretView, error)
	UpdateWebOAuth(cmd UpdateWebOAuthCommand) error
	DeleteWebOAuth(cmd DeleteWebOAuthCommand) error
}
