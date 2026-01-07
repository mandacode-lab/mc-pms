package client_project

import (
	"context"

	"github.com/mandacode-lab/mc-pms/internal/domain/project"
	"github.com/go-mandacode/merr"
)

func (a *Application) IsValidProjectNamespace(ctx context.Context, projectID string, namespaceID string) (bool, error) {
	proj, err := a.projectRepo.GetByID(ctx, project.NewProjectID(projectID))
	if err != nil {
		return false, merr.New(merr.ErrNotFound, "Project not found", err)
	}

	return proj.NamespaceID().String() == namespaceID, nil
}
