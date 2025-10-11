package out

import (
	"context"

	"github.com/mandacode-com/mandacode-ssam/internal/domain/entity"
	vo "github.com/mandacode-com/mandacode-ssam/internal/domain/vo"
)

type ClientAppListFilter struct {
	ServiceID *vo.ServiceID
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
	Create(ctx context.Context, tx Tx, clientApp *entity.ClientApp) (*entity.ClientApp, error)
	Update(ctx context.Context, tx Tx, clientApp *entity.ClientApp) error
	FindByID(ctx context.Context, id vo.ClientAppID) (*entity.ClientApp, error)
	Delete(ctx context.Context, tx Tx, clientApp *entity.ClientApp) error
}

type ClientAppQueryRepository interface {
	FindByID(ctx context.Context, id vo.ClientAppID) (*entity.ClientApp, error)
	FindByPublicID(ctx context.Context, publicID vo.ClientAppPublicID) (*entity.ClientApp, error)
	FindByName(ctx context.Context, name string) (*entity.ClientApp, error)
	List(ctx context.Context, filter *ClientAppListFilter, options *ClientAppListOptions) ([]*entity.ClientApp, int, error)
}
