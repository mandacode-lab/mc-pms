package servicemgmt

import (
	"time"

	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
)

type ServiceResponse struct {
	ServiceID   string    `json:"service_id" example:"srv_1234567890abcdef"`
	Name        string    `json:"name" example:"My Service"`
	Description string    `json:"description" example:"Service description"`
	IsActive    bool      `json:"is_active" example:"true"`
	CreatedAt   time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

func toServiceResponse(svc *in.ServiceView) ServiceResponse {
	return ServiceResponse{
		ServiceID:   svc.ServiceID.Value(),
		Name:        svc.Name.Value(),
		Description: svc.Desc.String(),
		IsActive:    svc.IsActive,
		CreatedAt:   svc.CreatedAt,
		UpdatedAt:   svc.UpdatedAt,
	}
}
