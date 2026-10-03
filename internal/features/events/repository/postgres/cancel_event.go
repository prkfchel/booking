package events_postgres_repository

import (
	"booking/internal/core/domain"
	core_errors "booking/internal/core/errors"
	core_postgres_pool "booking/internal/core/repository/postgres/pool"
	"context"
	"errors"
	"fmt"
)

func (r *EventsRepository) CancelEvent(
	ctx context.Context,
	id int,
) (domain.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, r.Pool.OpTimeout())
	defer cancel()

	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return domain.Event{}, fmt.Errorf(
			"failed to create transaction: %w",
			err,
		)
	}
	defer tx.Rollback(ctx)

	query1 := `
		SELECT canceled, event_time <= now() FROM booking.events
		WHERE id = $1
		FOR UPDATE;
	`

	row := tx.QueryRow(ctx, query1, id)

	var canceled, passed bool

	err = row.Scan(
		&canceled,
		&passed,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Event{}, fmt.Errorf(
				"event with id='%d': %w",
				id,
				core_errors.ErrNotFound,
			)
		}
		return domain.Event{}, fmt.Errorf(
			"event canceled scan error: %w",
			err,
		)
	}

	if canceled {
		return domain.Event{}, fmt.Errorf(
			"event with id='%d' already canceled: %w",
			id,
			core_errors.ErrConflict,
		)
	}

	if passed {
		return domain.Event{}, fmt.Errorf(
			"event is already passed: %w",
			core_errors.ErrConflict,
		)
	}

	query2 := `
		UPDATE booking.events
		SET 
			version = version + 1,
			canceled = TRUE,
			canceled_at = now()
		WHERE id = $1
		RETURNING 
			id, 
			version, 
			title, 
			description, 
			total_tickets, 
			created_at, 
			event_address, 
			event_time, 
			owners_user_id,
			canceled,
			canceled_at;
	`

	row = tx.QueryRow(ctx, query2, id)

	var model EventModel

	err = row.Scan(
		&model.ID,
		&model.Version,
		&model.Title,
		&model.Description,
		&model.TotalTickets,
		&model.CreatedAt,
		&model.EventAddress,
		&model.EventTime,
		&model.OwnersUserId,
		&model.Canceled,
		&model.CanceledAt,
	)
	if err != nil {
		return domain.Event{}, fmt.Errorf(
			"event scan error: %w",
			err,
		)
	}

	query3 := `
		UPDATE booking.bookings b
		SET 
			canceled = TRUE,
			canceled_at = now()
		WHERE b.event_id = $1 AND b.canceled_at IS NULL;
	`

	if _, err := tx.Exec(ctx, query3, id); err != nil {
		return domain.Event{}, fmt.Errorf(
			"failed to cancel books: %w",
			err,
		)
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.Event{}, fmt.Errorf(
			"failed to commit: %w",
			err,
		)
	}

	eventDomain := DomainFromModel(model)

	return eventDomain, nil
}
