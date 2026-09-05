package users_service

import (
	"booking/internal/core/domain"
	"context"
	"fmt"
)

func (s *UsersService) CreateUser(
	user domain.User,
	ctx context.Context,
) (domain.User, error) {
	if err := user.Validate(); err != nil {
		return domain.User{}, fmt.Errorf(
			"validation error: %w",
			err,
		)
	}

	user, err := s.Repository.CreateUser(ctx, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("create user repository error: %w", err)
	}

	return user, nil
}
