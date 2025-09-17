package client_access

import (
	"context"

	clientappval "github.com/mandacode-com/mandacode-ssam/internal/domain/clientapp/value"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/merr"
)

func (u *Usecase) VerifyClient(ctx context.Context, req *in.VerifyClientRequest) (*in.VerifyClientResponse, error) {
	// Parse client ID
	clientID, err := clientappval.ParsePublicID(req.ClientID)
	if err != nil {
		return nil, merr.New(merr.ErrBadRequest, ErrInvalidClientIDMsg, err)
	}

	// Find client app by public ID
	clientApp, err := u.clientAppQueryRepo.FindByPublicID(ctx, clientID)
	if err != nil {
		return nil, merr.New(merr.ErrNotFound, ErrClientAppNotFoundMsg, err)
	}

	// Check if client app is active
	if !clientApp.IsActive() {
		return &in.VerifyClientResponse{
			IsValid:       false,
			ServiceID:     "",
			ServiceInfo:   in.ServiceInfo{},
			ClientAppInfo: in.ClientAppInfo{},
		}, nil
	}

	// Verify secret
	err = u.hasher.Compare(ctx, clientApp.SecretHash(), req.ClientSecret)
	if err != nil {
		return &in.VerifyClientResponse{
			IsValid:       false,
			ServiceID:     "",
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
		return &in.VerifyClientResponse{
			IsValid:       false,
			ServiceID:     service.PublicID().String(),
			ServiceInfo:   in.ServiceInfo{},
			ClientAppInfo: in.ClientAppInfo{},
		}, nil
	}

	// Return valid result with service and client app info
	return &in.VerifyClientResponse{
		IsValid:       true,
		ServiceID:     service.PublicID().String(),
		ServiceInfo:   toServiceInfo(service),
		ClientAppInfo: toClientAppInfo(clientApp, service.PublicID()),
	}, nil
}
