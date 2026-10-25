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

	queryAnonymizeUser := `
		UPDATE booking.users
		SET 
			version = version + 1,
			name = 'deleted user',
			phone_number = NULL,
			mail = NULL,
			deleted_at = now()
		WHERE id = $1 AND deleted_at IS NULL;
	`

	cmdTag, err := tx.Exec(ctx, queryAnonymizeUser, userID)
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

	queryCancelEvent := `
		UPDATE booking.events
		SET 
			version = version + 1,
			canceled_at = now()
		WHERE owners_user_id = $1 AND canceled_at IS NULL AND event_time > now();
	`

	if _, err = tx.Exec(ctx, queryCancelEvent, userID); err != nil {
		return fmt.Errorf(
			"failed to cancel events: %w",
			err,
		)
	}

	queryCancelEventsBookings := `
		UPDATE booking.bookings b
		SET 
			canceled_at = now()
		FROM booking.events e
		WHERE b.event_id = e.id  
			AND e.owners_user_id = $1 
			AND b.canceled_at IS NULL;
	`

	if _, err = tx.Exec(ctx, queryCancelEventsBookings, userID); err != nil {
		return fmt.Errorf(
			"failed to cancel books: %w",
			err,
		)
	}

	queryCancelUsersBookings := `
		UPDATE booking.bookings b
		SET 
			canceled_at = now()
		WHERE b.user_id = $1 AND b.canceled_at IS NULL;
	`

	if _, err := tx.Exec(ctx, queryCancelUsersBookings, userID); err != nil {
		return fmt.Errorf(
			"failed to cancel users books: %w",
			err,
		)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf(
			"failed to commit delete user: %w",
			err,
		)
	}

	return nil
}
