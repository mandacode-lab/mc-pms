package vo

import (
	"github.com/google/uuid"
)

type ClientAppPublicID struct {
	value uuid.UUID
}

func NewClientAppPublicID(value uuid.UUID) ClientAppPublicID {
	return ClientAppPublicID{value: value}
}

func ParseClientAppPublicID(value string) (ClientAppPublicID, error) {
	parsed, err := uuid.Parse(value)
	if err != nil {
		return ClientAppPublicID{}, err
	}
	return ClientAppPublicID{value: parsed}, nil
}

func (pid ClientAppPublicID) Value() uuid.UUID {
	return pid.value
}

func (pid ClientAppPublicID) String() string {
	return pid.value.String()
}
