package clientapp_mgmt

import (
	"context"

	"github.com/mandacode-com/merr"
	"github.com/mandacode-com/mandacode-service-hub/internal/domain/clientapp"
	clientappval "github.com/mandacode-com/mandacode-service-hub/internal/domain/clientapp/value"
	"github.com/mandacode-com/mandacode-service-hub/internal/port/in"
	"github.com/mandacode-com/mandacode-service-hub/internal/port/out"
	"github.com/mandacode-com/mandacode-service-hub/pkg/utils"
)

func (u *Usecase) CreateClientApp(ctx context.Context, cmd *in.CreateClientAppCommand) (*in.CreateClientAppView, error) {
	// Validate service exists
	service, err := u.serviceQueryRepo.FindByPublicID(ctx, cmd.ServiceID)
	if err != nil {
		return nil, merr.New(merr.ErrNotFound, ErrServiceNotFoundMsg, err)
	}

	// Validate input
	if cmd.Name == "" {
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
	desc := utils.StringNil(cmd.Desc)

	clientAppEntity, returnedSecret := clientapp.DraftClientApp(service.ID(), cmd.Name, desc, plainSecret, secretHash)

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

	// Convert to view
	return &in.CreateClientAppView{
		ClientAppInfo: toClientAppInfo(savedClientApp, cmd.ServiceID),
		Secret:        returnedSecret,
	}, nil
}
