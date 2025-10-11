package in

import (
	"context"

	vo "github.com/mandacode-com/mandacode-ssam/internal/domain/vo"
)

type CreateClientAppRequest struct {
	ServiceID vo.ServicePublicID
	Name      string
	Desc      string
}

type CreateClientAppResponse struct {
	ClientAppInfo
	Secret []byte
}

type DeleteClientAppRequest struct {
	ClientAppID vo.ClientAppPublicID
}

type RefreshSecretRequest struct {
	ClientAppID vo.ClientAppPublicID
}

type RefreshSecretResponse struct {
	Secret []byte
}

type UpdateClientAppRequest struct {
	ClientAppID vo.ClientAppPublicID
	NewName     *string
	NewDesc     *string
	NewIsActive *bool
}

type UpdateClientAppResponse struct {
	ClientAppInfo
}

type FindClientAppRequest struct {
	ClientAppID *vo.ClientAppPublicID
	Name        *string
}

type FindClientAppResponse struct {
	ClientAppInfo
}

type ListClientAppsRequest struct {
	ServiceID    *vo.ServicePublicID
	NameContains *string
	IsActive     *bool
	Limit        int
	Offset       int
}

type ListClientAppsResponse struct {
	ServiceID  vo.ServicePublicID
	ClientApps []ClientAppInfo
}

type ClientAppMgmtUsecase interface {
	CreateClientApp(ctx context.Context, req *CreateClientAppRequest) (*CreateClientAppResponse, error)
	DeleteClientApp(ctx context.Context, req *DeleteClientAppRequest) error
	RefreshSecret(ctx context.Context, req *RefreshSecretRequest) (*RefreshSecretResponse, error)
	UpdateClientApp(ctx context.Context, req *UpdateClientAppRequest) (*UpdateClientAppResponse, error)
	FindClientApp(ctx context.Context, req *FindClientAppRequest) (*FindClientAppResponse, error)
	ListClientApps(ctx context.Context, req *ListClientAppsRequest) (*ListClientAppsResponse, error)
}
