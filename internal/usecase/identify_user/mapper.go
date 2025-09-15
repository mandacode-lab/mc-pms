package identify_user

import (
	serviceval "github.com/mandacode-com/mandacode-service-hub/internal/domain/service/value"
	"github.com/mandacode-com/mandacode-service-hub/internal/domain/useridentity"
	"github.com/mandacode-com/mandacode-service-hub/internal/domain/userinfo"
	"github.com/mandacode-com/mandacode-service-hub/internal/port/in"
)

func toUserIdentityResult(userIdentity *useridentity.UserIdentity, userInfo *userinfo.UserInfo, servicePublicID serviceval.PublicID, rawData map[string]any) *in.UserIdentityResult {
	return &in.UserIdentityResult{
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
