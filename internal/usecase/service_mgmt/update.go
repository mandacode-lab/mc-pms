package servicemgmt

import (
	"context"

	vo "github.com/mandacode-com/mandacode-ssam/internal/domain/vo"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
	"github.com/mandacode-com/merr"
)

func (u *Usecase) UpdateService(ctx context.Context, req *in.UpdateServiceRequest) (*in.UpdateServiceResponse, error) {
	// Find the service by public ID
	service, err := u.serviceQueryRepo.FindByPublicID(ctx, req.ServiceID)
	if err != nil {
		return nil, merr.New(merr.ErrNotFound, ErrServiceNotFoundMsg, err)
	}

	// Apply updates using domain logic
	if req.NewName != nil {
		serviceName, err := vo.NewServiceName(*req.NewName)
		if err != nil {
			return nil, merr.New(merr.ErrBadRequest, ErrInvalidServiceNameMsg, err)
		}
		service.UpdateName(serviceName)
	}

	if req.NewDesc != nil {
		desc, err := vo.NewServiceDescription(*req.NewDesc)
		if err != nil {
			return nil, merr.New(merr.ErrBadRequest, "invalid service description", err)
		}
		service.UpdateDescription(desc)
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
