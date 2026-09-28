package events_service

import (
	"booking/internal/core/domain"
	core_errors "booking/internal/core/errors"
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

	if event.Canceled {
		return domain.Event{}, fmt.Errorf(
			"imposible to patch canceled evet: %w",
			core_errors.ErrInvalidArgument,
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
