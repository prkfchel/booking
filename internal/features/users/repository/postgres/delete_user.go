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

	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf(
			"failed to create transaction: %w",
			err,
		)
	}

	defer tx.Rollback(ctx)

	query1 := `
		UPDATE booking.users
		SET 
			version = version + 1,
			name = 'deleted user',
			phone_number = NULL,
			mail = NULL,
			deleted_at = now()
		WHERE id = $1 AND deleted_at IS NULL;
	`

	cmdTag, err := tx.Exec(ctx, query1, userID)
	if err != nil {
		return fmt.Errorf(
			"failed to delete user: %w",
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

	query2 := `
		UPDATE booking.events
		SET 
			version = version + 1,
			canceled = TRUE,
			canceled_at = now()
		WHERE owners_user_id = $1 AND canceled_at IS NULL;
	`

	if _, err = tx.Exec(ctx, query2, userID); err != nil {
		return fmt.Errorf(
			"failed to cancel events: %w",
			err,
		)
	}

	query3 := `
		UPDATE booking.bookings b
		SET 
			canceled = TRUE,
			canceled_at = now()
		FROM booking.events e
		WHERE b.event_id = e.id  
			AND e.owners_user_id = $1 
			AND b.canceled_at IS NULL;
	`

	if _, err = tx.Exec(ctx, query3, userID); err != nil {
		return fmt.Errorf(
			"failed to cancel books: %w",
			err,
		)
	}

	query4 := `
		UPDATE booking.bookings b
		SET 
			canceled = TRUE,
			canceled_at = now()
		WHERE b.user_id = $1 AND b.canceled_at IS NULL;
	`

	if _, err := tx.Exec(ctx, query4, userID); err != nil {
		return fmt.Errorf(
			"failed to cancel users books: %w",
			err,
		)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf(
			"failed to coomit delete user: %w",
			err,
		)
	}

	return nil
}
