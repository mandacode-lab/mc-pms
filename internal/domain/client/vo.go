package client

import (
	"errors"
)

// ID is a type alias for client app internal ID (int64)
type ID int64

// NewID creates a new ClientAppID with validation
func NewID(value int64) (ID, error) {
	id := ID(value)
	if valid, err := id.IsValid(); !valid {
		return 0, err
	}
	return id, nil
}

func (id ID) Value() int64 {
	return int64(id)
}

func (id ID) IsValid() (bool, error) {
	if id <= 0 {
		return false, errors.New("client app ID must be greater than zero")
	}
	return true, nil
}

// PublicID
type PublicID string

func NewPublicID(value string) (PublicID, error) {
	publicID := PublicID(value)
	if valid, err := publicID.IsValid(); !valid {
		return "", err
	}
	return publicID, nil
}

func (id PublicID) Value() string {
	return string(id)
}

func (id PublicID) IsValid() (bool, error) {
	if id == "" {
		return false, errors.New("public ID cannot be empty")
	}
	return true, nil
}

// Name
type Name string

func NewName(value string) (Name, error) {
	name := Name(value)
	if valid, err := name.IsValid(); !valid {
		return "", err
	}
	return name, nil
}

func (n Name) Value() string {
	return string(n)
}

func (n Name) IsValid() (bool, error) {
	if n == "" {
		return false, errors.New("name cannot be empty")
	}
	return true, nil
}

// Description
type Description struct {
	value *string
}

func NewDescription(value string) (Description, error) {
	desc := Description{value: &value}
	if valid, err := desc.IsValid(); !valid {
		return Description{}, err
	}
	return desc, nil
}

func (d Description) Value() *string {
	return d.value
}

func (d Description) IsValid() (bool, error) {
	return true, nil
}

// SecretHash
type SecretHash []byte

func NewSecretHash(value []byte) (SecretHash, error) {
	secretHash := SecretHash(value)
	if valid, err := secretHash.IsValid(); !valid {
		return nil, err
	}
	return secretHash, nil
}

func (sh SecretHash) Value() []byte {
	return []byte(sh)
}

func (sh SecretHash) IsValid() (bool, error) {
	if len(sh) == 0 {
		return false, errors.New("secret hash cannot be empty")
	}
	return true, nil
}
