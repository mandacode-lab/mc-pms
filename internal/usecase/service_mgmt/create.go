package service_mgmt

import (
	"context"

	"github.com/mandacode-com/merr"
	"github.com/mandacode-com/mandacode-service-hub/internal/domain/service"
	"github.com/mandacode-com/mandacode-service-hub/internal/port/in"
	"github.com/mandacode-com/mandacode-service-hub/internal/port/out"
	"github.com/mandacode-com/mandacode-service-hub/pkg/utils"
)

func (u *Usecase) CreateService(ctx context.Context, cmd *in.CreateServiceCommand) (*in.CreateServiceResult, error) {
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

	// Convert to result
	return &in.CreateServiceResult{
		ServiceInfo: toServiceInfo(savedService),
	}, nil
}
