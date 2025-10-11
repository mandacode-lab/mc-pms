package vo

import (
	"errors"
	"strconv"
)

// ClientAppID is a type alias for client app internal ID (int64)
type ClientAppID int64

var (
	ErrInvalidClientAppID = errors.New("invalid client app ID")
)

// NewClientAppID creates a new ClientAppID with validation
func NewClientAppID(value int64) (ClientAppID, error) {
	if value <= 0 {
		return 0, ErrInvalidClientAppID
	}
	return ClientAppID(value), nil
}

// Int64 returns the int64 representation
func (id ClientAppID) Int64() int64 {
	return int64(id)
}

// String returns the string representation
func (id ClientAppID) String() string {
	return strconv.FormatInt(int64(id), 10)
}
