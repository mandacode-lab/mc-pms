package admin_project

import (
	"context"

	"github.com/mandacode-lab/mc-pms/internal/domain/project"
	"github.com/go-mandacode/merr"
)

func (a *Application) Delete(ctx context.Context, projectID string) error {
	if err := a.projectRepo.Delete(ctx, project.NewProjectID(projectID)); err != nil {
		return merr.New(merr.ErrInternalServerError, "Failed to delete project", err)
	}
	return nil
}
