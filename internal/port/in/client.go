package in

import (
	"context"
	"time"

	"github.com/mandacode-com/mandacode-ssam/internal/domain/service"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/svcclient"
)

type VerifyClientInput struct {
	ClientID     svcclient.PublicID
	ClientSecret []byte // raw secret
}

type VerifyClientResult struct {
	ServiceID service.PublicID
}

type ClientAccessUsecase interface {
	VerifyClient(ctx context.Context, req *VerifyClientInput) (*VerifyClientResult, error)
}

type MgmtClientInfo struct {
	ClientID  svcclient.PublicID
	Name      svcclient.Name
	Desc      svcclient.Description
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateClientInput struct {
	ServiceID service.PublicID
	Name      service.Name
	Desc      service.Description
}

type CreateClientResult struct {
	MgmtClientInfo
	Secret []byte
}

type DeleteClientInput struct {
	ClientID svcclient.PublicID
}

type RefreshSecretInput struct {
	ClientID svcclient.PublicID
}

type RefreshSecretResult struct {
	RawSecret []byte
}

type UpdateClientInput struct {
	ClientID    svcclient.PublicID
	NewName     *string
	NewDesc     *string
	NewIsActive *bool
}

type UpdateClientResult struct {
	MgmtClientInfo
}

type FindClientInput struct {
	ClientID *svcclient.PublicID
	Name     *svcclient.Name
}

type FindClientResult struct {
	MgmtClientInfo
}

type ListClientInput struct {
	ServiceID    service.PublicID
	NameContains *string
	IsActive     *bool
	Limit        int
	Offset       int
}

type ListClientResult struct {
	ServiceID service.PublicID
	Clients   []MgmtClientInfo
}

type ClientMgmtUsecase interface {
	CreateClient(ctx context.Context, req *CreateClientInput) (*CreateClientResult, error)
	DeleteClient(ctx context.Context, req *DeleteClientInput) error
	RefreshSecret(ctx context.Context, req *RefreshSecretInput) (*RefreshSecretResult, error)
	UpdateClient(ctx context.Context, req *UpdateClientInput) (*UpdateClientResult, error)
	FindClient(ctx context.Context, req *FindClientInput) (*FindClientResult, error)
	ListClients(ctx context.Context, req *ListClientInput) (*ListClientResult, error)
}
