package in

import (
	"time"
)

type ServiceInfo struct {
	ServiceID string
	Name      string
	Desc      string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ClientAppInfo struct {
	ServiceID   string
	ClientAppID string
	Name        string
	Desc        string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UserInfo struct {
	ServiceID string
	UserID    string
	Nickname  string
	Email     string
	Provider  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type WebOAuthInfo struct {
	ClientAppID   string
	Provider      string
	OAuthClientID string
	RedirectURI   string
	Scopes        []string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
