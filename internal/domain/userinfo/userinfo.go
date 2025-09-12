package userinfo

import (
	"time"

	"github.com/google/uuid"
	"github.com/mandacode-com/serengeti-integrated/internal/domain/shared"
	useridentityval "github.com/mandacode-com/serengeti-integrated/internal/domain/useridentity/value"
	userinfoval "github.com/mandacode-com/serengeti-integrated/internal/domain/userinfo/value"
)

type UserInfo struct {
	id             userinfoval.ID
	publicID       userinfoval.PublicID
	userIdentityID useridentityval.ID
	nickname       string
	email          string
	rawData        []byte
	createdAt      time.Time
	updatedAt      time.Time
	events         []shared.DomainEvent
}

func NewUserInfo(
	id userinfoval.ID,
	publicID userinfoval.PublicID,
	userIdentityID useridentityval.ID,
	nickname string,
	email string,
	rawData []byte,
	createdAt time.Time,
	updatedAt time.Time,
) *UserInfo {
	return &UserInfo{
		id:             id,
		publicID:       publicID,
		userIdentityID: userIdentityID,
		nickname:       nickname,
		email:          email,
		rawData:        rawData,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
		events:         make([]shared.DomainEvent, 0),
	}
}

func DraftUserInfo(
	userIdentityID useridentityval.ID,
	nickname string,
	email string,
	rawData []byte,
) *UserInfo {
	now := time.Now().UTC()
	id := userinfoval.NewID(0) // ID will be set by the database
	publicID := userinfoval.NewPublicID(uuid.New())

	ui := NewUserInfo(
		id,
		publicID,
		userIdentityID,
		nickname,
		email,
		rawData,
		now,
		now,
	)

	ui.raise(NewUserInfoCreatedEvent(ui.publicID.String(), nickname))
	return ui
}

// Getter methods

func (ui *UserInfo) ID() userinfoval.ID {
	return ui.id
}

func (ui *UserInfo) PublicID() userinfoval.PublicID {
	return ui.publicID
}

func (ui *UserInfo) UserIdentityID() useridentityval.ID {
	return ui.userIdentityID
}

func (ui *UserInfo) Nickname() string {
	return ui.nickname
}

func (ui *UserInfo) Email() string {
	return ui.email
}

func (ui *UserInfo) RawData() []byte {
	return ui.rawData
}

func (ui *UserInfo) CreatedAt() time.Time {
	return ui.createdAt
}

func (ui *UserInfo) UpdatedAt() time.Time {
	return ui.updatedAt
}

// Business methods

func (ui *UserInfo) UpdateNickname(nickname string) {
	if ui.nickname != nickname {
		ui.nickname = nickname
		ui.updatedAt = time.Now().UTC()
		ui.raise(NewUserInfoUpdatedEvent(ui.publicID.String(), "nickname_changed"))
	}
}

func (ui *UserInfo) UpdateEmail(email string) {
	ui.email = email
	ui.updatedAt = time.Now().UTC()
	ui.raise(NewUserInfoUpdatedEvent(ui.publicID.String(), "email_changed"))
}

func (ui *UserInfo) UpdateRawData(rawData []byte) {
	ui.rawData = rawData
	ui.updatedAt = time.Now().UTC()
	ui.raise(NewUserInfoUpdatedEvent(ui.publicID.String(), "raw_data_updated"))
}

// Event methods

func (ui *UserInfo) raise(e shared.DomainEvent) {
	ui.events = append(ui.events, e)
}

func (ui *UserInfo) PullEvents() []shared.DomainEvent {
	events := ui.events
	ui.events = nil
	return events
}
