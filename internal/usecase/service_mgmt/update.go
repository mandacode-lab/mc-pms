package service_mgmt

import (
	"context"

	serviceval "github.com/mandacode-com/mandacode-ssam/internal/domain/service/value"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
	"github.com/mandacode-com/merr"
)

func (u *Usecase) UpdateService(ctx context.Context, req *in.UpdateServiceRequest) (*in.UpdateServiceResponse, error) {
	// Parse service ID
	serviceID, err := serviceval.ParsePublicID(req.ServiceID)
	if err != nil {
		return nil, merr.New(merr.ErrBadRequest, ErrInvalidServiceIDMsg, err)
	}

	// Find the service by public ID
	service, err := u.serviceQueryRepo.FindByPublicID(ctx, serviceID)
	if err != nil {
		return nil, merr.New(merr.ErrNotFound, ErrServiceNotFoundMsg, err)
	}

	// Apply updates using domain logic
	if req.NewName != nil {
		serviceName, err := serviceval.NewName(*req.NewName)
		if err != nil {
			return nil, merr.New(merr.ErrBadRequest, ErrInvalidServiceNameMsg, err)
		}
		service.UpdateName(serviceName)
	}

	if req.NewDesc != nil {
		service.UpdateDescription(req.NewDesc)
	}

	if req.NewIsActive != nil {
		if *req.NewIsActive {
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
	return &in.UpdateServiceResponse{
		ServiceInfo: toServiceInfo(service),
	}, nil
}
