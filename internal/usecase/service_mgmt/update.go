package service_mgmt

import (
	"context"

	"github.com/mandacode-com/merr"
	"github.com/mandacode-com/mandacode-service-hub/internal/port/in"
	"github.com/mandacode-com/mandacode-service-hub/internal/port/out"
)

func (u *Usecase) UpdateService(ctx context.Context, cmd *in.UpdateServiceCommand) (*in.UpdateServiceResult, error) {
	// Find the service by public ID
	service, err := u.serviceQueryRepo.FindByPublicID(ctx, cmd.ServiceID)
	if err != nil {
		return nil, merr.New(merr.ErrNotFound, ErrServiceNotFoundMsg, err)
	}

	// Apply updates using domain logic
	if cmd.NewName != nil {
		service.UpdateName(*cmd.NewName)
	}

	if cmd.NewDesc != nil {
		service.UpdateDescription(cmd.NewDesc)
	}

	if cmd.NewIsActive != nil {
		if *cmd.NewIsActive {
			service.Activate()
		} else {
			service.Deactivate()
		}
	}

	// Update in repository within transaction
	err = u.txManager.WithTx(ctx, func(tx out.Tx) error {
		return u.serviceRepo.Update(ctx, tx, service)
	})
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}

	// Convert to result
	return &in.UpdateServiceResult{
		ServiceInfo: toServiceInfo(service),
	}, nil
}
