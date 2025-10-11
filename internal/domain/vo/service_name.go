package vo

import (
	"errors"
	"strings"
)

var (
	ErrServiceNameEmpty   = errors.New("service name cannot be empty")
	ErrServiceNameTooLong = errors.New("service name is too long")
)

const (
	ServiceNameMaxLength = 255
)

// ServiceName is a type alias for service name with validation
type ServiceName string

// NewServiceName creates a new ServiceName with validation
func NewServiceName(name string) (ServiceName, error) {
	trimmed := strings.TrimSpace(name)

	if trimmed == "" {
		return "", ErrServiceNameEmpty
	}

	if len(trimmed) > ServiceNameMaxLength {
		return "", ErrServiceNameTooLong
	}

	return ServiceName(trimmed), nil
}

// String returns the string representation
func (n ServiceName) String() string {
	return string(n)
}
