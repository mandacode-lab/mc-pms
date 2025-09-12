package in

import (
	"context"
	
	serviceval "github.com/mandacode-com/serengeti-integrated/internal/domain/service/value"
)

type CreateServiceCommand struct {
	Name        serviceval.Name
	Description string
}

type CreateServiceView struct {
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

type UpdateServiceView struct {
	ServiceInfo
}

type ServiceMgmtUsecase interface {
	CreateService(ctx context.Context, cmd CreateServiceCommand) (CreateServiceView, error)
	DeleteService(ctx context.Context, cmd DeleteServiceCommand) error
	UpdateService(ctx context.Context, cmd UpdateServiceCommand) (UpdateServiceView, error)
}
