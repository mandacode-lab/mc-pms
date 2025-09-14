package out

import (
	"context"

	"github.com/mandacode-com/serengeti/internal/domain/clientapp"
	clientappval "github.com/mandacode-com/serengeti/internal/domain/clientapp/value"
	serviceval "github.com/mandacode-com/serengeti/internal/domain/service/value"
)

type ClientAppListFilter struct {
	ServiceID *serviceval.ID
	Name      *string
	IsActive  *bool
}

type ClientAppListOrder string

const (
	ClientAppListOrderNameAsc  ClientAppListOrder = "name_asc"
	ClientAppListOrderNameDesc ClientAppListOrder = "name_desc"
)

type ClientAppListOptions struct {
	Limit  int
	Offset int
	Order  ClientAppListOrder
}


type ClientAppRepository interface {
	Create(ctx context.Context, tx Tx, clientApp *clientapp.ClientApp) (*clientapp.ClientApp, error)
	Update(ctx context.Context, tx Tx, clientApp *clientapp.ClientApp) error
	FindByID(ctx context.Context, id clientappval.ID) (*clientapp.ClientApp, error)
	Delete(ctx context.Context, tx Tx, clientApp *clientapp.ClientApp) error
}

type ClientAppQueryRepository interface {
	FindByID(ctx context.Context, id clientappval.ID) (*clientapp.ClientApp, error)
	FindByPublicID(ctx context.Context, publicID clientappval.PublicID) (*clientapp.ClientApp, error)
	FindByName(ctx context.Context, name string) (*clientapp.ClientApp, error)
	FindByServiceID(ctx context.Context, serviceID serviceval.ID) ([]*clientapp.ClientApp, error)
	List(ctx context.Context, filter *ClientAppListFilter, options *ClientAppListOptions) ([]*clientapp.ClientApp, int, error)
}