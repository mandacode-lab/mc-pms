package out

import (
	"context"
)

type (
	Resource string
	Action   string
)

const (
	ResourceService       Resource = "service"
	ResourceServiceAccess Resource = "service_access"
)

const (
	ActionRead   Action = "read"
	ActionWrite  Action = "write"
	ActionDelete Action = "delete"
)

type Permission struct {
	Resource Resource
	Action   Action
	Context  map[string]any
}

func (p Permission) Validate() bool {
	validResources := map[Resource]bool{
		ResourceService: true,
	}

	validActions := map[Action]bool{
		ActionRead:   true,
		ActionWrite:  true,
		ActionDelete: true,
	}

	return validResources[p.Resource] && validActions[p.Action]
}

// IAMService defines methods for interacting with an Identity and Access Management system.
// TODO:
// Implement with basic auth to adapter
type IAMService interface {
	ReadResourcePermissions(ctx context.Context, userID string, resource Resource) ([]Permission, error)
	HasPermission(ctx context.Context, userID string, permission Permission) (bool, error)
}
