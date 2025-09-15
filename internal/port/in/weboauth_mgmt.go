package in

import (
	"context"
)

type RegisterWebOAuthRequest struct {
	ClientAppID   string
	RedirectURI   string
	Provider      string
	Scopes        []string
	OAuthClientID string
	OAuthSecret   []byte
}

type ReadWebOAuthRequest struct {
	ClientAppID string
	Provider    *string
}

type ReadWebOAuthResponse struct {
	WebOAuths []WebOAuthInfo
}

type UpdateWebOAuthRequest struct {
	ClientAppID   string
	Provider      string
	RedirectURI   *string
	Scopes        *[]string
	OAuthClientID *string
	OAuthSecret   *[]byte
}

type UpdateWebOAuthResponse struct {
	WebOAuthInfo
}

type DeleteWebOAuthRequest struct {
	ClientAppID string
	Provider    string
}

type ReadWebOAuthSecretRequest struct {
	ClientAppID string
	Provider    string
}

type ReadWebOAuthSecretResponse struct {
	OAuthSecret []byte
}

type WebOAuthMgmtUsecase interface {
	RegisterWebOAuth(ctx context.Context, req *RegisterWebOAuthRequest) error
	ReadWebOAuth(ctx context.Context, req *ReadWebOAuthRequest) (*ReadWebOAuthResponse, error)
	ReadWebOAuthSecret(ctx context.Context, req *ReadWebOAuthSecretRequest) (*ReadWebOAuthSecretResponse, error)
	UpdateWebOAuth(ctx context.Context, req *UpdateWebOAuthRequest) (*UpdateWebOAuthResponse, error)
	DeleteWebOAuth(ctx context.Context, req *DeleteWebOAuthRequest) error
}
