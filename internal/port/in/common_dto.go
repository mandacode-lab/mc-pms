package in

import (
	"time"

	vo "github.com/mandacode-com/mandacode-ssam/internal/domain/value_object"
)

type ServiceInfo struct {
	ServiceID vo.ServicePublicID
	Name      string
	Desc      string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ClientAppInfo struct {
	ServiceID   vo.ServicePublicID
	ClientAppID vo.ClientAppPublicID
	Name        string
	Desc        string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
