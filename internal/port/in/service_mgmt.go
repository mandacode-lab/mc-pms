package in

import serviceval "github.com/mandacode-com/serengeti-integrated/internal/domain/service/value"

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
	CreateService(cmd CreateServiceCommand) (CreateServiceView, error)
	DeleteService(cmd DeleteServiceCommand) error
	UpdateService(cmd UpdateServiceCommand) (UpdateServiceView, error)
}
