package clientmgmt

import (
	"context"

	"github.com/mandacode-com/mandacode-ssam/internal/domain/tx"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/merr"
)

func (u *Usecase) RefreshSecret(ctx context.Context, req *in.RefreshSecretInput) (*in.RefreshSecretResult, error) {
	// Find the client app by public ID
	svcClient, err := u.clientQueryRepo.FindByPublicID(ctx, req.ClientID)
	if err != nil {
		return nil, merr.New(merr.ErrNotFound, ErrClientNotFoundMsg, err)
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

	// Regenerate secret using domain logic
	svcClient.RegenerateSecret(hash)

	// Update in repository within transaction
	err = u.txManager.WithTx(ctx, func(tx tx.Tx) error {
		return u.clientRepo.Update(ctx, tx, svcClient)
	})
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}

	return &in.RefreshSecretResult{
		RawSecret: plainSecret,
	}, nil
}
