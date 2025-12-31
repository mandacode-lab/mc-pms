package app

import (
	"context"
	"time"
)

type AdminProjectResult struct {
	ProjectID   string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateProjectInput struct {
	NamespaceID string
	Name        string
	Description string
}

type UpdateProjectInput struct {
	ProjectID   string
	Name        *string
	Description *string
}

type FindProjectInput struct {
	NameContains string
	NamespaceID  string
}

type AdminProject interface {
	// Command
	Create(ctx context.Context, input CreateProjectInput) (*AdminProjectResult, error)
	Update(ctx context.Context, input UpdateProjectInput) (*AdminProjectResult, error)
	Delete(ctx context.Context, projectID string) error
	// Query
	GetByID(ctx context.Context, projectID string) (*AdminProjectResult, error)
	Find(ctx context.Context, input FindProjectInput) ([]*AdminProjectResult, error)
	ListAll(ctx context.Context) ([]*AdminProjectResult, error)
}
