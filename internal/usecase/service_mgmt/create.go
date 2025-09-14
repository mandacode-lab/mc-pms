package service_mgmt

import (
	"context"

	"github.com/mandacode-com/merr"
	"github.com/mandacode-com/serengeti/internal/domain/service"
	"github.com/mandacode-com/serengeti/internal/port/in"
	"github.com/mandacode-com/serengeti/internal/port/out"
	"github.com/mandacode-com/serengeti/pkg/utils"
)

func (u *Usecase) CreateService(ctx context.Context, cmd *in.CreateServiceCommand) (*in.CreateServiceView, error) {
	// Check if service with same name already exists
	existing, err := u.serviceQueryRepo.FindByName(ctx, cmd.Name)
	if err == nil && existing != nil {
		return nil, merr.New(merr.ErrConflict, ErrServiceAlreadyExistsMsg, nil)
	}

	// Create new service using domain logic
	desc := utils.StringNil(cmd.Description)
	serviceEntity := service.DraftService(cmd.Name, desc)

	// Save to repository within transaction
	var savedService *service.Service
	err = u.txManager.WithTx(ctx, func(tx out.Tx) error {
		saved, err := u.serviceRepo.Create(ctx, tx, serviceEntity)
		if err != nil {
			return merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
		}
		savedService = saved
		return nil
	})
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}

	// Convert to view
	return &in.CreateServiceView{
		ServiceInfo: toServiceInfo(savedService),
	}, nil
}
