package in

import (
	"context"

	serviceval "github.com/mandacode-com/mandacode-ssam/internal/domain/service/value"
)

type CreateServiceRequest struct {
	Name        string
	Description string
}

type CreateServiceResponse struct {
	ServiceInfo
}

type DeleteServiceRequest struct {
	ServiceID serviceval.PublicID
}

type UpdateServiceRequest struct {
	ServiceID   serviceval.PublicID
	NewName     *string
	NewDesc     *string
	NewIsActive *bool
}

type UpdateServiceResponse struct {
	ServiceInfo
}

type ServiceMgmtUsecase interface {
	CreateService(ctx context.Context, req *CreateServiceRequest) (*CreateServiceResponse, error)
	DeleteService(ctx context.Context, req *DeleteServiceRequest) error
	UpdateService(ctx context.Context, req *UpdateServiceRequest) (*UpdateServiceResponse, error)
}
