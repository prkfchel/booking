package users_postgres_repository

import (
	"booking/internal/core/domain"
	core_errors "booking/internal/core/errors"
	core_postgres_pool "booking/internal/core/repository/postgres/pool"
	"context"
	"errors"
	"fmt"
)

func (r *UsersRepository) CreateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.Pool.OpTimeout())
	defer cancel()

	queryCreateUser := `
		INSERT INTO booking.users (name, phone_number, mail)
		VALUES ($1, $2, $3)
		RETURNING id, version, name, phone_number, mail, created_at;
	`

	row := r.Pool.QueryRow(
		ctx,
		queryCreateUser,
		user.Name,
		user.PhoneNumber,
		user.Mail,
	)

	var model UserModel

	err := row.Scan(
		&model.ID,
		&model.Version,
		&model.Name,
		&model.PhoneNumber,
		&model.Mail,
		&model.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrViolatesUnique) {
			return domain.User{}, fmt.Errorf(
				"phone number or mail is already in use: %v: %w",
				err,
				core_errors.ErrConflict,
			)
		}
		if errors.Is(err, core_postgres_pool.ErrViolatesCheck) {
			return domain.User{}, fmt.Errorf(
				"violates check constraint: %v: %w",
				err,
				core_errors.ErrInvalidArgument,
			)
		}
		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	createdUser := DomainFromModel(model)

	return createdUser, nil
}
