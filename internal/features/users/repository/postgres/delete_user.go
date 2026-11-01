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
		WHERE owners_user_id = $1 AND canceled_at IS NULL AND start_at > now();
	`

	if _, err = tx.Exec(ctx, queryCancelEvent, userID); err != nil {
		return fmt.Errorf(
			"failed to cancel events: %w",
			err,
		)
	}

	queryCancelBookings := `
		WITH to_cancel AS (
			SELECT b.id
			FROM booking.bookings b
			JOIN booking.events e ON e.id = b.event_id
			WHERE b.canceled_at IS NULL
				AND e.start_at > now()
				AND (b.user_id = $1 OR e.owners_user_id = $1)
			ORDER BY b.id
			FOR UPDATE OF b
		)
		UPDATE booking.bookings b
		SET canceled_at = now()
		FROM to_cancel
		WHERE b.id = to_cancel.id;
	`

	if _, err = tx.Exec(ctx, queryCancelBookings, userID); err != nil {
		return fmt.Errorf(
			"failed to cancel bookings: %w",
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
