package initialize

import (
	"context"

	"github.com/mandacode-com/mandacode-ssam/internal/domain/clientapp"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/service"
	serviceval "github.com/mandacode-com/mandacode-ssam/internal/domain/service/value"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
	"github.com/mandacode-com/merr"
)

func (u *Usecase) InitializeSystem(ctx context.Context, req *in.InitializeSystemRequest) (*in.InitializeSystemResponse, error) {
	// Create Service name domain value
	serviceName, err := serviceval.NewName(req.ServiceName)
	if err != nil {
		return nil, merr.New(merr.ErrBadRequest, ErrInvalidServiceNameMsg, err)
	}

	// Check if serviceEntity with the same name already exists
	var serviceEntity *service.Service
	serviceEntity, err = u.serviceQueryRepo.FindByName(ctx, serviceName)
	if err != nil {
		return nil, merr.New(merr.ErrConflict, "failed to check existing service", err)
	}

	// Create new service if not exists
	if serviceEntity != nil {
		description := "System initialized service"

		// Create service draft
		serviceEntity = service.DraftService(serviceName, &description)

		// Create service with transaction
		err = u.txManager.WithTx(ctx, func(tx out.Tx) error {
			createdService, err := u.serviceRepo.Create(ctx, tx, serviceEntity)
			if err != nil {
				return merr.New(merr.ErrInternalServerError, "failed to create service", err)
			}
			serviceEntity = createdService
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	// Check if client app with the same name already exists for the service
	var clientAppEntity *clientapp.ClientApp
	clientAppEntity, err = u.clientAppQueryRepo.FindByName(ctx, req.ClientAppName)
	if err != nil {
		return nil, merr.New(merr.ErrConflict, "failed to check existing client app", err)
	}

	// Create new client app if not exists
	if clientAppEntity == nil {
		description := "System initialized client app"
		secretHash, err := u.hasher.Hash(ctx, req.ClientAppSecret)
		if err != nil {
			return nil, merr.New(merr.ErrInternalServerError, "failed to hash client app secret", err)
		}

		// Create client app draft
		clientAppDraft := clientapp.DraftClientApp(
			serviceEntity.ID(),
			req.ClientAppName,
			&description,
			secretHash,
		)

		// Create client app with transaction
		err = u.txManager.WithTx(ctx, func(tx out.Tx) error {
			createdClientApp, err := u.clientAppRepo.Create(ctx, tx, clientAppDraft)
			if err != nil {
				return merr.New(merr.ErrInternalServerError, ErrClientAppCreationFailedMsg, err)
			}
			clientAppEntity = createdClientApp
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return &in.InitializeSystemResponse{
		ServiceID:   serviceEntity.PublicID().String(),
		ClientAppID: clientAppEntity.PublicID().String(),
	}, nil
}
