package vo

import (
	"errors"
	"strings"
)

var (
	ErrServiceNameEmpty   = errors.New("service name cannot be empty")
	ErrServiceNameTooLong = errors.New("service name is too long")
)

type ServiceName struct {
	value string
}

func NewServiceName(name string) (ServiceName, error) {
	trimmed := strings.TrimSpace(name)

	if trimmed == "" {
		return ServiceName{}, ErrServiceNameEmpty
	}

	if len(trimmed) > 255 {
		return ServiceName{}, ErrServiceNameTooLong
	}

	return ServiceName{value: trimmed}, nil
}

func (n ServiceName) Value() string {
	return n.value
}

func (n ServiceName) String() string {
	return n.value
}
