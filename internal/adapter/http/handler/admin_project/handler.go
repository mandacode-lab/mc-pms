package admin_project

import (
	"github.com/mandacode-lab/mc-pms/internal/port/app"
)

type Handler struct {
	adminProject app.AdminProject
}

func NewHandler(adminProject app.AdminProject) *Handler {
	return &Handler{
		adminProject: adminProject,
	}
}

func toProjectResponse(p *app.AdminProjectResult) *ProjectResponse {
	return &ProjectResponse{
		ProjectID:   p.ProjectID,
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
