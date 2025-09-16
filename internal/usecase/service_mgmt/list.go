package service_mgmt

import (
	"context"

	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
	"github.com/mandacode-com/merr"
)

func (u *Usecase) ListServices(ctx context.Context, req *in.ListServicesRequest) (*in.ListServicesResponse, error) {
	if req == nil {
		return nil, merr.New(merr.ErrBadRequest, "invalid request", nil)
	}

	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 20 // Default limit
	}

	if req.Offset < 0 {
		req.Offset = 0
	}

	filter := &out.ServiceListFilter{
		Name:     req.NameContains,
		IsActive: req.IsActive,
	}

	options := &out.ServiceListOptions{
		Limit:  req.Limit,
		Offset: req.Offset,
		Order:  out.ServiceListOrderNameAsc,
	}

	services, total, err := u.serviceQueryRepo.List(ctx, filter, options)
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}

	_ = total // We can use this for pagination info later

	serviceInfos := make([]in.ServiceInfo, len(services))
	for i, svc := range services {
		serviceInfos[i] = toServiceInfo(svc)
	}

	return &in.ListServicesResponse{
		Services: serviceInfos,
	}, nil
}