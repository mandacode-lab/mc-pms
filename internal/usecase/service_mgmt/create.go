package servicemgmt

import (
	"context"

	"github.com/mandacode-com/mandacode-ssam/internal/domain/entity"
	vo "github.com/mandacode-com/mandacode-ssam/internal/domain/value_object"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
	"github.com/mandacode-com/mandacode-ssam/pkg/utils"
	"github.com/mandacode-com/merr"
)

func (u *Usecase) CreateService(ctx context.Context, req *in.CreateServiceRequest) (*in.CreateServiceResponse, error) {
	// Parse service name
	serviceName, err := vo.NewServiceName(req.Name)
	if err != nil {
		return nil, merr.New(merr.ErrBadRequest, ErrInvalidServiceNameMsg, err)
	}

	// Check if service with same name already exists
	existing, err := u.serviceQueryRepo.FindByName(ctx, serviceName)
	if err == nil && existing != nil {
		return nil, merr.New(merr.ErrConflict, ErrServiceAlreadyExistsMsg, nil)
	}

	// Create new service using domain logic
	desc := utils.StringNil(req.Description)
	serviceEntity := entity.DraftService(serviceName, desc)

	// Save to repository within transaction
	var savedService *entity.Service
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
	return &in.CreateServiceResponse{
		ServiceInfo: toServiceInfo(savedService),
	}, nil
}
