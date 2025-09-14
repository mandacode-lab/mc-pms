package weboauth_mgmt

import (
	clientappval "github.com/mandacode-com/serengeti/internal/domain/clientapp/value"
	"github.com/mandacode-com/serengeti/internal/domain/weboauth"
	"github.com/mandacode-com/serengeti/internal/port/in"
)

func toWebOAuthInfo(wo *weboauth.WebOAuth, clientAppPublicID clientappval.PublicID) in.WebOAuthInfo {
	return in.WebOAuthInfo{
		ClientAppID:   clientAppPublicID,
		Provider:      wo.Provider(),
		OAuthClientID: wo.OAuthClientID(),
		RedirectURI:   wo.RedirectURI(),
		Scopes:        wo.Scopes(),
		CreatedAt:     wo.CreatedAt(),
		UpdatedAt:     wo.UpdatedAt(),
	}
}
