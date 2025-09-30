package out

import (
	"context"

	"github.com/mandacode-com/mandacode-ssam/internal/domain/entity"
	vo "github.com/mandacode-com/mandacode-ssam/internal/domain/value_object"
)

type ServiceListFilter struct {
	Name     *string
	IsActive *bool
}

type ServiceListOrder string

const (
	ServiceListOrderNameAsc  ServiceListOrder = "name_asc"
	ServiceListOrderNameDesc ServiceListOrder = "name_desc"
)

type ServiceListOptions struct {
	Limit  int
	Offset int
	Order  ServiceListOrder
}


type ServiceRepository interface {
	Create(ctx context.Context, tx Tx, service *entity.Service) (*entity.Service, error)
	Update(ctx context.Context, tx Tx, service *entity.Service) error
	FindByID(ctx context.Context, id vo.ServiceID) (*entity.Service, error)
	Delete(ctx context.Context, tx Tx, service *entity.Service) error
}

type ServiceQueryRepository interface {
	FindByID(ctx context.Context, id vo.ServiceID) (*entity.Service, error)
	FindByPublicID(ctx context.Context, publicID vo.ServicePublicID) (*entity.Service, error)
	FindByName(ctx context.Context, name vo.ServiceName) (*entity.Service, error)
	List(ctx context.Context, filter *ServiceListFilter, options *ServiceListOptions) ([]*entity.Service, int, error)
}