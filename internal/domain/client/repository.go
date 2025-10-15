package client

import (
	"context"

	"github.com/mandacode-com/mandacode-ssam/internal/domain/service"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/tx"
)

type ClientListFilter struct {
	ServiceID *service.ID
	Name      *Name
	IsActive  *bool
}

type ClientListOrder string

const (
	ClientListOrderNameAsc  ClientListOrder = "name_asc"
	ClientListOrderNameDesc ClientListOrder = "name_desc"
)

type ClientListOptions struct {
	Limit  int
	Offset int
	Order  ClientListOrder
}

type ClientRepository interface {
	Create(ctx context.Context, tx tx.Tx, client *Client) (*Client, error)
	Update(ctx context.Context, tx tx.Tx, client *Client) error
	Delete(ctx context.Context, tx tx.Tx, client *Client) error
	FindByID(ctx context.Context, id ID) (*Client, error)
}

type ClientQueryRepository interface {
	FindByID(ctx context.Context, id ID) (*Client, error)
	FindByPublicID(ctx context.Context, publicID PublicID) (*Client, error)
	FindByName(ctx context.Context, name Name) (*Client, error)
	List(ctx context.Context, filter *ClientListFilter, options *ClientListOptions) ([]*Client, int, error)
}
