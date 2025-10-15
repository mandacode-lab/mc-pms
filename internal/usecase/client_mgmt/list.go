package clientmgmt

import (
	"context"

	"github.com/mandacode-com/mandacode-ssam/internal/domain/client"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/service"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/merr"
)

func (u *Usecase) ListClients(ctx context.Context, req *in.ListClientsInput) (*in.ListClientsResult, error) {
	if req == nil {
		return nil, merr.New(merr.ErrBadRequest, "invalid request", nil)
	}

	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 20 // Default limit
	}

	if req.Offset < 0 {
		req.Offset = 0
	}

	// var serviceID service.PublicID
	var serviceInternalID *service.ID

	// Convert PublicID to internal ID for repository query
	if req.ServiceID != nil {
		serviceEntity, err := u.serviceQueryRepo.FindByPublicID(ctx, *req.ServiceID)
		if err != nil {
			return nil, merr.New(merr.ErrNotFound, ErrServiceNotFoundMsg, err)
		}
		internalID := serviceEntity.ID()
		serviceInternalID = &internalID
	}

	filter := &client.ClientListFilter{
		ServiceID:    serviceInternalID,
		NameContains: req.NameContains,
		IsActive:     req.IsActive,
	}

	options := &client.ClientListOptions{
		Limit:  req.Limit,
		Offset: req.Offset,
		Order:  client.ClientListOrderNameAsc,
	}

	result, total, err := u.clientQueryRepo.List(ctx, filter, options)
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}

	_ = total // We can use this for pagination info later

	return &in.ListClientsResult{
		Clients: toClientInfos(result),
	}, nil
}
