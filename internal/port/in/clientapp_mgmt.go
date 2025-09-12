package in

import (
	clientappval "github.com/mandacode-com/serengeti-integrated/internal/domain/clientapp/value"
	serviceval "github.com/mandacode-com/serengeti-integrated/internal/domain/service/value"
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
	CreateClientApp(cmd CreateClientAppCommand) (CreateClientAppView, error)
	DeleteClientApp(cmd DeleteClientAppCommand) error
	RefreshSecret(cmd RefreshSecretCommand) (RefreshSecretView, error)
	UpdateClientApp(cmd UpdateClientAppCommand) (UpdateClientAppView, error)
	ListClientApps(cmd ListClientAppsCommand) (ListClientAppsView, error)
}
