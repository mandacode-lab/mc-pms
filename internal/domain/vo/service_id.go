package vo

import (
	"errors"
	"strconv"
)

// ServiceID is a type alias for service internal ID (int64)
type ServiceID int64

var (
	ErrInvalidServiceID = errors.New("invalid service ID")
)

// NewServiceID creates a new ServiceID with validation
func NewServiceID(value int64) (ServiceID, error) {
	if value <= 0 {
		return 0, ErrInvalidServiceID
	}
	return ServiceID(value), nil
}

// Int64 returns the int64 representation
func (id ServiceID) Int64() int64 {
	return int64(id)
}

// String returns the string representation
func (id ServiceID) String() string {
	return strconv.FormatInt(int64(id), 10)
}
