package value

import (
	"errors"
	"strconv"
)

var ErrInvalidID = errors.New("invalid clientapp ID")

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