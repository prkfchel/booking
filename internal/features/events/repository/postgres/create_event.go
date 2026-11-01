package events_postgres_repository

import (
	"booking/internal/core/domain"
	core_errors "booking/internal/core/errors"
	core_postgres_pool "booking/internal/core/repository/postgres/pool"
	"context"
	"errors"
	"fmt"
)

func (r *EventsRepository) CreateEvent(
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

	queryCheckOwner := `
		SELECT 1 FROM booking.users
		WHERE id = $1 AND deleted_at IS NULL
		FOR SHARE;
	`

	var one int
	row := tx.QueryRow(ctx, queryCheckOwner, event.OwnerUserID)

	err = row.Scan(&one)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Event{}, fmt.Errorf(
				"user with id='%d': %w",
				event.OwnerUserID,
				core_errors.ErrNotFound,
			)
		}

		return domain.Event{}, fmt.Errorf(
			"scan error: %w",
			err,
		)
	}

	queryCreateEvent := `
		INSERT INTO booking.events (
			title, 
			description, 
			total_tickets, 
			event_address, 
			start_at, 
			owners_user_id
		)
		VALUES($1, $2, $3, $4, $5, $6)
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
	row = tx.QueryRow(
		ctx,
		queryCreateEvent,
		event.Title,
		event.Description,
		event.TotalTickets,
		event.Address,
		event.StartsAt,
		event.OwnerUserID,
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
		if errors.Is(err, core_postgres_pool.ErrViolatesForeignKey) {
			return domain.Event{}, fmt.Errorf(
				"%v: user with id='%d': %w",
				err,
				event.OwnerUserID,
				core_errors.ErrNotFound,
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
			"commit create event error: %w",
			err,
		)
	}

	createdEvent := DomainFromModel(model)

	return createdEvent, nil

}
