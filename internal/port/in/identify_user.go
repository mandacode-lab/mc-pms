package in

import (
	"context"
	
	clientappval "github.com/mandacode-com/mandacode-ssam/internal/domain/clientapp/value"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/shared"
)

type UserIdentityResult struct {
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
}

type GetAuthURLResult struct {
	AuthURL string
	State   string
}

type IdentifyUserUsecase interface {
	IdentifyByCode(ctx context.Context, cmd *IdentifyByCode) (*UserIdentityResult, error)
	IdentifyByToken(ctx context.Context, cmd *IdentifyByToken) (*UserIdentityResult, error)
	GetAuthURL(ctx context.Context, cmd *GetAuthURL) (*GetAuthURLResult, error)
}
