package value

import (
	"errors"
	"strings"
)

var (
	ErrNameEmpty   = errors.New("service name cannot be empty")
	ErrNameTooLong = errors.New("service name is too long")
)

type Name struct {
	value string
}

func NewName(name string) (Name, error) {
	trimmed := strings.TrimSpace(name)

	if trimmed == "" {
		return Name{}, ErrNameEmpty
	}

	if len(trimmed) > 255 {
		return Name{}, ErrNameTooLong
	}

	return Name{value: trimmed}, nil
}

func (n Name) Value() string {
	return n.value
}

func (n Name) String() string {
	return n.value
}
