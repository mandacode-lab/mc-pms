package in

import (
	"context"

	clientappval "github.com/mandacode-com/mandacode-ssam/internal/domain/clientapp/value"
	serviceval "github.com/mandacode-com/mandacode-ssam/internal/domain/service/value"
)

type CreateClientAppRequest struct {
	ServiceID serviceval.PublicID
	Name      string
	Desc      string
}

type CreateClientAppResponse struct {
	ClientAppInfo
	Secret []byte
}

type DeleteClientAppRequest struct {
	ClientAppID clientappval.PublicID
}

type RefreshSecretRequest struct {
	ClientAppID clientappval.PublicID
}

type RefreshSecretResponse struct {
	Secret []byte
}

type UpdateClientAppRequest struct {
	ClientAppID clientappval.PublicID
	NewName     *string
	NewDesc     *string
	NewIsActive *bool
}

type UpdateClientAppResponse struct {
	ClientAppInfo
}

type ListClientAppsRequest struct {
	ServiceID serviceval.PublicID
}

type ListClientAppsResponse struct {
	ServiceID  serviceval.PublicID
	ClientApps []ClientAppInfo
}

type ClientAppMgmtUsecase interface {
	CreateClientApp(ctx context.Context, req *CreateClientAppRequest) (*CreateClientAppResponse, error)
	DeleteClientApp(ctx context.Context, req *DeleteClientAppRequest) error
	RefreshSecret(ctx context.Context, req *RefreshSecretRequest) (*RefreshSecretResponse, error)
	UpdateClientApp(ctx context.Context, req *UpdateClientAppRequest) (*UpdateClientAppResponse, error)
	ListClientApps(ctx context.Context, req *ListClientAppsRequest) (*ListClientAppsResponse, error)
}
