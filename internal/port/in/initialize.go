package in

import "context"

type InitializeSystemRequest struct {
	ServiceName     string
	ClientAppName   string
	ClientAppID     string
	ClientAppSecret []byte
}

type InitializeSystemResponse struct {
	ServiceID   string
	ClientAppID string
}

type InitializeUsecase interface {
	InitializeSystem(ctx context.Context, req *InitializeSystemRequest) (*InitializeSystemResponse, error)
}

