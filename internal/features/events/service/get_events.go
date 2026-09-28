package events_service

import (
	"booking/internal/core/domain"
	"context"
	"fmt"
)

func (s *EventsService) GetEvents(
	ctx context.Context,
	filter domain.EventsQueryFilter,
) ([]domain.Event, error) {
	if err := filter.Validate(); err != nil {
		return nil, fmt.Errorf(
			"failed validation: %w",
			err,
		)
	}

	events, err := s.Repository.GetEvents(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get events from repository: %w",
			err,
		)
	}

	return events, nil
}
