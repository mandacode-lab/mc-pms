package value

import (
	"github.com/google/uuid"
)

type PublicID struct {
	value uuid.UUID
}

func NewPublicID(value uuid.UUID) PublicID {
	return PublicID{value: value}
}

func ParsePublicID(value string) (PublicID, error) {
	parsed, err := uuid.Parse(value)
	if err != nil {
		return PublicID{}, ErrInvalidID
	}
	return PublicID{value: parsed}, nil
}

func (pid PublicID) Value() uuid.UUID {
	return pid.value
}

func (pid PublicID) String() string {
	return pid.value.String()
}