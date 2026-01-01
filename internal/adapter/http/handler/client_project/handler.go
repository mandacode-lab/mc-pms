package client_project

import (
	"github.com/mandacode-com/mandacode-pms/internal/port/app"
)

type Handler struct {
	clientProject app.ClientProject
}

func NewHandler(clientProject app.ClientProject) *Handler {
	return &Handler{
		clientProject: clientProject,
	}
}

func toProjectResponse(p *app.ClientProjectResult) *ProjectResponse {
	return &ProjectResponse{
		ProjectID:   p.ProjectID,
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		NamespaceID: p.NamespaceID,
	}
}
