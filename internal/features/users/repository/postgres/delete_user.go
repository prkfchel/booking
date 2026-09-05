package users_postgres_repository

import (
	core_errors "booking/internal/core/errors"
	"context"
	"fmt"
)

func (r *UsersRepository) DeleteUser(
	ctx context.Context,
	userID int,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.Pool.OpTimeout())
	defer cancel()

	query := `
		DELETE FROM booking.users
		WHERE id = $1;
	`

	cmdTag, err := r.Pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf(
			"exec query: %w",
			err,
		)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf(
			"user with id='%d': %w",
			userID,
			core_errors.ErrNotFound,
		)
	}

	return nil
}
