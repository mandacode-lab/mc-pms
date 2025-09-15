package in

import (
	"context"
	
	serviceval "github.com/mandacode-com/mandacode-service-hub/internal/domain/service/value"
)

type CreateServiceCommand struct {
	Name        serviceval.Name
	Description string
}

type CreateServiceResult struct {
	ServiceInfo
}

type DeleteServiceCommand struct {
	ServiceID serviceval.PublicID
}

type UpdateServiceCommand struct {
	ServiceID   serviceval.PublicID
	NewName     *serviceval.Name
	NewDesc     *string
	NewIsActive *bool
}

type UpdateServiceResult struct {
	ServiceInfo
}

type ServiceMgmtUsecase interface {
	CreateService(ctx context.Context, cmd *CreateServiceCommand) (*CreateServiceResult, error)
	DeleteService(ctx context.Context, cmd *DeleteServiceCommand) error
	UpdateService(ctx context.Context, cmd *UpdateServiceCommand) (*UpdateServiceResult, error)
}
