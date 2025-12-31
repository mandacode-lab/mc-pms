package admin_ns

import "time"

// NamespaceResponse is the common response DTO for namespace
type NamespaceResponse struct {
	NamespaceID string    `json:"namespace_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
