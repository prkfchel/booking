package bookings_postgres_repository

import (
	"booking/internal/core/domain"
	"context"
	"fmt"
	"strings"
)

func (r *BookRepository) GetBooks(
	ctx context.Context,
	filter domain.BookingsFilter,
) ([]domain.Booking, error) {
	ctx, cancel := context.WithTimeout(ctx, r.Pool.OpTimeout())
	defer cancel()

	queryGetBookings := `
		SELECT 
			id,
			event_id,
			user_id,
			created_at,
			canceled_at
		FROM booking.bookings
		%s
		ORDER BY id ASC 
		LIMIT $1
		OFFSET $2;
	`

	args := []any{filter.Limit, filter.Offset}
	conditions := []string{}
	argsCount := 2

	if filter.UserID != nil {
		args = append(args, filter.UserID)
		argsCount++
		conditions = append(conditions, fmt.Sprintf("event_id = $%d", argsCount))
	}

	if filter.UserID != nil {
		args = append(args, filter.UserID)
		argsCount++
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argsCount))
	}

	if filter.DateFrom != nil {
		args = append(args, filter.DateFrom)
		argsCount++
		conditions = append(conditions, fmt.Sprintf("created_at > $%d", argsCount))
	}

	if filter.DateTo != nil {
		args = append(args, filter.DateTo)
		argsCount++
		conditions = append(conditions, fmt.Sprintf("created_at < $%d", argsCount))
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	queryGetBookings = fmt.Sprintf(queryGetBookings, where)

	rows, err := r.Pool.Query(ctx, queryGetBookings, args...)
	if err != nil {
		return nil, fmt.Errorf(
			"select error: %w",
			err,
		)
	}

	defer rows.Close()

	var models []BookModel

	for rows.Next() {
		var model BookModel

		err = rows.Scan(
			&model.ID,
			&model.EventID,
			&model.UserID,
			&model.CreatedAt,
			&model.CanceledAt,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"scan error: %w",
				err,
			)
		}

		models = append(models, model)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"next rows: %w",
			err,
		)
	}

	books := domainsFromModels(models)

	return books, nil

}
