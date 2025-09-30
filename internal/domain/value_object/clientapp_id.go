package vo

import (
	"strconv"
)

type ClientAppID struct {
	value int64
}

func NewClientAppID(value int64) ClientAppID {
	return ClientAppID{value: value}
}

func (id ClientAppID) Value() int64 {
	return id.value
}

func (id ClientAppID) String() string {
	return strconv.FormatInt(id.value, 10)
}
