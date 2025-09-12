package userinfoval

import (
	"errors"
	"strconv"

	"github.com/google/uuid"
)

var (
	ErrInvalidID = errors.New("invalid userinfo ID")
)

type ID struct {
	value int64
}

func NewID(value int64) ID {
	return ID{value: value}
}

func (id ID) Value() int64 {
	return id.value
}

func (id ID) String() string {
	return strconv.FormatInt(id.value, 10)
}

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