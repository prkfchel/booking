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
		SELECT total_tickets FROM booking.events
		WHERE id = $1 AND canceled = FALSE
		FOR UPDATE;
	`

	var totalTickets int
	row := tx.QueryRow(ctx, query1, book.EventID)

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

	query2 := `
		SELECT COUNT(*) FROM booking.bookings
		WHERE event_id = $1 AND canceled = FALSE;
	`

	var soldTickets int
	row = tx.QueryRow(ctx, query2, book.EventID)

	err = row.Scan(&soldTickets)
	if err != nil {
		return domain.Book{}, fmt.Errorf(
			"sold tickets scan error: %w",
			err,
		)
	}

	if soldTickets >= totalTickets {
		return domain.Book{}, fmt.Errorf(
			"no avaliable tickets: %w",
			core_errors.ErrConflict,
		)
	}

	query3 := `
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
		query3,
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
