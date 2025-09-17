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

type PermissionAllowInfo struct {
	Allowed    bool
	Action     Action
	Version    string
}

// IAMService defines methods for interacting with an Identity and Access Management system.
type IAMService interface {
	HasPermission(ctx context.Context, userID string, permission Permission) (PermissionAllowInfo, error)
}
