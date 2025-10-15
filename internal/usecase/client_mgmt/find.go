package clientmgmt

import (
	"context"

	"github.com/mandacode-com/mandacode-ssam/internal/domain/client"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/merr"
)

func (u *Usecase) FindClient(ctx context.Context, req *in.FindClientInput) (*in.FindClientResult, error) {
	if req == nil {
		return nil, merr.New(merr.ErrBadRequest, "invalid request", nil)
	}

	if req.ClientID == nil && req.Name == nil {
		return nil, merr.New(merr.ErrBadRequest, "either client app ID or name must be provided", nil)
	}

	var svcClient *client.Client
	var err error

	if req.ClientID != nil {
		svcClient, err = u.clientQueryRepo.FindByPublicID(ctx, *req.ClientID)
	} else if req.Name != nil {
		svcClient, err = u.clientQueryRepo.FindByName(ctx, *req.Name)
	}

	if err != nil {
		return nil, merr.New(merr.ErrNotFound, ErrClientNotFoundMsg, err)
	}

	if svcClient == nil {
		return nil, merr.New(merr.ErrNotFound, ErrClientNotFoundMsg, nil)
	}

	// Need to get service public ID from the service
	service, err := u.serviceQueryRepo.FindByID(ctx, svcClient.ServiceID())
	if err != nil {
		return nil, merr.New(merr.ErrNotFound, ErrServiceNotFoundMsg, err)
	}

	return &in.FindClientResult{
		ServiceID:      service.PublicID(),
		MgmtClientInfo: toClientInfo(svcClient),
	}, nil
}
