package servicemgmt

import (
	"context"

	"github.com/mandacode-com/mandacode-ssam/internal/domain/entity"
	vo "github.com/mandacode-com/mandacode-ssam/internal/domain/value_object"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/merr"
)

func (u *Usecase) FindService(ctx context.Context, req *in.FindServiceRequest) (*in.FindServiceByNameResponse, error) {
	if req == nil {
		return nil, merr.New(merr.ErrBadRequest, "invalid request", nil)
	}

	if req.Name == nil && req.PublicID == nil {
		return nil, merr.New(merr.ErrBadRequest, "either name or public ID must be provided", nil)
	}

	var service *entity.Service
	var err error

	if req.PublicID != nil {
		service, err = u.serviceQueryRepo.FindByPublicID(ctx, *req.PublicID)
	} else if req.Name != nil {
		name, nameErr := vo.NewServiceName(*req.Name)
		if nameErr != nil {
			return nil, merr.New(merr.ErrBadRequest, "invalid service name", nameErr)
		}
		service, err = u.serviceQueryRepo.FindByName(ctx, name)
	}

	if err != nil {
		return nil, merr.New(merr.ErrNotFound, ErrServiceNotFoundMsg, err)
	}

	if service == nil {
		return nil, merr.New(merr.ErrNotFound, ErrServiceNotFoundMsg, nil)
	}

	return &in.FindServiceByNameResponse{
		ServiceInfo: toServiceInfo(service),
	}, nil
}
