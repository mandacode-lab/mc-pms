package in

import (
	"context"

	clientappval "github.com/mandacode-com/mandacode-ssam/internal/domain/clientapp/value"
	serviceval "github.com/mandacode-com/mandacode-ssam/internal/domain/service/value"
)

type VerifyClientRequest struct {
	ClientID     clientappval.PublicID
	ClientSecret []byte
}

type VerifyClientResponse struct {
	ServiceID serviceval.PublicID
}

type ClientAccessUsecase interface {
	VerifyClient(ctx context.Context, req *VerifyClientRequest) (*VerifyClientResponse, error)
}
