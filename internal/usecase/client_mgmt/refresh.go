package clientmgmt

import (
	"context"

	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
	"github.com/mandacode-com/merr"
)

func (u *Usecase) RefreshSecret(ctx context.Context, req *in.RefreshSecretRequest) (*in.RefreshSecretResponse, error) {
	// Find the client app by public ID
	clientApp, err := u.clientAppQueryRepo.FindByPublicID(ctx, req.ClientAppID)
	if err != nil {
		return nil, merr.New(merr.ErrNotFound, ErrClientAppNotFoundMsg, err)
	}

	// Generate secret bytes
	secretBytes, err := u.secretGen.Generate(ctx)
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}

	// Encode secret bytes to plain secret
	plainSecret, err := u.encoder.Encode(secretBytes)
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}

	// Hash the plain secret
	hash, err := u.hasher.Hash(ctx, plainSecret)
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}

	// Regenerate secret using domain logic
	clientApp.RegenerateSecret(hash)

	// Update in repository within transaction
	err = u.txManager.WithTx(ctx, func(tx out.Tx) error {
		return u.clientAppRepo.Update(ctx, tx, clientApp)
	})
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}

	return &in.RefreshSecretResponse{
		Secret: plainSecret,
	}, nil
}
