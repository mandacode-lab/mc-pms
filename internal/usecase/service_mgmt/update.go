package servicemgmt

import (
	"context"

	"github.com/mandacode-com/mandacode-ssam/internal/domain/tx"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/merr"
)

func (u *Usecase) UpdateService(ctx context.Context, req *in.UpdateServiceInput) (*in.UpdateServiceResult, error) {
	// Find the serviceModel by public ID
	serviceModel, err := u.serviceQueryRepo.FindByPublicID(ctx, req.ServiceID)
	if err != nil {
		return nil, merr.New(merr.ErrNotFound, ErrServiceNotFoundMsg, err)
	}

	// Apply updates using domain logic
	if req.NewName != nil {
		err := serviceModel.Rename(*req.NewName)
		if err != nil {
			return nil, merr.New(merr.ErrBadRequest, "invalid service name", err)
		}
	}

	if req.NewDesc != nil {
		err := serviceModel.UpdateDescription(*req.NewDesc)
		if err != nil {
			return nil, merr.New(merr.ErrBadRequest, "invalid service description", err)
		}
	}

	if req.NewIsActive != nil {
		if *req.NewIsActive {
			err := serviceModel.Activate()
			if err != nil {
				return nil, merr.New(merr.ErrBadRequest, "invalid service state", err)
			}
		} else {
			err := serviceModel.Deactivate()
			if err != nil {
				return nil, merr.New(merr.ErrBadRequest, "invalid service state", err)
			}
		}
	}

	// Update in repository within transaction
	err = u.txManager.WithTx(ctx, func(tx tx.Tx) error {
		return u.serviceRepo.Update(ctx, tx, serviceModel)
	})
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}

	// Convert to result
	return &in.UpdateServiceResult{
		UpdatedAt: serviceModel.UpdatedAt(),
	}, nil
}
