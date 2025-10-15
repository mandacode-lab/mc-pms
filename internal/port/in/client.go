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

type ClientView struct {
	ServiceID service.PublicID
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
	*ClientView
	Secret []byte
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
	UpdatedAt time.Time
}

type FindClientInput struct {
	ClientID *client.PublicID
	Name     *client.Name
}

type FindClientResult struct {
	*ClientView
	ServiceID service.PublicID
}

type ListClientsInput struct {
	ServiceID    *service.PublicID
	NameContains *string
	IsActive     *bool
	Limit        int
	Offset       int
}

type ListClientsResult struct {
	Clients []*ClientView
}

type ClientMgmtUsecase interface {
	CreateClient(ctx context.Context, req *CreateClientInput) (*CreateClientResult, error)
	DeleteClient(ctx context.Context, req *DeleteClientInput) error
	RefreshSecret(ctx context.Context, req *RefreshSecretInput) (*RefreshSecretResult, error)
	UpdateClient(ctx context.Context, req *UpdateClientInput) (*UpdateClientResult, error)
	FindClient(ctx context.Context, req *FindClientInput) (*FindClientResult, error)
	ListClients(ctx context.Context, req *ListClientsInput) (*ListClientsResult, error)
}
