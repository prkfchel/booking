package bookings_postgres_repository

import (
	"booking/internal/core/domain"
	core_errors "booking/internal/core/errors"
	core_postgres_pool "booking/internal/core/repository/postgres/pool"
	"context"
	"errors"
	"fmt"
	"time"
)

func (r *BookRepository) CancelBook(
	ctx context.Context,
	bookID int,
) (domain.Booking, error) {
	ctx, cancel := context.WithTimeout(ctx, r.Pool.OpTimeout())
	defer cancel()

	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return domain.Booking{}, fmt.Errorf(
			"failed to create transaction: %w",
			err,
		)
	}
	defer tx.Rollback(ctx)

	queryCheckBooking := `
		SELECT b.canceled_at,
			e.start_at <= now() AS event_started
		FROM booking.bookings b
		JOIN booking.events e ON e.id = b.event_id
		WHERE b.id = $1
		FOR UPDATE OF b;
	`

	row := tx.QueryRow(ctx, queryCheckBooking, bookID)

	var (
		canceledAt   *time.Time
		eventStarted bool
	)

	err = row.Scan(
		&canceledAt,
		&eventStarted,
	)

	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Booking{}, fmt.Errorf(
				"booking with id='%d': %w: %v",
				bookID,
				core_errors.ErrNotFound,
				err,
			)
		}
		return domain.Booking{}, fmt.Errorf(
			"scan error: %w",
			err,
		)
	}

	if canceledAt != nil {
		return domain.Booking{}, fmt.Errorf(
			"booking is already canceled: %w",
			core_errors.ErrConflict,
		)
	}

	if eventStarted {
		return domain.Booking{}, fmt.Errorf(
			"event has already started, booking can't be canceled: %w",
			core_errors.ErrConflict,
		)
	}

	queryCancelBooking := `
		UPDATE booking.bookings
		SET 
			canceled_at = now()
		WHERE id = $1
		RETURNING
			id,
			event_id,
			user_id,
			created_at,
			canceled_at;
	`

	row = tx.QueryRow(ctx, queryCancelBooking, bookID)

	var model BookModel

	err = row.Scan(
		&model.ID,
		&model.EventID,
		&model.UserID,
		&model.CreatedAt,
		&model.CanceledAt,
	)
	if err != nil {
		return domain.Booking{}, fmt.Errorf(
			"scan error: %w",
			err,
		)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return domain.Booking{}, fmt.Errorf(
			"commit cancel booking error: %w",
			err,
		)
	}

	book := domainFromModel(model)

	return book, nil
}
