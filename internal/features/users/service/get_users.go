package users_service

import (
	"booking/internal/core/domain"
	"context"
	"fmt"
)

func (s *UsersService) GetUsers(
	ctx context.Context,
	limit *int,
	offset *int,
) ([]domain.User, error) {
	users, err := s.Repository.GetUsers(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get users from repository: %w",
			err,
		)
	}

	return users, nil
}
