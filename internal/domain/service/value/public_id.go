package value

import (
	"errors"

	"github.com/google/uuid"
)

var ErrInvalidPublicID = errors.New("invalid service public ID")

type PublicID struct {
	value uuid.UUID
}

func NewPublicID(value uuid.UUID) PublicID {
	return PublicID{value: value}
}

func ParsePublicID(value string) (PublicID, error) {
	parsed, err := uuid.Parse(value)
	if err != nil {
		return PublicID{}, ErrInvalidPublicID
	}
	return PublicID{value: parsed}, nil
}

func (pid PublicID) Value() uuid.UUID {
	return pid.value
}

func (pid PublicID) String() string {
	return pid.value.String()
}
