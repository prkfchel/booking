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

	queryCheckCanceledBooking := `
		SELECT canceled_at FROM booking.bookings
		WHERE id = $1 
		FOR UPDATE;
	`

	row := tx.QueryRow(ctx, queryCheckCanceledBooking, bookID)

	var canceledAt *time.Time

	err = row.Scan(
		&canceledAt,
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
			"book is already canceled: %w",
			core_errors.ErrConflict,
		)
	}

	queryCancelBooking := `
		UPDATE booking.bookings
		SET 
			canceled_at = $1
		WHERE id = $2
		RETURNING
			id,
			event_id,
			user_id,
			created_at,
			canceled_at;
	`

	row = tx.QueryRow(ctx, queryCancelBooking, time.Now(), bookID)

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
