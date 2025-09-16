package in

import (
	"context"
)

type VerifyClientRequest struct {
	ClientID     string
	ClientSecret []byte
}

type VerifyClientResponse struct {
	IsValid   bool
	ServiceID string
	ServiceInfo
	ClientAppInfo
}

type ClientAccessUsecase interface {
	VerifyClient(ctx context.Context, req *VerifyClientRequest) (*VerifyClientResponse, error)
}
