package in

import (
	"time"

	clientappval "github.com/mandacode-com/mandacode-ssam/internal/domain/clientapp/value"
	serviceval "github.com/mandacode-com/mandacode-ssam/internal/domain/service/value"
)

type ServiceInfo struct {
	ServiceID serviceval.PublicID
	Name      string
	Desc      string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ClientAppInfo struct {
	ServiceID   serviceval.PublicID
	ClientAppID clientappval.PublicID
	Name        string
	Desc        string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
