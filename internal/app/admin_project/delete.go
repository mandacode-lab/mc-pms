package admin_project

import (
	"context"

	"github.com/mandacode-com/mandacode-pms/internal/domain/project"
)

func (a *Application) Delete(ctx context.Context, projectID string) error {
	return a.projectRepo.Delete(ctx, project.NewProjectID(projectID))
}
