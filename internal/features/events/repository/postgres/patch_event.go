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

	queryPatchEvent := `
		UPDATE booking.events
		SET 
			version = version + 1,
			title = $1,
			description = $2,
			total_tickets = $3,
			event_address = $4,
			event_time = $5
			%s
		WHERE id = $6 AND version = $7
		RETURNING id, 
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

	args := []any{
		event.Title,
		event.Description,
		event.TotalTickets,
		event.Address,
		event.StartsAt,
		event.ID,
		event.Version,
	}

	if event.Canceled {
		queryPatchEvent = fmt.Sprintf(
			queryPatchEvent,
			", canceled = $8, canceled_at = $9",
		)
		args = append(args, event.Canceled, time.Now())
	} else {
		queryPatchEvent = fmt.Sprintf(queryPatchEvent, "")
	}

	row := r.Pool.QueryRow(
		ctx,
		queryPatchEvent,
		args...,
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
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Event{}, fmt.Errorf(
				"%v: event with id='%d': %w",
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

	patchedEvent := DomainFromModel(model)

	return patchedEvent, nil
}
