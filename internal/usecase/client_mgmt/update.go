package clientmgmt

import (
	"context"

	vo "github.com/mandacode-com/mandacode-ssam/internal/domain/vo"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
	"github.com/mandacode-com/merr"
)

func (u *Usecase) UpdateClientApp(ctx context.Context, req *in.UpdateClientAppRequest) (*in.UpdateClientAppResponse, error) {
	// Find the client app by public ID
	clientApp, err := u.clientAppQueryRepo.FindByPublicID(ctx, req.ClientAppID)
	if err != nil {
		return nil, merr.New(merr.ErrNotFound, ErrClientAppNotFoundMsg, err)
	}

	// Get service for public ID conversion
	service, err := u.serviceQueryRepo.FindByID(ctx, clientApp.ServiceID())
	if err != nil {
		return nil, merr.New(merr.ErrNotFound, ErrServiceNotFoundMsg, err)
	}

	// Apply updates using domain logic
	if req.NewName != nil {
		clientApp.UpdateName(*req.NewName)
	}

	if req.NewDesc != nil {
		desc, err := vo.NewClientAppDescription(*req.NewDesc)
		if err != nil {
			return nil, merr.New(merr.ErrBadRequest, "invalid client app description", err)
		}
		clientApp.UpdateDescription(desc)
	}

	if req.NewIsActive != nil {
		if *req.NewIsActive {
			clientApp.Activate()
		} else {
			clientApp.Deactivate()
		}
	}

	// Update in repository within transaction
	err = u.txManager.WithTx(ctx, func(tx out.Tx) error {
		return u.clientAppRepo.Update(ctx, tx, clientApp)
	})
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}

	// Convert to result
	return &in.UpdateClientAppResponse{
		ClientAppInfo: toClientAppInfo(clientApp, service.PublicID()),
	}, nil
}
