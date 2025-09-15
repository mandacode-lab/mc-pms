package in

import (
	"context"
)

type CreateClientAppRequest struct {
	ServiceID string
	Name      string
	Desc      string
}

type CreateClientAppResponse struct {
	ClientAppInfo
	Secret []byte
}

type DeleteClientAppRequest struct {
	ClientAppID string
}

type RefreshSecretRequest struct {
	ClientAppID string
}

type RefreshSecretResponse struct {
	Secret []byte
}

type UpdateClientAppRequest struct {
	ClientAppID string
	NewName     *string
	NewDesc     *string
	NewIsActive *bool
}

type UpdateClientAppResponse struct {
	ClientAppInfo
}

type ListClientAppsRequest struct {
	ServiceID string
}

type ListClientAppsResponse struct {
	ServiceID  string
	ClientApps []ClientAppInfo
}

type ClientAppMgmtUsecase interface {
	CreateClientApp(ctx context.Context, req *CreateClientAppRequest) (*CreateClientAppResponse, error)
	DeleteClientApp(ctx context.Context, req *DeleteClientAppRequest) error
	RefreshSecret(ctx context.Context, req *RefreshSecretRequest) (*RefreshSecretResponse, error)
	UpdateClientApp(ctx context.Context, req *UpdateClientAppRequest) (*UpdateClientAppResponse, error)
	ListClientApps(ctx context.Context, req *ListClientAppsRequest) (*ListClientAppsResponse, error)
}
