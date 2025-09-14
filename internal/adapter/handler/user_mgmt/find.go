package user_mgmt

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/merr"
	_ "github.com/mandacode-com/merr/middleware"
	serviceval "github.com/mandacode-com/serengeti/internal/domain/service/value"
	"github.com/mandacode-com/serengeti/internal/domain/shared"
	useridentityval "github.com/mandacode-com/serengeti/internal/domain/useridentity/value"
	"github.com/mandacode-com/serengeti/internal/port/in"
)

type UserResponse struct {
	ServiceID string    `json:"service_id" example:"srv_1234567890abcdef"`
	UserID    string    `json:"user_id" example:"usr_abcdef1234567890"`
	Nickname  string    `json:"nickname" example:"john_doe"`
	Email     string    `json:"email" example:"john.doe@example.com"`
	Provider  string    `json:"provider" example:"google"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

type FindUserInfoResponse struct {
	Users []UserResponse `json:"users"`
}

func toUserResponse(info in.UserInfo) UserResponse {
	return UserResponse{
		ServiceID: info.ServiceID.String(),
		UserID:    info.UserID.String(),
		Nickname:  info.Nickname,
		Email:     info.Email,
		Provider:  string(info.Provider),
		CreatedAt: info.CreatedAt,
		UpdatedAt: info.UpdatedAt,
	}
}

func toFindUserInfoResponse(view *in.FindUserInfoView) *FindUserInfoResponse {
	users := make([]UserResponse, len(view.Users))
	for i, info := range view.Users {
		users[i] = toUserResponse(info)
	}
	return &FindUserInfoResponse{
		Users: users,
	}
}

// FindUserInfo finds user information based on query parameters
// @Summary Find users
// @Description Find user information based on various filter criteria
// @Tags users
// @Produce json
// @Param service_id query string false "Service ID"
// @Param user_id query string false "User ID"
// @Param provider query string false "OAuth Provider" Enums(google, kakao, naver)
// @Param email query string false "User email"
// @Param nickname query string false "User nickname"
// @Success 200 {object} FindUserInfoResponse "List of matching users"
// @Failure 400 {object} merrmid.ErrorResponse "Invalid request"
// @Failure 500 {object} merrmid.ErrorResponse "Internal server error"
// @Router /users [get]
func (h *Handler) FindUserInfo(c *gin.Context) {
	ctx := c.Request.Context()

	var serviceID *serviceval.PublicID
	if serviceIDStr := c.Query("service_id"); serviceIDStr != "" {
		sid, err := serviceval.ParsePublicID(serviceIDStr)
		if err != nil {
			err := merr.New(merr.ErrBadRequest, "invalid service_id format", err)
			c.Error(err)
			return
		}
		serviceID = &sid
	}

	var userID *useridentityval.PublicID
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		uid, err := useridentityval.ParsePublicID(userIDStr)
		if err != nil {
			err := merr.New(merr.ErrBadRequest, "invalid user_id format", err)
			c.Error(err)
			return
		}
		userID = &uid
	}

	var provider *shared.Provider
	if providerStr := c.Query("provider"); providerStr != "" {
		p := shared.Provider(providerStr)
		if p.IsValid() {
			provider = &p
		}
	}

	var email *string
	if emailStr := c.Query("email"); emailStr != "" {
		email = &emailStr
	}

	var nickname *string
	if nicknameStr := c.Query("nickname"); nicknameStr != "" {
		nickname = &nicknameStr
	}

	query := &in.FindUserInfoQuery{
		ServiceID: serviceID,
		UserID:    userID,
		Provider:  provider,
		Email:     email,
		Nickname:  nickname,
	}

	view, err := h.userMgmt.FindUserInfo(ctx, query)
	if err != nil {
		c.Error(err)
		return
	}

	response := toFindUserInfoResponse(view)
	c.JSON(http.StatusOK, response)
}
