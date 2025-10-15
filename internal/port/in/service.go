package in

import (
	"context"
	"time"

	"github.com/mandacode-com/mandacode-ssam/internal/domain/service"
)

type ServiceView struct {
	ServiceID service.PublicID
	Name      service.Name
	Desc      service.Description
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateServiceInput struct {
	Name        service.Name
	Description service.Description
}

type CreateServiceResult struct {
	ServiceView
}

type DeleteServiceInput struct {
	ServiceID service.PublicID
}

type UpdateServiceInput struct {
	ServiceID   service.PublicID
	NewName     *service.Name
	NewDesc     *service.Description
	NewIsActive *bool
}

type UpdateServiceResult struct {
	UpdatedAt time.Time
}

type FindServiceInput struct {
	Name     *service.Name
	PublicID *service.PublicID
}

type FindServiceByNameResult struct {
	ServiceView
}

type ListServicesInput struct {
	NameContains *string
	IsActive     *bool
	Limit        int
	Offset       int
}

type ListServicesResult struct {
	Services []ServiceView
}

type ServiceMgmtUsecase interface {
	CreateService(ctx context.Context, req *CreateServiceInput) (*CreateServiceResult, error)
	DeleteService(ctx context.Context, req *DeleteServiceInput) error
	UpdateService(ctx context.Context, req *UpdateServiceInput) (*UpdateServiceResult, error)
	FindService(ctx context.Context, req *FindServiceInput) (*FindServiceByNameResult, error)
	ListServices(ctx context.Context, req *ListServicesInput) (*ListServicesResult, error)
}
