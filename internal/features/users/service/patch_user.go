package users_service

import (
	"booking/internal/core/domain"
	"context"
	"fmt"
)

func (s *UsersService) PatchUser(
	ctx context.Context,
	userID int,
	userPatch domain.PatchUser,
) (domain.User, error) {
	user, err := s.Repository.GetUser(ctx, userID)
	if err != nil {
		return domain.User{}, fmt.Errorf(
			"failed to get user with id='%d' from repository: %w",
			userID,
			err,
		)
	}

	err = user.AppyPatch(userPatch)
	if err != nil {
		return domain.User{}, fmt.Errorf(
			"failed to apply patch: %w",
			err,
		)
	}

	patchedUser, err := s.Repository.PatchUser(ctx, user)

	if err != nil {
		return domain.User{}, fmt.Errorf(
			"failed to patch user in repository: %w",
			err,
		)
	}

	return patchedUser, nil
}
