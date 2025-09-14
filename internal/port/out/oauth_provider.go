package out

import (
	"context"
	"encoding/json"
)

type OAuthUserInfo struct {
	ProviderID string
	Email      string
	Nickname   string
	RawData    json.RawMessage
}

type OAuthAPI interface {
	GetAccessToken(ctx context.Context, code string, clientID string, clientSecret []byte, redirectURI string) (string, error)
	GetUserInfo(ctx context.Context, accessToken string) (*OAuthUserInfo, error)
	GetAuthURL(ctx context.Context, clientID string, scopes []string, redirectURI string, state string) string
}
