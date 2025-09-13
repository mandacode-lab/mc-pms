package user_mgmt

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/serengeti-integrated/internal/domain/shared"
	serviceval "github.com/mandacode-com/serengeti-integrated/internal/domain/service/value"
	useridentityval "github.com/mandacode-com/serengeti-integrated/internal/domain/useridentity/value"
	"github.com/mandacode-com/serengeti-integrated/internal/port/in"
)

func (h *Handler) FindUserInfo(c *gin.Context) {
	ctx := c.Request.Context()

	var serviceID *serviceval.PublicID
	if serviceIDStr := c.Query("service_id"); serviceIDStr != "" {
		sid, err := serviceval.ParsePublicID(serviceIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid service_id format"})
			return
		}
		serviceID = &sid
	}

	var userID *useridentityval.PublicID
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		uid, err := useridentityval.ParsePublicID(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id format"})
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

	c.JSON(http.StatusOK, view)
}