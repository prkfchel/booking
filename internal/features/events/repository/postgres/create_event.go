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

func (r *EventsRepository) CreateEvent(
	ctx context.Context,
	event domain.Event,
) (domain.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, r.Pool.OpTimeout())
	defer cancel()

	queryCreateEvent := `
		INSERT INTO booking.events (
			title, 
			description, 
			total_tickets, 
			event_address, 
			event_time, 
			owners_user_id,
			created_at
		)
		VALUES($1, $2, $3, $4, $5, $6, $7)
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
	row := r.Pool.QueryRow(
		ctx,
		queryCreateEvent,
		event.Title,
		event.Description,
		event.TotalTickets,
		event.Address,
		event.StartsAt,
		event.OwnerUserID,
		time.Now(),
	)

	var model EventModel

	err := row.Scan(
		&model.ID,
		&model.Version,
		&model.Title,
		&model.Description,
		&model.TotalTickets,
		&model.CreatedAt,
		&model.Address,
		&model.StartsAt,
		&model.OwnerUserID,
		&model.Canceled,
		&model.CanceledAt,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrViolatesForeignKey) {
			return domain.Event{}, fmt.Errorf(
				"%v: user with id='%d': %w",
				err,
				event.ID,
				core_errors.ErrNotFound,
			)
		}
		return domain.Event{}, fmt.Errorf(
			"scan error: %w",
			err,
		)
	}

	createdEvent := DomainFromModel(model)

	return createdEvent, nil

}
