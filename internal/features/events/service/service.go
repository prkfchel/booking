package events_service

import (
	"booking/internal/core/domain"
	"context"
)

type EventsService struct {
	Repository EventsRepository
}

func NewEventsService(resitory EventsRepository) *EventsService {
	return &EventsService{
		Repository: resitory,
	}
}

type EventsRepository interface {
	CreateEvent(
		ctx context.Context,
		event domain.Event,
	) (domain.Event, error)

	GetEvent(
		ctx context.Context,
		id int,
	) (domain.Event, error)

	PatchEvent(
		ctx context.Context,
		event domain.Event,
	) (domain.Event, error)

	GetEvents(
		ctx context.Context,
		filter domain.EventsQueryFilter,
	) ([]domain.Event, error)
}
