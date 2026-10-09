package events_postgres_repository

import (
	"booking/internal/core/domain"
	"context"
	"fmt"
	"strings"
)

func (r *EventsRepository) GetEvents(
	ctx context.Context,
	filter domain.EventsFilter,
) ([]domain.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, r.Pool.OpTimeout())
	defer cancel()

	queryGetEvents := `
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
			canceled_at
		FROM booking.events
		%s
		ORDER BY id ASC
		LIMIT $1
		OFFSET $2;
	`

	args := []any{filter.Limit, filter.Offset}
	conditions := []string{}
	argsCount := 2

	if filter.Title != nil {
		argsCount++
		conditions = append(conditions, fmt.Sprintf("title = $%d", argsCount))
		args = append(args, *filter.Title)
	}
	if filter.Address != nil {
		argsCount++
		conditions = append(conditions, fmt.Sprintf("event_address = $%d", argsCount))
		args = append(args, *filter.Address)
	}
	if filter.DateFrom != nil {
		argsCount++
		conditions = append(conditions, fmt.Sprintf("event_time > $%d", argsCount))
		args = append(args, *filter.DateFrom)
	}
	if filter.DateTo != nil {
		argsCount++
		conditions = append(conditions, fmt.Sprintf("event_time < $%d", argsCount))
		args = append(args, *filter.DateTo)
	}
	if filter.OwnerUserID != nil {
		argsCount++
		conditions = append(conditions, fmt.Sprintf("owners_user_id = $%d", argsCount))
		args = append(args, *filter.OwnerUserID)
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}
	queryGetEvents = fmt.Sprintf(queryGetEvents, where)

	rows, err := r.Pool.Query(ctx, queryGetEvents, args...)
	if err != nil {
		return nil, fmt.Errorf(
			"select error: %w",
			err,
		)
	}

	defer rows.Close()

	var models []EventModel

	for rows.Next() {
		var model EventModel

		err = rows.Scan(
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

	events := DomainsFromModels(models)

	return events, nil
}
