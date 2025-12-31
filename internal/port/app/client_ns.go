package app

import (
	"context"
	"time"
)

type ClientNamespaceResult struct {
	NamespaceID string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ClientNamespace interface {
	// Query
	GetByID(ctx context.Context, namespaceID string) (*ClientNamespaceResult, error)
}
