package vo

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

// ClientAppPublicID is a type alias for client app public ID (UUID string)
type ClientAppPublicID string

var (
	ErrInvalidClientAppPublicIDFormat = errors.New("invalid client app public ID format")
	ErrClientAppPublicIDEmpty         = errors.New("client app public ID cannot be empty")
)

// NewClientAppPublicID creates a new ClientAppPublicID with validation
func NewClientAppPublicID(value string) (ClientAppPublicID, error) {
	trimmed := strings.TrimSpace(value)

	if trimmed == "" {
		return "", ErrClientAppPublicIDEmpty
	}

	// Validate UUID format
	if _, err := uuid.Parse(trimmed); err != nil {
		return "", ErrInvalidClientAppPublicIDFormat
	}

	return ClientAppPublicID(trimmed), nil
}

// ParseClientAppPublicID is an alias for NewClientAppPublicID (for backward compatibility)
func ParseClientAppPublicID(value string) (ClientAppPublicID, error) {
	return NewClientAppPublicID(value)
}

// String returns the string representation
func (id ClientAppPublicID) String() string {
	return string(id)
}

// UUID returns the uuid.UUID representation
func (id ClientAppPublicID) UUID() (uuid.UUID, error) {
	return uuid.Parse(string(id))
}
