package in

import (
	"context"

	vo "github.com/mandacode-com/mandacode-ssam/internal/domain/value_object"
)

type VerifyClientRequest struct {
	ClientID     vo.ClientAppPublicID
	ClientSecret []byte
}

type VerifyClientResponse struct {
	ServiceID vo.ServicePublicID
}

type ClientAccessUsecase interface {
	VerifyClient(ctx context.Context, req *VerifyClientRequest) (*VerifyClientResponse, error)
}
