package service_mgmt

import (
	"context"

	"github.com/mandacode-com/merr"
	"github.com/mandacode-com/mandacode-service-hub/internal/port/in"
	"github.com/mandacode-com/mandacode-service-hub/internal/port/out"
)

func (u *Usecase) DeleteService(ctx context.Context, cmd *in.DeleteServiceCommand) error {
	// Find the service by public ID
	service, err := u.serviceQueryRepo.FindByPublicID(ctx, cmd.ServiceID)
	if err != nil {
		return merr.New(merr.ErrNotFound, ErrServiceNotFoundMsg, err)
	}

	// Delete within transaction
	err = u.txManager.WithTx(ctx, func(tx out.Tx) error {
		return u.serviceRepo.Delete(ctx, tx, service)
	})
	if err != nil {
		return merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}

	return nil
}
