package clientmgmt

import (
	"context"
	"time"

	"github.com/mandacode-com/mandacode-ssam/internal/domain/client"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/tx"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/merr"
)

func (u *Usecase) CreateClient(ctx context.Context, req *in.CreateClientInput) (*in.CreateClientResult, error) {
	// Validate service exists
	service, err := u.serviceQueryRepo.FindByPublicID(ctx, req.ServiceID)
	if err != nil {
		return nil, merr.New(merr.ErrNotFound, ErrServiceNotFoundMsg, err)
	}

	// Validate input
	if req.Name == "" {
		return nil, merr.New(merr.ErrBadRequest, ErrInvalidClientNameMsg, nil)
	}

	// Generate secret bytes
	secretBytes, err := u.secretGen.Generate(ctx)
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}

	// Hash the secret bytes
	hash, err := u.hasher.Hash(ctx, secretBytes)
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}

	// Encode secret bytes to plain secret (for response only)
	plainSecret, err := u.encoder.Encode(secretBytes)
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}

	// Create new client app using domain logic
	id, err := client.NewID(1) // ID will be set by repository
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}
	publicIDStr, err := u.clientIDGen.Generate()
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}
	publicID, err := client.NewPublicID(publicIDStr)
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}
	now := time.Now().UTC()

	svcClient := client.NewClient(
		id,
		publicID,
		req.Name,
		req.Desc,
		hash,
		service.ID(),
		true,
		now,
		now,
	)
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}

	// Save to repository within transaction
	var savedClient *client.Client
	err = u.txManager.WithTx(ctx, func(tx tx.Tx) error {
		savedClient, err = u.clientRepo.Create(ctx, tx, svcClient)
		if err != nil {
			return merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Convert to result
	return &in.CreateClientResult{
		Secret:         plainSecret,
		ServiceID:      service.PublicID(),
		MgmtClientInfo: toClientInfo(savedClient),
	}, nil
}
