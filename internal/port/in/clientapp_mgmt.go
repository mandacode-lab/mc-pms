package in

import (
	"context"

	clientappval "github.com/mandacode-com/mandacode-service-hub/internal/domain/clientapp/value"
	serviceval "github.com/mandacode-com/mandacode-service-hub/internal/domain/service/value"
)

type CreateClientAppCommand struct {
	ServiceID serviceval.PublicID
	Name      string
	Desc      string
}

type CreateClientAppResult struct {
	ClientAppInfo
	Secret []byte
}

type DeleteClientAppCommand struct {
	ClientAppID clientappval.PublicID
}

type RefreshSecretCommand struct {
	ClientAppID clientappval.PublicID
}

type RefreshSecretResult struct {
	Secret []byte
}

type UpdateClientAppCommand struct {
	ClientAppID clientappval.PublicID
	NewName     *string
	NewDesc     *string
	NewIsActive *bool
}

type UpdateClientAppResult struct {
	ClientAppInfo
}

type ListClientAppsCommand struct {
	ServiceID serviceval.PublicID
}

type ListClientAppsResult struct {
	ServiceID  serviceval.PublicID
	ClientApps []ClientAppInfo
}


type ClientAppMgmtUsecase interface {
	CreateClientApp(ctx context.Context, cmd *CreateClientAppCommand) (*CreateClientAppResult, error)
	DeleteClientApp(ctx context.Context, cmd *DeleteClientAppCommand) error
	RefreshSecret(ctx context.Context, cmd *RefreshSecretCommand) (*RefreshSecretResult, error)
	UpdateClientApp(ctx context.Context, cmd *UpdateClientAppCommand) (*UpdateClientAppResult, error)
	ListClientApps(ctx context.Context, cmd *ListClientAppsCommand) (*ListClientAppsResult, error)
}
