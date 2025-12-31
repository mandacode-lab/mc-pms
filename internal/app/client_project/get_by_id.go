package client_project

import (
	"context"

	"github.com/mandacode-com/mandacode-pms/internal/domain/project"
	"github.com/mandacode-com/mandacode-pms/internal/port/app"
)

func (a *Application) GetByID(ctx context.Context, projectID string) (*app.ClientProjectResult, error) {
	proj, err := a.projectRepo.GetByID(ctx, project.NewProjectID(projectID))
	if err != nil {
		return nil, err
	}

	return &app.ClientProjectResult{
		ProjectID:   proj.ID().String(),
		Name:        proj.Name().String(),
		Description: proj.Description().String(),
		CreatedAt:   proj.CreatedAt(),
		UpdatedAt:   proj.UpdatedAt(),
		NamespaceID: proj.NamespaceID().String(),
	}, nil
}
