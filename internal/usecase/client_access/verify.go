package clientaccess

import (
	"context"

	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/merr"
)

func (u *Usecase) VerifyClient(ctx context.Context, req *in.VerifyClientRequest) (*in.VerifyClientResponse, error) {
	// Find client app by public ID
	clientApp, err := u.clientAppQueryRepo.FindByPublicID(ctx, req.ClientID)
	if err != nil {
		return nil, merr.New(merr.ErrNotFound, ErrClientAppNotFoundMsg, err)
	}

	// Check if client app is active
	if !clientApp.IsActive() {
		return nil, merr.New(merr.ErrForbidden, ErrClientAppInactiveMsg, nil)
	}

	// Decode the base64-encoded client secret to get original secret bytes
	secretBytes, err := u.encoder.Decode(req.ClientSecret)
	if err != nil {
		return nil, merr.New(merr.ErrBadRequest, "invalid client_secret format", err)
	}

	// Verify secret: compare hash with decoded secret bytes
	err = u.hasher.Compare(ctx, clientApp.SecretHash(), secretBytes)
	if err != nil {
		return nil, merr.New(merr.ErrUnauthorized, ErrInvalidClientSecretMsg, err)
	}

	// Find service information
	service, err := u.serviceQueryRepo.FindByID(ctx, clientApp.ServiceID())
	if err != nil {
		return nil, merr.New(merr.ErrNotFound, ErrServiceNotFoundMsg, err)
	}

	// Check if service is active
	if !service.IsActive() {
		return nil, merr.New(merr.ErrForbidden, ErrServiceInactiveMsg, nil)
	}

	// Return valid result with service ID
	return &in.VerifyClientResponse{
		ServiceID: service.PublicID(),
	}, nil
}
