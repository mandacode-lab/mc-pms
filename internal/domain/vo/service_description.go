package vo

import (
	"errors"
	"strings"
)

var (
	ErrServiceDescriptionTooLong = errors.New("service description is too long")
)

const (
	ServiceDescriptionMaxLength = 500
)

// ServiceDescription is a type alias for service description text with validation
type ServiceDescription string

// NewServiceDescription creates a new ServiceDescription with validation
func NewServiceDescription(desc string) (ServiceDescription, error) {
	trimmed := strings.TrimSpace(desc)

	// Empty is allowed for description
	if trimmed == "" {
		return "", nil
	}

	if len(trimmed) > ServiceDescriptionMaxLength {
		return "", ErrServiceDescriptionTooLong
	}

	return ServiceDescription(trimmed), nil
}

// String returns the string representation
func (d ServiceDescription) String() string {
	return string(d)
}

// Ptr returns a pointer to the string value, or nil if empty
func (d ServiceDescription) Ptr() *string {
	if d == "" {
		return nil
	}
	s := string(d)
	return &s
}

// ServiceDescriptionFromPtr creates a ServiceDescription from a string pointer
func ServiceDescriptionFromPtr(s *string) (ServiceDescription, error) {
	if s == nil {
		return "", nil
	}
	return NewServiceDescription(*s)
}
