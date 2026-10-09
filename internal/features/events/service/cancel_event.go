package events_service

import (
	"booking/internal/core/domain"
	"context"
	"fmt"
)

func (s *EventsService) CancelEvent(
	ctx context.Context,
	eventID int,
) (domain.Event, error) {
	event, err := s.Repository.CancelEvent(ctx, eventID)
	if err != nil {
		return domain.Event{}, fmt.Errorf(
			"failed to cancel event in repository: %w",
			err,
		)
	}

	return event, nil
}
