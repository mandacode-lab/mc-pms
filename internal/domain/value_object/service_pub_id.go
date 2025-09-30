package vo

import (
	"github.com/google/uuid"
)

type ServicePublicID struct {
	value uuid.UUID
}

func NewServicePublicID(value uuid.UUID) ServicePublicID {
	return ServicePublicID{value: value}
}

func ParseServicePublicID(value string) (ServicePublicID, error) {
	parsed, err := uuid.Parse(value)
	if err != nil {
		return ServicePublicID{}, err
	}
	return ServicePublicID{value: parsed}, nil
}

func (pid ServicePublicID) Value() uuid.UUID {
	return pid.value
}

func (pid ServicePublicID) String() string {
	return pid.value.String()
}
