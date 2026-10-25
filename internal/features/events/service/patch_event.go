package events_service

import (
	"booking/internal/core/domain"
	"context"
	"fmt"
)

func (s *EventsService) PatchEvent(
	ctx context.Context,
	id int,
	patch domain.PatchEvent,
) (domain.Event, error) {
	event, err := s.Repository.GetEvent(ctx, id)
	if err != nil {
		return domain.Event{}, fmt.Errorf(
			"failed to get user with id='%d' from repository: %w",
			id,
			err,
		)
	}

	err = event.ApplyPatch(patch)
	if err != nil {
		return domain.Event{}, fmt.Errorf(
			"failed to apply patch: %w",
			err,
		)
	}

	patchedEvent, err := s.Repository.PatchEvent(ctx, event)
	if err != nil {
		return domain.Event{}, fmt.Errorf(
			"failed to patch event in repository: %w",
			err,
		)
	}

	return patchedEvent, nil
}
