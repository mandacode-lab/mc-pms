package clientapp_mgmt

import (
	"context"

	"github.com/mandacode-com/mandacode-ssam/internal/domain/clientapp"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/merr"
)

func (u *Usecase) FindClientApp(ctx context.Context, req *in.FindClientAppRequest) (*in.FindClientAppResponse, error) {
	if req == nil {
		return nil, merr.New(merr.ErrBadRequest, "invalid request", nil)
	}

	if req.ClientAppID == nil && req.Name == nil {
		return nil, merr.New(merr.ErrBadRequest, "either client app ID or name must be provided", nil)
	}

	var clientApp *clientapp.ClientApp
	var err error

	if req.ClientAppID != nil {
		clientApp, err = u.clientAppQueryRepo.FindByPublicID(ctx, *req.ClientAppID)
	} else if req.Name != nil {
		clientApp, err = u.clientAppQueryRepo.FindByName(ctx, *req.Name)
	}

	if err != nil {
		return nil, merr.New(merr.ErrNotFound, ErrClientAppNotFoundMsg, err)
	}

	if clientApp == nil {
		return nil, merr.New(merr.ErrNotFound, ErrClientAppNotFoundMsg, nil)
	}

	// Need to get service public ID from the service
	serviceEntity, err := u.serviceQueryRepo.FindByID(ctx, clientApp.ServiceID())
	if err != nil {
		return nil, merr.New(merr.ErrNotFound, ErrServiceNotFoundMsg, err)
	}

	return &in.FindClientAppResponse{
		ClientAppInfo: toClientAppInfo(clientApp, serviceEntity.PublicID()),
	}, nil
}