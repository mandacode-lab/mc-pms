package in

import (
	clientappval "github.com/mandacode-com/serengeti-integrated/internal/domain/clientapp/value"
	"github.com/mandacode-com/serengeti-integrated/internal/domain/shared"
)

type UserIdentityView struct {
	UserInfo
	RawData map[string]interface{}
}

type IdentifyByCode struct {
	Provider        shared.Provider
	OAuthCode       string
	ClientAppID     clientappval.PublicID
	ClientAppSecret []byte
}

type IdentifyByToken struct {
	Provider        shared.Provider
	Token           string
	ClientApp       clientappval.PublicID
	ClientAppSecret []byte
}

type IdentifyUserUsecase interface {
	IdentifyByCode(cmd IdentifyByCode) (UserIdentityView, error)
	IdentifyByToken(cmd IdentifyByToken) (UserIdentityView, error)
}
