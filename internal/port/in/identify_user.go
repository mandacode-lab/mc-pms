package in

import (
	"context"
	
	clientappval "github.com/mandacode-com/serengeti-integrated/internal/domain/clientapp/value"
	"github.com/mandacode-com/serengeti-integrated/internal/domain/shared"
)

type UserIdentityView struct {
	UserInfo
	RawData map[string]any
}

type IdentifyByCode struct {
	Provider        shared.Provider
	OAuthCode       string
	State           string
	ClientAppID     clientappval.PublicID
	ClientAppSecret []byte
}

type IdentifyByToken struct {
	Provider        shared.Provider
	Token           string
	ClientApp       clientappval.PublicID
	ClientAppSecret []byte
}

type GetAuthURL struct {
	Provider        shared.Provider
	ClientAppID     clientappval.PublicID
	ClientAppSecret []byte
	RedirectURI     string
	Scopes          []string
}

type GetAuthURLView struct {
	AuthURL string
	State   string
}

type IdentifyUserUsecase interface {
	IdentifyByCode(ctx context.Context, cmd *IdentifyByCode) (*UserIdentityView, error)
	IdentifyByToken(ctx context.Context, cmd *IdentifyByToken) (*UserIdentityView, error)
	GetAuthURL(ctx context.Context, cmd *GetAuthURL) (*GetAuthURLView, error)
}
