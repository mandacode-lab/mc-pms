package service

import (
	"context"

	"github.com/mandacode-com/mandacode-ssam/internal/domain/tx"
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
	Create(ctx context.Context, tx tx.Tx, service *Service) (*Service, error)
	Update(ctx context.Context, tx tx.Tx, service *Service) error
	FindByID(ctx context.Context, id ID) (*Service, error)
	Delete(ctx context.Context, tx tx.Tx, service *Service) error
}

type ServiceQueryRepository interface {
	FindByID(ctx context.Context, id ID) (*Service, error)
	FindByPublicID(ctx context.Context, publicID PublicID) (*Service, error)
	FindByName(ctx context.Context, name Name) (*Service, error)
	List(ctx context.Context, filter *ServiceListFilter, options *ServiceListOptions) ([]*Service, int, error)
}
