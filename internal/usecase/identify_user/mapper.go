package identify_user

import (
	serviceval "github.com/mandacode-com/serengeti/internal/domain/service/value"
	"github.com/mandacode-com/serengeti/internal/domain/useridentity"
	"github.com/mandacode-com/serengeti/internal/domain/userinfo"
	"github.com/mandacode-com/serengeti/internal/port/in"
)

func toUserIdentityView(userIdentity *useridentity.UserIdentity, userInfo *userinfo.UserInfo, servicePublicID serviceval.PublicID, rawData map[string]any) *in.UserIdentityView {
	return &in.UserIdentityView{
		UserInfo: in.UserInfo{
			ServiceID: servicePublicID,
			UserID:    userIdentity.PublicID(),
			Nickname:  userInfo.Nickname(),
			Email:     userInfo.Email(),
			Provider:  userIdentity.Provider(),
			CreatedAt: userIdentity.CreatedAt(),
			UpdatedAt: userIdentity.UpdatedAt(),
		},
		RawData: rawData,
	}
}
