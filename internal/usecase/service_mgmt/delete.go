package service_mgmt

import (
	"context"

	serviceval "github.com/mandacode-com/mandacode-ssam/internal/domain/service/value"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
	"github.com/mandacode-com/merr"
)

func (u *Usecase) DeleteService(ctx context.Context, req *in.DeleteServiceRequest) error {
	// Parse service ID
	serviceID, err := serviceval.ParsePublicID(req.ServiceID)
	if err != nil {
		return merr.New(merr.ErrBadRequest, ErrInvalidServiceIDMsg, err)
	}

	// Find the service by public ID
	service, err := u.serviceQueryRepo.FindByPublicID(ctx, serviceID)
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
