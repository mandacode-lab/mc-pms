package in

import (
	"context"
)

type CreateServiceRequest struct {
	Name        string
	Description string
}

type CreateServiceResponse struct {
	ServiceInfo
}

type DeleteServiceRequest struct {
	ServiceID string
}

type UpdateServiceRequest struct {
	ServiceID   string
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
