package servicemgmt

import (
	"context"
	"time"

	"github.com/mandacode-com/mandacode-ssam/internal/domain/service"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/tx"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/merr"
)

func (u *Usecase) CreateService(ctx context.Context, req *in.CreateServiceInput) (*in.CreateServiceResult, error) {
	// Check if service with same name already exists
	existing, err := u.serviceQueryRepo.FindByName(ctx, req.Name)
	if err == nil && existing != nil {
		return nil, merr.New(merr.ErrConflict, ErrServiceAlreadyExistsMsg, nil)
	}

	id, err := service.NewID(1) // ID will be set by repository
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}
	publicIDStr, err := u.serviceIDGen.Generate()
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}
	publicID, err := service.NewPublicID(publicIDStr)
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}
	now := time.Now().UTC()
	serviceModel := service.NewService(
		id,
		publicID,
		req.Name,
		req.Description,
		true, // isActive
		now,
		now,
	)

	// Save to repository within transaction
	var savedService *service.Service
	err = u.txManager.WithTx(ctx, func(tx tx.Tx) error {
		saved, err := u.serviceRepo.Create(ctx, tx, serviceModel)
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
		ServiceView: toServiceInfo(savedService),
	}, nil
}
