package events_postgres_repository

import (
	"booking/internal/core/domain"
	core_errors "booking/internal/core/errors"
	core_postgres_pool "booking/internal/core/repository/postgres/pool"
	"context"
	"errors"
	"fmt"
)

func (r *EventsRepository) GetEvent(
	ctx context.Context,
	id int,
) (domain.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, r.Pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, 
			version, 
			title, 
			description, 
			total_tickets, 
			created_at, 
			event_address, 
			event_time, 
			owners_user_id,
			canceled,
			canceled_at FROM booking.events
		WHERE id = $1;
	`

	row := r.Pool.QueryRow(ctx, query, id)

	var model EventModel

	err := row.Scan(
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
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Event{}, fmt.Errorf(
				"event with id='%d' %w: %v",
				id,
				core_errors.ErrNotFound,
				err,
			)
		}
		return domain.Event{}, fmt.Errorf(
			"scan error: %w",
			err,
		)
	}

	eventDomain := DomainFromModel(model)

	return eventDomain, nil
}
