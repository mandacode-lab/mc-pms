package clientmgmt

import (
	"context"

	"github.com/mandacode-com/mandacode-ssam/internal/domain/entity"
	vo "github.com/mandacode-com/mandacode-ssam/internal/domain/vo"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
	"github.com/mandacode-com/merr"
)

func (u *Usecase) CreateClientApp(ctx context.Context, req *in.CreateClientAppRequest) (*in.CreateClientAppResponse, error) {
	// Validate service exists
	service, err := u.serviceQueryRepo.FindByPublicID(ctx, req.ServiceID)
	if err != nil {
		return nil, merr.New(merr.ErrNotFound, ErrServiceNotFoundMsg, err)
	}

	// Validate input
	if req.Name == "" {
		return nil, merr.New(merr.ErrBadRequest, ErrInvalidClientAppNameMsg, nil)
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
	desc, err := vo.NewClientAppDescription(req.Desc)
	if err != nil {
		return nil, merr.New(merr.ErrBadRequest, "invalid client app description", err)
	}

	clientAppEntity, err := entity.DraftClientApp(service.ID(), req.Name, desc, hash)
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}

	// Save to repository within transaction
	var savedClientApp *entity.ClientApp
	err = u.txManager.WithTx(ctx, func(tx out.Tx) error {
		savedClientApp, err = u.clientAppRepo.Create(ctx, tx, clientAppEntity)
		if err != nil {
			return merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Convert to result
	return &in.CreateClientAppResponse{
		ClientAppInfo: toClientAppInfo(savedClientApp, service.PublicID()),
		Secret:        plainSecret,
	}, nil
}
