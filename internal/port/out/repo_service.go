package out

import (
	"context"

	"github.com/mandacode-com/serengeti/internal/domain/service"
	serviceval "github.com/mandacode-com/serengeti/internal/domain/service/value"
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
	Create(ctx context.Context, tx Tx, service *service.Service) (*service.Service, error)
	Update(ctx context.Context, tx Tx, service *service.Service) error
	FindByID(ctx context.Context, id serviceval.ID) (*service.Service, error)
	Delete(ctx context.Context, tx Tx, service *service.Service) error
}

type ServiceQueryRepository interface {
	FindByID(ctx context.Context, id serviceval.ID) (*service.Service, error)
	FindByPublicID(ctx context.Context, publicID serviceval.PublicID) (*service.Service, error)
	FindByName(ctx context.Context, name serviceval.Name) (*service.Service, error)
	List(ctx context.Context, filter *ServiceListFilter, options *ServiceListOptions) ([]*service.Service, int, error)
}