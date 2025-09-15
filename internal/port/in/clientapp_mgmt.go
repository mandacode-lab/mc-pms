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

type CreateClientAppView struct {
	ClientAppInfo
	Secret []byte
}

type DeleteClientAppCommand struct {
	ClientAppID clientappval.PublicID
}

type RefreshSecretCommand struct {
	ClientAppID clientappval.PublicID
}

type RefreshSecretView struct {
	Secret []byte
}

type UpdateClientAppCommand struct {
	ClientAppID clientappval.PublicID
	NewName     *string
	NewDesc     *string
	NewIsActive *bool
}

type UpdateClientAppView struct {
	ClientAppInfo
}

type ListClientAppsCommand struct {
	ServiceID serviceval.PublicID
}

type ListClientAppsView struct {
	ServiceID  serviceval.PublicID
	ClientApps []ClientAppInfo
}


type ClientAppMgmtUsecase interface {
	CreateClientApp(ctx context.Context, cmd *CreateClientAppCommand) (*CreateClientAppView, error)
	DeleteClientApp(ctx context.Context, cmd *DeleteClientAppCommand) error
	RefreshSecret(ctx context.Context, cmd *RefreshSecretCommand) (*RefreshSecretView, error)
	UpdateClientApp(ctx context.Context, cmd *UpdateClientAppCommand) (*UpdateClientAppView, error)
	ListClientApps(ctx context.Context, cmd *ListClientAppsCommand) (*ListClientAppsView, error)
}
