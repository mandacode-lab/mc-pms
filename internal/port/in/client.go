package in

import (
	"context"
	"time"

	"github.com/mandacode-com/mandacode-ssam/internal/domain/client"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/service"
)

type VerifyClientInput struct {
	ClientID     client.PublicID
	ClientSecret []byte // raw secret
}

type VerifyClientResult struct {
	ServiceID service.PublicID
}

type ClientAccessUsecase interface {
	VerifyClient(ctx context.Context, req *VerifyClientInput) (*VerifyClientResult, error)
}

type MgmtClientInfo struct {
	ClientID  client.PublicID
	Name      client.Name
	Desc      client.Description
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateClientInput struct {
	ServiceID service.PublicID
	Name      client.Name
	Desc      client.Description
}

type CreateClientResult struct {
	MgmtClientInfo
	Secret    []byte
	ServiceID service.PublicID
}

type DeleteClientInput struct {
	ClientID client.PublicID
}

type RefreshSecretInput struct {
	ClientID client.PublicID
}

type RefreshSecretResult struct {
	RawSecret []byte
}

type UpdateClientInput struct {
	ClientID    client.PublicID
	NewName     *client.Name
	NewDesc     *client.Description
	NewIsActive *bool
}

type UpdateClientResult struct {
	MgmtClientInfo
}

type FindClientInput struct {
	ClientID *client.PublicID
	Name     *client.Name
}

type FindClientResult struct {
	ServiceID service.PublicID
	MgmtClientInfo
}

type ListClientsInput struct {
	ServiceID    *service.PublicID
	NameContains *string
	IsActive     *bool
	Limit        int
	Offset       int
}

type ListClientsResult struct {
	Clients []MgmtClientInfo
}

type ClientMgmtUsecase interface {
	CreateClient(ctx context.Context, req *CreateClientInput) (*CreateClientResult, error)
	DeleteClient(ctx context.Context, req *DeleteClientInput) error
	RefreshSecret(ctx context.Context, req *RefreshSecretInput) (*RefreshSecretResult, error)
	UpdateClient(ctx context.Context, req *UpdateClientInput) (*UpdateClientResult, error)
	FindClient(ctx context.Context, req *FindClientInput) (*FindClientResult, error)
	ListClients(ctx context.Context, req *ListClientsInput) (*ListClientsResult, error)
}
