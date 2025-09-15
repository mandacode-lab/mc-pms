package out

import "github.com/google/uuid"

type RoleService interface {
	InitializeRole(userID uuid.UUID) error
	DeleteRole(userID uuid.UUID) error
}
