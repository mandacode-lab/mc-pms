package vo

import (
	"strconv"
)

type ServiceID struct {
	value int64
}

func NewServiceID(value int64) ServiceID {
	return ServiceID{value: value}
}

func (id ServiceID) Value() int64 {
	return id.value
}

func (id ServiceID) String() string {
	return strconv.FormatInt(id.value, 10)
}
