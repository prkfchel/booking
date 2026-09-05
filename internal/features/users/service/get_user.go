package users_service

import (
	"booking/internal/core/domain"
	"context"
	"fmt"
)

func (s *UsersService) GetUser(
	ctx context.Context,
	userID int,
) (domain.User, error) {
	user, err := s.Repository.GetUser(ctx, userID)
	if err != nil {
		return domain.User{}, fmt.Errorf(
			"failed to get user from repository: %w",
			err,
		)
	}

	return user, nil
}
