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
	eventID int,
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

	queryEventCheck := `
		SELECT start_at <= now(), canceled_at IS NOT NULL
		FROM booking.events
		WHERE id = $1
		FOR UPDATE;
	`

	row := tx.QueryRow(ctx, queryEventCheck, eventID)

	var (
		passed          bool
		alreadyCanceled bool
	)

	err = row.Scan(
		&passed,
		&alreadyCanceled,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Event{}, fmt.Errorf(
				"event with id='%d': %w: %v",
				eventID,
				core_errors.ErrNotFound,
				err,
			)
		}
		return domain.Event{}, fmt.Errorf(
			"event canceled scan error: %w",
			err,
		)
	}

	if passed {
		return domain.Event{}, fmt.Errorf(
			"event with id='%d' is already passed: %w",
			eventID,
			core_errors.ErrConflict,
		)
	}

	if alreadyCanceled {
		return domain.Event{}, fmt.Errorf(
			"event with id='%d' is already canceled: %w",
			eventID,
			core_errors.ErrConflict,
		)
	}

	queryCancelEvent := `
		UPDATE booking.events
		SET 
			version = version + 1,
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
			start_at, 
			owners_user_id,
			canceled_at;
	`

	row = tx.QueryRow(ctx, queryCancelEvent, eventID)

	var model EventModel

	err = row.Scan(
		&model.ID,
		&model.Version,
		&model.Title,
		&model.Description,
		&model.TotalTickets,
		&model.CreatedAt,
		&model.Address,
		&model.StartsAt,
		&model.OwnerUserID,
		&model.CanceledAt,
	)
	if err != nil {
		return domain.Event{}, fmt.Errorf(
			"event scan error: %w",
			err,
		)
	}

	queryCancelEventsBookings := `
		WITH to_cancel AS (
			SELECT id
			FROM booking.bookings
			WHERE event_id = $1 AND canceled_at IS NULL
			ORDER BY id
			FOR UPDATE
		)
		UPDATE booking.bookings b
		SET canceled_at = now()
		FROM to_cancel
		WHERE b.id = to_cancel.id;
	`

	if _, err := tx.Exec(ctx, queryCancelEventsBookings, eventID); err != nil {
		return domain.Event{}, fmt.Errorf(
			"failed to cancel bookings: %w",
			err,
		)
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.Event{}, fmt.Errorf(
			"failed to commit: %w",
			err,
		)
	}

	event := DomainFromModel(model)

	return event, nil
}
