package events_postgres_repository

import (
	"booking/internal/core/domain"
	core_errors "booking/internal/core/errors"
	core_postgres_pool "booking/internal/core/repository/postgres/pool"
	"context"
	"errors"
	"fmt"
	"time"
)

func (r *EventsRepository) PatchEvent(
	ctx context.Context,
	event domain.Event,
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

	queryCheckEvent := `
		SELECT canceled_at, start_at <= now()
		FROM booking.events
		WHERE id = $1
		FOR UPDATE;
	`

	var (
		canceledAt *time.Time
		passed     bool
	)
	row := tx.QueryRow(ctx, queryCheckEvent, event.ID)

	err = row.Scan(
		&canceledAt,
		&passed,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Event{}, fmt.Errorf(
				"event with id ='%d': %w",
				event.ID,
				core_errors.ErrNotFound,
			)
		}

		return domain.Event{}, fmt.Errorf(
			"scan error: %w",
			err,
		)
	}

	if canceledAt != nil {
		return domain.Event{}, fmt.Errorf(
			"you can't patch canceled event: %w",
			core_errors.ErrConflict,
		)
	}

	if passed {
		return domain.Event{}, fmt.Errorf(
			"event already passed: %w",
			core_errors.ErrConflict,
		)
	}

	queryBookedTickets := `
		SELECT COUNT(*) FROM booking.bookings 
		WHERE event_id = $1 AND canceled_at IS NULL;
	`
	var bookedTickets int
	row = tx.QueryRow(ctx, queryBookedTickets, event.ID)

	if err = row.Scan(&bookedTickets); err != nil {
		return domain.Event{}, fmt.Errorf(
			"scan error: %w",
			err,
		)
	}

	if bookedTickets > event.TotalTickets {
		return domain.Event{}, fmt.Errorf(
			"already booked tickets more than total: %w",
			core_errors.ErrConflict,
		)
	}

	queryPatchEvent := `
		UPDATE booking.events
		SET 
			version = version + 1,
			title = $1,
			description = $2,
			total_tickets = $3,
			event_address = $4,
			start_at = $5
		WHERE id = $6 AND version = $7 AND canceled_at IS NULL
		RETURNING id, 
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

	args := []any{
		event.Title,
		event.Description,
		event.TotalTickets,
		event.Address,
		event.StartsAt,
		event.ID,
		event.Version,
	}

	row = tx.QueryRow(
		ctx,
		queryPatchEvent,
		args...,
	)

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
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Event{}, fmt.Errorf(
				"%v: event with id='%d': %w",
				err,
				event.ID,
				core_errors.ErrConflict,
			)
		}
		if errors.Is(err, core_postgres_pool.ErrViolatesCheck) {
			return domain.Event{}, fmt.Errorf(
				"violates check constraint: %v: %w",
				err,
				core_errors.ErrInvalidArgument,
			)
		}
		return domain.Event{}, fmt.Errorf(
			"scan error: %w",
			err,
		)
	}

	if err = tx.Commit(ctx); err != nil {
		return domain.Event{}, fmt.Errorf(
			"failed to commit patch event: %w",
			err,
		)
	}

	patchedEvent := DomainFromModel(model)

	return patchedEvent, nil
}
