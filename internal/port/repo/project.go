package repo

import (
	"context"

	"github.com/mandacode-com/mandacode-project/internal/domain/ns"
	"github.com/mandacode-com/mandacode-project/internal/domain/project"
)

type ProjectFilter struct {
	NameContains string
	NamespaceID  ns.NamespaceID
}

type ProjectRepo interface {
	// Command
	Upsert(ctx context.Context, p *project.Project) error
	Delete(ctx context.Context, id project.ProjectID) error
	// Query
	GetByID(ctx context.Context, id project.ProjectID) (*project.Project, error)
	Find(ctx context.Context, filter ProjectFilter) ([]*project.Project, error)
	ListAll(ctx context.Context) ([]*project.Project, error)
}
