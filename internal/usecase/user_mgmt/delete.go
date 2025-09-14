package user_mgmt

import (
	"context"

	"github.com/mandacode-com/merr"
	"github.com/mandacode-com/serengeti/internal/port/in"
	"github.com/mandacode-com/serengeti/internal/port/out"
)

func (u *Usecase) DeleteUser(ctx context.Context, cmd *in.DeleteUserCommand) error {
	// Find UserIdentity by public ID
	userIdentity, err := u.userIdentityQueryRepo.FindByPublicID(ctx, cmd.UserID)
	if err != nil {
		return merr.New(merr.ErrNotFound, ErrUserNotFoundMsg, err)
	}

	// Find UserInfo by UserIdentity ID
	userInfo, err := u.userInfoQueryRepo.FindByUserIdentityID(ctx, userIdentity.ID())
	if err != nil {
		return merr.New(merr.ErrNotFound, ErrUserNotFoundMsg, err)
	}

	// Delete both UserInfo and UserIdentity within transaction
	err = u.txManager.WithTx(ctx, func(tx out.Tx) error {
		// Delete UserInfo first (foreign key constraint)
		err := u.userInfoRepo.Delete(ctx, tx, userInfo)
		if err != nil {
			return err
		}

		// Delete UserIdentity
		err = u.userIdentityRepo.Delete(ctx, tx, userIdentity)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}

	return nil
}
