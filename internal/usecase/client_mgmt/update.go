package clientmgmt

import (
	"context"

	"github.com/mandacode-com/mandacode-ssam/internal/domain/tx"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/merr"
)

func (u *Usecase) UpdateClient(ctx context.Context, req *in.UpdateClientInput) (*in.UpdateClientResult, error) {
	// Find the client app by public ID
	svcClient, err := u.clientQueryRepo.FindByPublicID(ctx, req.ClientID)
	if err != nil {
		return nil, merr.New(merr.ErrNotFound, ErrClientNotFoundMsg, err)
	}

	// Apply updates using domain logic
	if req.NewName != nil {
		svcClient.UpdateName(*req.NewName)
	}

	if req.NewDesc != nil {
		svcClient.UpdateDescription(*req.NewDesc)
	}

	if req.NewIsActive != nil {
		if *req.NewIsActive {
			svcClient.Activate()
		} else {
			svcClient.Deactivate()
		}
	}

	// Update in repository within transaction
	err = u.txManager.WithTx(ctx, func(tx tx.Tx) error {
		return u.clientRepo.Update(ctx, tx, svcClient)
	})
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}

	// Convert to result
	return &in.UpdateClientResult{
		MgmtClientInfo: toClientInfo(svcClient),
	}, nil
}
