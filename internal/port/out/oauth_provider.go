package out

import (
	"encoding/json"
)

type OAuthUserInfo struct {
	ProviderID string
	Email      string
	RawData    json.RawMessage
}

type OAuthAPI interface {
	GetAccessToken(code string, clientID string, clientSecret []byte) (string, error)
	GetUserInfo(accessToken string) (*OAuthUserInfo, error)
	GetAuthURL(clientID string, scopes []string, redirectURI string, state string)
}
