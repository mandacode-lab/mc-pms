package in

import (
	"context"
)

type UserIdentityResponse struct {
	UserInfo
	RawData map[string]any
}

type IdentifyByCodeRequest struct {
	Provider        string
	OAuthCode       string
	State           string
	ClientAppID     string
	ClientAppSecret []byte
}

type IdentifyByTokenRequest struct {
	Provider        string
	Token           string
	ClientAppID     string
	ClientAppSecret []byte
}

type GetAuthURLRequest struct {
	Provider        string
	ClientAppID     string
	ClientAppSecret []byte
}

type GetAuthURLResponse struct {
	AuthURL string
	State   string
}

type IdentifyUserUsecase interface {
	IdentifyByCode(ctx context.Context, req *IdentifyByCodeRequest) (*UserIdentityResponse, error)
	IdentifyByToken(ctx context.Context, req *IdentifyByTokenRequest) (*UserIdentityResponse, error)
	GetAuthURL(ctx context.Context, req *GetAuthURLRequest) (*GetAuthURLResponse, error)
}
