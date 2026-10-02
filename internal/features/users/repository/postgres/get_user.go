package users_postgres_repository

import (
	"booking/internal/core/domain"
	core_errors "booking/internal/core/errors"
	core_postgres_pool "booking/internal/core/repository/postgres/pool"
	"context"
	"errors"
	"fmt"
)

func (r *UsersRepository) GetUser(
	ctx context.Context,
	userID int,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.Pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, version, name, phone_number, mail, created_at FROM booking.users
		WHERE id = $1 AND deleted_at IS NULL;
	`

	row := r.Pool.QueryRow(ctx, query, userID)

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
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.User{}, fmt.Errorf(
				"user with id='%d' %w: %v",
				userID,
				core_errors.ErrNotFound,
				err,
			)
		}
		return domain.User{}, fmt.Errorf(
			"select error: %w",
			err,
		)
	}

	userDomain := DomainFromModel(model)

	return userDomain, nil
}
