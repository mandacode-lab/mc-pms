package iam

import (
	"context"

	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
)

// NoOpIAMService is a no-operation IAM service that always allows all permissions
type NoOpIAMService struct{}

func NewNoOpIAMService() *NoOpIAMService {
	return &NoOpIAMService{}
}

func (s *NoOpIAMService) HasPermission(ctx context.Context, userID string, permission out.Permission) (bool, error) {
	// Always return true to allow all requests when IAM is disabled
	return true, nil
}
