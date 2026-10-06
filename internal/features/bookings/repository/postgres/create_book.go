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
	book domain.Book,
) (domain.Book, error) {
	ctx, cancel := context.WithTimeout(ctx, r.Pool.OpTimeout())
	defer cancel()

	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return domain.Book{}, fmt.Errorf(
			"failed to create transaction: %w",
			err,
		)
	}

	defer tx.Rollback(ctx)

	query1 := `
		SELECT 1 FROM booking.users
		WHERE id = $1 AND deleted_at IS NULL
		FOR SHARE;
	`

	var one int
	row := tx.QueryRow(ctx, query1, book.UserID)

	err = row.Scan(&one)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Book{}, fmt.Errorf(
				"user with id='%d': %w",
				book.UserID,
				core_errors.ErrNotFound,
			)
		}
		return domain.Book{}, fmt.Errorf(
			"scan error: %w",
			err,
		)
	}

	query2 := `
		SELECT total_tickets FROM booking.events
		WHERE id = $1 AND canceled = FALSE AND event_time > now()
		FOR UPDATE;
	`

	var totalTickets int
	row = tx.QueryRow(ctx, query2, book.EventID)

	err = row.Scan(&totalTickets)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Book{}, fmt.Errorf(
				"event with id='%d' is canceled scan error: %w: event: %w",
				book.EventID,
				err,
				core_errors.ErrNotFound,
			)
		}
		return domain.Book{}, fmt.Errorf(
			"total tickets scan error: %w",
			err,
		)
	}

	query3 := `
		SELECT 
			COUNT(*) FILTER (WHERE user_id = $1),
			COUNT(*) 
		FROM booking.bookings
		WHERE event_id = $2 AND canceled_at IS NULL;
	`

	var userBooked, totalBooked int
	row = tx.QueryRow(ctx, query3, book.UserID, book.EventID)

	err = row.Scan(&userBooked, &totalBooked)
	if err != nil {
		return domain.Book{}, fmt.Errorf(
			"scan error: %w",
			err,
		)
	}

	if userBooked >= maxUserBooksForEvent {
		return domain.Book{}, fmt.Errorf(
			"no avaliable tickets for user: %w",
			core_errors.ErrConflict,
		)
	}

	if totalBooked >= totalTickets {
		return domain.Book{}, fmt.Errorf(
			"no avaliable tickets: %w",
			core_errors.ErrConflict,
		)
	}

	query4 := `
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
			canceled,
			canceled_at;
	`

	row = tx.QueryRow(
		ctx,
		query4,
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
		&model.Canceled,
		&model.CanceledAt,
	)
	if err != nil {
		return domain.Book{}, fmt.Errorf(
			"scan error: %w",
			err,
		)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return domain.Book{}, fmt.Errorf(
			"commit create booking error: %w",
			err,
		)
	}

	domain := domainFromModel(model)

	return domain, nil
}
