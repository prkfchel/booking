package users_service

import (
	"context"
	"fmt"
)

func (s *UsersService) DeleteUser(
	ctx context.Context,
	userID int,
) error {
	err := s.Repository.DeleteUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("delete user from repository: %w", err)
	}

	return nil
}
