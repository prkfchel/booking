package events_http_transport

import (
	"booking/internal/core/domain"
	core_http_server "booking/internal/core/transport/http/server"
	"context"
	"net/http"
)

type EventsHTTPHandler struct {
	Sevice EventsService
}

func NewEventsHTTPHandler(service EventsService) *EventsHTTPHandler {
	return &EventsHTTPHandler{
		Sevice: service,
	}
}

type EventsService interface {
	CreateEvent(
		ctx context.Context,
		event domain.Event,
	) (domain.Event, error)

	GetEvent(
		ctx context.Context,
		EventID int,
	) (domain.Event, error)

	PatchEvent(
		ctx context.Context,
		EventID int,
		patch domain.PatchEvent,
	) (domain.Event, error)

	GetEvents(
		ctx context.Context,
		filter domain.EventsFilter,
	) ([]domain.Event, error)

	CancelEvent(
		ctx context.Context,
		eventID int,
	) (domain.Event, error)
}

func (h *EventsHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/events",
			Handler: h.CreateEvent,
		},
		{
			Method:  http.MethodGet,
			Path:    "/events/{id}",
			Handler: h.GetEvent,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/events/{id}",
			Handler: h.PatchEvent,
		},
		{
			Method:  http.MethodGet,
			Path:    "/events",
			Handler: h.GetEvents,
		},
		{
			Method:  http.MethodPost,
			Path:    "/events/{id}/cancel",
			Handler: h.CancelEvent,
		},
	}
}
