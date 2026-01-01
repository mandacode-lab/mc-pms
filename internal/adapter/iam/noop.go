package iam

import (
	"context"

	"github.com/mandacode-com/mandacode-pms/internal/port/iam"
)

// NoopIAMService is a no-operation IAM service that always allows access
// This is useful for development and testing
type NoopIAMService struct{}

func NewNoopIAMService() *NoopIAMService {
	return &NoopIAMService{}
}

// CheckAuthorization always returns allowed
func (s *NoopIAMService) CheckAuthorization(ctx context.Context, req *iam.AuthorizationRequest) (*iam.AuthorizationResult, error) {
	return &iam.AuthorizationResult{
		Allowed: true,
		Reason:  "",
	}, nil
}
