package events_service

import (
	"booking/internal/core/domain"
	"context"
	"fmt"
)

func (s *EventsService) CreateEvent(
	ctx context.Context,
	event domain.Event,
) (domain.Event, error) {
	if err := event.Validate(); err != nil {
		return domain.Event{}, fmt.Errorf(
			"validation error: %w",
			err,
		)
	}

	eventDoamin, err := s.Repository.CreateEvent(ctx, event)
	if err != nil {
		return domain.Event{}, fmt.Errorf(
			"failed to create event: %w",
			err,
		)
	}

	return eventDoamin, nil
}
