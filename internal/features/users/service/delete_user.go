package users_service

import (
	core_errors "booking/internal/core/errors"
	"context"
	"fmt"
)

func (s *UsersService) DeleteUser(
	ctx context.Context,
	userID int,
) error {
	if _, err := s.Repository.GetUser(ctx, userID); err != nil {
		return fmt.Errorf(
			"user with id='%d': %w",
			userID,
			core_errors.ErrNotFound,
		)
	}

	err := s.Repository.DeleteUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("delete user from repository: %w", err)
	}

	return nil
}
