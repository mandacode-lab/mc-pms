package clientapp_mgmt

import (
	"context"

	"github.com/mandacode-com/mandacode-ssam/internal/domain/clientapp"
	clientappval "github.com/mandacode-com/mandacode-ssam/internal/domain/clientapp/value"
	serviceval "github.com/mandacode-com/mandacode-ssam/internal/domain/service/value"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
	"github.com/mandacode-com/mandacode-ssam/pkg/utils"
	"github.com/mandacode-com/merr"
)

func (u *Usecase) CreateClientApp(ctx context.Context, req *in.CreateClientAppRequest) (*in.CreateClientAppResponse, error) {
	// Parse service ID
	serviceID, err := serviceval.ParsePublicID(req.ServiceID)
	if err != nil {
		return nil, merr.New(merr.ErrBadRequest, ErrInvalidServiceIDMsg, err)
	}

	// Validate service exists
	service, err := u.serviceQueryRepo.FindByPublicID(ctx, serviceID)
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

	// Encode secret bytes to plain secret
	plainSecret := u.encoder.Encode(secretBytes)

	// Hash the plain secret
	hash, err := u.hasher.Hash(ctx, plainSecret)
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}

	// Create secret hash from the computed hash
	secretHash := clientappval.NewSecretHash(hash)

	// Create new client app using domain logic
	desc := utils.StringNil(req.Desc)

	clientAppEntity, returnedSecret := clientapp.DraftClientApp(service.ID(), req.Name, desc, plainSecret, secretHash)

	// Save to repository within transaction
	var savedClientApp *clientapp.ClientApp
	err = u.txManager.WithTx(ctx, func(tx out.Tx) error {
		saved, err := u.clientAppRepo.Create(ctx, tx, clientAppEntity)
		if err != nil {
			return merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
		}
		savedClientApp = saved
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Convert to result
	return &in.CreateClientAppResponse{
		ClientAppInfo: toClientAppInfo(savedClientApp, service.PublicID()),
		Secret:        returnedSecret,
	}, nil
}
