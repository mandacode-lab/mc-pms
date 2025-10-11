package vo

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

// ServicePublicID is a type alias for service public ID (UUID string)
type ServicePublicID string

var (
	ErrInvalidServicePublicIDFormat = errors.New("invalid service public ID format")
	ErrServicePublicIDEmpty         = errors.New("service public ID cannot be empty")
)

// NewServicePublicID creates a new ServicePublicID with validation
func NewServicePublicID(value string) (ServicePublicID, error) {
	trimmed := strings.TrimSpace(value)

	if trimmed == "" {
		return "", ErrServicePublicIDEmpty
	}

	// Validate UUID format
	if _, err := uuid.Parse(trimmed); err != nil {
		return "", ErrInvalidServicePublicIDFormat
	}

	return ServicePublicID(trimmed), nil
}

// String returns the string representation
func (id ServicePublicID) String() string {
	return string(id)
}

// UUID returns the uuid.UUID representation
func (id ServicePublicID) UUID() (uuid.UUID, error) {
	return uuid.Parse(string(id))
}
