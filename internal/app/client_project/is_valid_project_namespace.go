package client_project

import (
	"context"

	"github.com/mandacode-com/mandacode-project/internal/domain/project"
)

func (a *Application) IsValidProjectNamespace(ctx context.Context, projectID string, namespaceID string) (bool, error) {
	proj, err := a.projectRepo.GetByID(ctx, project.NewProjectID(projectID))
	if err != nil {
		return false, err
	}

	return proj.NamespaceID().String() == namespaceID, nil
}
