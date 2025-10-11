package vo

import (
	"errors"
	"strings"
)

var (
	ErrClientAppDescriptionTooLong = errors.New("client app description is too long")
)

const (
	ClientAppDescriptionMaxLength = 500
)

// ClientAppDescription is a type alias for client app description text with validation
type ClientAppDescription string

// NewClientAppDescription creates a new ClientAppDescription with validation
func NewClientAppDescription(desc string) (ClientAppDescription, error) {
	trimmed := strings.TrimSpace(desc)

	// Empty is allowed for description
	if trimmed == "" {
		return "", nil
	}

	if len(trimmed) > ClientAppDescriptionMaxLength {
		return "", ErrClientAppDescriptionTooLong
	}

	return ClientAppDescription(trimmed), nil
}

// String returns the string representation
func (d ClientAppDescription) String() string {
	return string(d)
}

// Ptr returns a pointer to the string value, or nil if empty
func (d ClientAppDescription) Ptr() *string {
	if d == "" {
		return nil
	}
	s := string(d)
	return &s
}

// ClientAppDescriptionFromPtr creates a ClientAppDescription from a string pointer
func ClientAppDescriptionFromPtr(s *string) (ClientAppDescription, error) {
	if s == nil {
		return "", nil
	}
	return NewClientAppDescription(*s)
}
