package client_access

import (
	"context"

	serviceval "github.com/mandacode-com/mandacode-ssam/internal/domain/service/value"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/merr"
)

func (u *Usecase) VerifyClient(ctx context.Context, cmd *in.VerifyClientCommand) (*in.VerifyClientResult, error) {
	// Find client app by public ID
	clientApp, err := u.clientAppQueryRepo.FindByPublicID(ctx, cmd.ClientID)
	if err != nil {
		return nil, merr.New(merr.ErrNotFound, ErrClientAppNotFoundMsg, err)
	}

	// Check if client app is active
	if !clientApp.IsActive() {
		return &in.VerifyClientResult{
			IsValid:       false,
			ServiceID:     serviceval.PublicID{},
			ServiceInfo:   in.ServiceInfo{},
			ClientAppInfo: in.ClientAppInfo{},
		}, nil
	}

	// Verify secret
	err = u.hasher.Compare(ctx, clientApp.SecretHash().Value(), cmd.Secret)
	if err != nil {
		return &in.VerifyClientResult{
			IsValid:       false,
			ServiceID:     serviceval.PublicID{},
			ServiceInfo:   in.ServiceInfo{},
			ClientAppInfo: in.ClientAppInfo{},
		}, nil
	}

	// Find service information
	service, err := u.serviceQueryRepo.FindByID(ctx, clientApp.ServiceID())
	if err != nil {
		return nil, merr.New(merr.ErrNotFound, ErrServiceNotFoundMsg, err)
	}

	// Check if service is active
	if !service.IsActive() {
		return &in.VerifyClientResult{
			IsValid:       false,
			ServiceID:     service.PublicID(),
			ServiceInfo:   in.ServiceInfo{},
			ClientAppInfo: in.ClientAppInfo{},
		}, nil
	}

	// Return valid result with service and client app info
	return &in.VerifyClientResult{
		IsValid:       true,
		ServiceID:     service.PublicID(),
		ServiceInfo:   toServiceInfo(service),
		ClientAppInfo: toClientAppInfo(clientApp, service.PublicID()),
	}, nil
}
