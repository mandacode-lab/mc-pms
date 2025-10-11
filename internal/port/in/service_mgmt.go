package in

import (
	"context"

	vo "github.com/mandacode-com/mandacode-ssam/internal/domain/vo"
)

type CreateServiceRequest struct {
	Name        string
	Description string
}

type CreateServiceResponse struct {
	ServiceInfo
}

type DeleteServiceRequest struct {
	ServiceID vo.ServicePublicID
}

type UpdateServiceRequest struct {
	ServiceID   vo.ServicePublicID
	NewName     *string
	NewDesc     *string
	NewIsActive *bool
}

type UpdateServiceResponse struct {
	ServiceInfo
}

type FindServiceRequest struct {
	Name     *string
	PublicID *vo.ServicePublicID
}

type FindServiceByNameResponse struct {
	ServiceInfo
}

type ListServicesRequest struct {
	NameContains *string
	IsActive     *bool
	Limit        int
	Offset       int
}

type ListServicesResponse struct {
	Services []ServiceInfo
}

type ServiceMgmtUsecase interface {
	CreateService(ctx context.Context, req *CreateServiceRequest) (*CreateServiceResponse, error)
	DeleteService(ctx context.Context, req *DeleteServiceRequest) error
	UpdateService(ctx context.Context, req *UpdateServiceRequest) (*UpdateServiceResponse, error)
	FindService(ctx context.Context, req *FindServiceRequest) (*FindServiceByNameResponse, error)
	ListServices(ctx context.Context, req *ListServicesRequest) (*ListServicesResponse, error)
}
