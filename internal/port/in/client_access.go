package in

import (
	"context"

	clientappval "github.com/mandacode-com/mandacode-ssam/internal/domain/clientapp/value"
	serviceval "github.com/mandacode-com/mandacode-ssam/internal/domain/service/value"
)

type VerifyClientCommand struct {
	ClientID clientappval.PublicID
	Secret   []byte
}

type VerifyClientResult struct {
	IsValid   bool
	ServiceID serviceval.PublicID
	ServiceInfo
	ClientAppInfo
}

type ClientAccessUsecase interface {
	VerifyClient(ctx context.Context, cmd *VerifyClientCommand) (*VerifyClientResult, error)
}
