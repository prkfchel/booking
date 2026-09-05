package users_service

import (
	"booking/internal/core/domain"
	"context"
)

type UsersService struct {
	Repository UsersRepository
}

func NewUsersService(repository UsersRepository) *UsersService {
	return &UsersService{
		Repository: repository,
	}
}

type UsersRepository interface {
	CreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)

	GetUsers(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]domain.User, error)

	GetUser(
		ctx context.Context,
		userID int,
	) (domain.User, error)

	DeleteUser(
		ctx context.Context,
		userID int,
	) error

	PatchUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)
}
