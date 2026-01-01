package app

import (
	"context"
	"time"
)

type ClientProjectResult struct {
	ProjectID   string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	NamespaceID string
}

type ClientProject interface {
	// Query
	GetByID(ctx context.Context, projectID string) (*ClientProjectResult, error)
	IsValidProjectNamespace(ctx context.Context, projectID string, namespaceID string) (bool, error)
}
