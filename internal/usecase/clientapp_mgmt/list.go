package clientapp_mgmt

import (
	"context"

	serviceval "github.com/mandacode-com/mandacode-ssam/internal/domain/service/value"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
	"github.com/mandacode-com/merr"
)

func (u *Usecase) ListClientApps(ctx context.Context, req *in.ListClientAppsRequest) (*in.ListClientAppsResponse, error) {
	if req == nil {
		return nil, merr.New(merr.ErrBadRequest, "invalid request", nil)
	}

	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 20 // Default limit
	}

	if req.Offset < 0 {
		req.Offset = 0
	}

	var serviceID serviceval.PublicID
	var serviceInternalID *serviceval.ID

	// Convert PublicID to internal ID for repository query
	if req.ServiceID != nil {
		serviceEntity, err := u.serviceQueryRepo.FindByPublicID(ctx, *req.ServiceID)
		if err != nil {
			return nil, merr.New(merr.ErrNotFound, ErrServiceNotFoundMsg, err)
		}
		internalID := serviceEntity.ID()
		serviceInternalID = &internalID
		serviceID = *req.ServiceID
	}

	filter := &out.ClientAppListFilter{
		ServiceID: serviceInternalID,
		Name:      req.NameContains,
		IsActive:  req.IsActive,
	}

	options := &out.ClientAppListOptions{
		Limit:  req.Limit,
		Offset: req.Offset,
		Order:  out.ClientAppListOrderNameAsc,
	}

	clientApps, total, err := u.clientAppQueryRepo.List(ctx, filter, options)
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}

	_ = total // We can use this for pagination info later

	return &in.ListClientAppsResponse{
		ServiceID:  serviceID,
		ClientApps: toClientAppInfos(clientApps, serviceID),
	}, nil
}