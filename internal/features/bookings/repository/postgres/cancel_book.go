package bookings_postgres_repository

import (
	"booking/internal/core/domain"
	core_errors "booking/internal/core/errors"
	"context"
	"fmt"
	"time"
)

func (r *BookRepository) CancelBook(
	ctx context.Context,
	id int,
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
		SELECT canceled FROM booking.bookings
		WHERE id = $1
	`

	row := tx.QueryRow(ctx, query1, id)

	var canceled bool

	err = row.Scan(
		&canceled,
	)

	if err != nil {
		return domain.Book{}, fmt.Errorf(
			"scan error: %w",
			err,
		)
	}

	if canceled {
		return domain.Book{}, fmt.Errorf(
			"book is already canceled: %w",
			core_errors.ErrConflict,
		)
	}

	query2 := `
		UPDATE booking.bookings
		SET 
			canceled = true,
			canceled_at = $1
		WHERE id = $2
		RETURNING
			id,
			event_id,
			user_id,
			created_at,
			canceled,
			canceled_at;
	`

	row = tx.QueryRow(ctx, query2, time.Now(), id)

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
		)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return domain.Book{}, fmt.Errorf(
			"commit cancel booking error: %w",
			err,
		)
	}

	book := domainFromModel(model)

	return book, nil
}
