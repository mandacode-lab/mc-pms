package userinfo

import (
	"github.com/mandacode-com/serengeti-integrated/internal/domain/shared"
)

const (
	UserInfoCreatedEventType = "userinfo.created"
	UserInfoUpdatedEventType = "userinfo.updated"
	UserInfoDeletedEventType = "userinfo.deleted"
)

type UserInfoCreatedEvent struct {
	shared.BaseEvent
	Nickname string `json:"nickname"`
}

func NewUserInfoCreatedEvent(aggregateID, nickname string) UserInfoCreatedEvent {
	return UserInfoCreatedEvent{
		BaseEvent: shared.NewBaseEvent(aggregateID, UserInfoCreatedEventType),
		Nickname:  nickname,
	}
}

type UserInfoUpdatedEvent struct {
	shared.BaseEvent
	UpdateType string `json:"update_type"`
}

func NewUserInfoUpdatedEvent(aggregateID, updateType string) UserInfoUpdatedEvent {
	return UserInfoUpdatedEvent{
		BaseEvent:  shared.NewBaseEvent(aggregateID, UserInfoUpdatedEventType),
		UpdateType: updateType,
	}
}

type UserInfoDeletedEvent struct {
	shared.BaseEvent
}

func NewUserInfoDeletedEvent(aggregateID string) UserInfoDeletedEvent {
	return UserInfoDeletedEvent{
		BaseEvent: shared.NewBaseEvent(aggregateID, UserInfoDeletedEventType),
	}
}
