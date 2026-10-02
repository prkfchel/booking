package users_postgres_repository

import (
	"booking/internal/core/domain"
	core_errors "booking/internal/core/errors"
	core_postgres_pool "booking/internal/core/repository/postgres/pool"
	"context"
	"errors"
	"fmt"
)

func (r *UsersRepository) PatchUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.Pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE booking.users 
		SET version = version + 1,
			name = $1,
			phone_number = $2,
			mail = $3
		WHERE id = $4 AND version = $5 AND deleted_at IS NULL
		RETURNING id, version, name, phone_number, mail, created_at;
	`

	row := r.Pool.QueryRow(ctx, query, user.Name, user.PhoneNumber, user.Mail, user.ID, user.Version)

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
		if errors.Is(core_postgres_pool.ErrNoRows, err) {
			return domain.User{}, fmt.Errorf(
				"user with id='%d' concurrently accessed: %w",
				user.ID,
				core_errors.ErrConflict,
			)
		}

		return domain.User{}, fmt.Errorf(
			"scan error: %w",
			err,
		)
	}

	patchedUserDomain := DomainFromModel(model)

	return patchedUserDomain, nil
}
