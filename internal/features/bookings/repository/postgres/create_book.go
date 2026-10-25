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

const maxUserBooksForEvent = 5

func (r *BookRepository) CreateBook(
	ctx context.Context,
	book domain.Booking,
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

	queryCheckUser := `
		SELECT 1 FROM booking.users
		WHERE id = $1 AND deleted_at IS NULL
		FOR SHARE;
	`

	var one int
	row := tx.QueryRow(ctx, queryCheckUser, book.UserID)

	err = row.Scan(&one)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Booking{}, fmt.Errorf(
				"user with id='%d': %w",
				book.UserID,
				core_errors.ErrNotFound,
			)
		}
		return domain.Booking{}, fmt.Errorf(
			"scan error: %w",
			err,
		)
	}

	queryGetTickets := `
		SELECT total_tickets, canceled_at FROM booking.events
		WHERE id = $1 AND event_time > now()
		FOR UPDATE;
	`

	var totalTickets int
	var canceledAt *time.Time
	row = tx.QueryRow(ctx, queryGetTickets, book.EventID)

	err = row.Scan(
		&totalTickets,
		&canceledAt,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Booking{}, fmt.Errorf(
				"event with id='%d': %w: %v",
				book.EventID,
				core_errors.ErrNotFound,
				err,
			)
		}
		return domain.Booking{}, fmt.Errorf(
			"total tickets scan error: %w",
			err,
		)
	}

	if canceledAt != nil {
		return domain.Booking{}, fmt.Errorf(
			"event with id='%d' is already canceled: %w",
			book.EventID,
			core_errors.ErrConflict,
		)
	}

	queryGetSoldTickets := `
		SELECT 
			COUNT(*) FILTER (WHERE user_id = $1),
			COUNT(*) 
		FROM booking.bookings
		WHERE event_id = $2 AND canceled_at IS NULL;
	`

	var userBooked, totalBooked int

	row = tx.QueryRow(ctx, queryGetSoldTickets, book.UserID, book.EventID)

	err = row.Scan(&userBooked, &totalBooked)
	if err != nil {
		return domain.Booking{}, fmt.Errorf(
			"scan error: %w",
			err,
		)
	}

	if userBooked >= maxUserBooksForEvent {
		return domain.Booking{}, fmt.Errorf(
			"no avaliable tickets for user: %w",
			core_errors.ErrConflict,
		)
	}

	if totalBooked >= totalTickets {
		return domain.Booking{}, fmt.Errorf(
			"no avaliable tickets: %w",
			core_errors.ErrConflict,
		)
	}

	queryCreateBooking := `
		INSERT INTO booking.bookings (
			event_id,
			user_id,
			created_at
		)
		VALUES($1, $2, $3)
		RETURNING
			id,
			event_id,
			user_id,
			created_at,
			canceled_at;
	`

	row = tx.QueryRow(
		ctx,
		queryCreateBooking,
		book.EventID,
		book.UserID,
		time.Now(),
	)

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
			"commit create booking error: %w",
			err,
		)
	}

	createdBook := domainFromModel(model)

	return createdBook, nil
}
