package clientapp_mgmt

import (
	"context"

	"github.com/mandacode-com/merr"
	"github.com/mandacode-com/serengeti-integrated/internal/port/in"
	"github.com/mandacode-com/serengeti-integrated/internal/port/out"
)

func (u *Usecase) UpdateClientApp(ctx context.Context, cmd *in.UpdateClientAppCommand) (*in.UpdateClientAppView, error) {
	// Find the client app by public ID
	clientApp, err := u.clientAppQueryRepo.FindByPublicID(ctx, cmd.ClientAppID)
	if err != nil {
		return nil, merr.New(merr.ErrNotFound, ErrClientAppNotFoundMsg, err)
	}

	// Get service for public ID conversion
	service, err := u.serviceQueryRepo.FindByID(ctx, clientApp.ServiceID())
	if err != nil {
		return nil, merr.New(merr.ErrNotFound, ErrServiceNotFoundMsg, err)
	}

	// Apply updates using domain logic
	if cmd.NewName != nil {
		clientApp.UpdateName(*cmd.NewName)
	}

	if cmd.NewDesc != nil {
		clientApp.UpdateDescription(cmd.NewDesc)
	}

	if cmd.NewIsActive != nil {
		if *cmd.NewIsActive {
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

	// Convert to view
	return &in.UpdateClientAppView{
		ClientAppInfo: toClientAppInfo(clientApp, service.PublicID()),
	}, nil
}
