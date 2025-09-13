package in

import (
	"time"

	clientappval "github.com/mandacode-com/serengeti-integrated/internal/domain/clientapp/value"
	serviceval "github.com/mandacode-com/serengeti-integrated/internal/domain/service/value"
	"github.com/mandacode-com/serengeti-integrated/internal/domain/shared"
	useridentityval "github.com/mandacode-com/serengeti-integrated/internal/domain/useridentity/value"
)

type ServiceInfo struct {
	ServiceID serviceval.PublicID
	Name      string
	Desc      string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ClientAppInfo struct {
	ServiceID   serviceval.PublicID
	ClientAppID clientappval.PublicID
	Name        string
	Desc        string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UserInfo struct {
	ServiceID serviceval.PublicID
	UserID    useridentityval.PublicID
	Nickname  string
	Email     string
	Provider  shared.Provider
	CreatedAt time.Time
	UpdatedAt time.Time
}

type WebOAuthInfo struct {
	ClientAppID   clientappval.PublicID
	Provider      shared.Provider
	OAuthClientID string
	RedirectURI   string
	Scopes        []string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
