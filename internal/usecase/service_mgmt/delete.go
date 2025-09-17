package service_mgmt

import (
	"context"

	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
	"github.com/mandacode-com/merr"
)

func (u *Usecase) DeleteService(ctx context.Context, req *in.DeleteServiceRequest) error {
	// Find the service by public ID
	service, err := u.serviceQueryRepo.FindByPublicID(ctx, req.ServiceID)
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
