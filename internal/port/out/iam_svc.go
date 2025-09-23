package out

import (
	"context"
)

type (
	Resource string
	Action   string
)

const (
	ResourceService   Resource = "mrn::ssam:Service"
	ResourceClientApp Resource = "mrn::ssam:ClientApp"
)

const (
	ActionRead   Action = "ssam:read"
	ActionWrite  Action = "ssam:write"
	ActionDelete Action = "ssam:delete"
)

type Permission struct {
	Resource Resource
	Action   Action
	Context  map[string]any
}

func (p Permission) Validate() bool {
	validResources := map[Resource]bool{
		ResourceService:   true,
		ResourceClientApp: true,
	}

	validActions := map[Action]bool{
		ActionRead:   true,
		ActionWrite:  true,
		ActionDelete: true,
	}

	return validResources[p.Resource] && validActions[p.Action]
}

// IAMService defines methods for interacting with an Identity and Access Management system.
type IAMService interface {
	HasPermission(ctx context.Context, userID string, permission Permission) (bool, error)
}
