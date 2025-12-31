package iam

// ResourceType represents the type of resource being accessed
type ResourceType string

const (
	// Namespace resources
	ResourceTypeNamespace ResourceType = "namespace"

	// Project resources
	ResourceTypeProject ResourceType = "project"
)

// Resource represents a specific resource instance
type Resource struct {
	Type ResourceType
	ID   string
}

func NewResource(resourceType ResourceType, id string) Resource {
	return Resource{
		Type: resourceType,
		ID:   id,
	}
}

func (r Resource) String() string {
	if r.ID == "*" {
		return string(r.Type) + ":*"
	}
	return string(r.Type) + ":" + r.ID
}
