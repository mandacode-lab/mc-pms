package in

import (
	"context"
)

type FindUserInfoRequest struct {
	ServiceID *string
	UserID    *string
	Provider  *string
	Email     *string
	Nickname  *string
}

type FindUserInfoResponse struct {
	Users []UserInfo
}

type DeleteUserRequest struct {
	UserID string
}

type UserMgmtUsecase interface {
	FindUserInfo(ctx context.Context, req *FindUserInfoRequest) (*FindUserInfoResponse, error)
	DeleteUser(ctx context.Context, req *DeleteUserRequest) error
}
