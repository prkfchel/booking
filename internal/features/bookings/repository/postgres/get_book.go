package bookings_postgres_repository

import (
	"booking/internal/core/domain"
	core_errors "booking/internal/core/errors"
	core_postgres_pool "booking/internal/core/repository/postgres/pool"
	"context"
	"errors"
	"fmt"
)

func (r *BookRepository) GetBook(
	ctx context.Context,
	bookID int,
) (domain.Booking, error) {
	ctx, cancel := context.WithTimeout(ctx, r.Pool.OpTimeout())
	defer cancel()

	queryGetBook := `
		SELECT 
			id,
			event_id,
			user_id,
			created_at,
			canceled,
			canceled_at
		FROM booking.bookings
		WHERE id = $1;
	`

	row := r.Pool.QueryRow(ctx, queryGetBook, bookID)

	var model BookModel

	err := row.Scan(
		&model.ID,
		&model.EventID,
		&model.UserID,
		&model.CreatedAt,
		&model.Canceled,
		&model.CanceledAt,
	)

	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Booking{}, fmt.Errorf(
				"book with id='%d' %w: %w",
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

	book := domainFromModel(model)

	return book, nil
}
