package events_service

import (
	"booking/internal/core/domain"
	"context"
	"fmt"
)

func (s *EventsService) GetEvent(
	ctx context.Context,
	eventID int,
) (domain.Event, error) {
	event, err := s.Repository.GetEvent(ctx, eventID)
	if err != nil {
		return domain.Event{}, fmt.Errorf(
			"failed to get event from repository: %w",
			err,
		)
	}

	return event, nil
}
