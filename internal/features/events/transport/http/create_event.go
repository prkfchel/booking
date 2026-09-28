package events_http_transport

import (
	"booking/internal/core/domain"
	core_logger "booking/internal/core/logger"
	core_http_request "booking/internal/core/transport/http/request"
	core_http_response "booking/internal/core/transport/http/response"
	"net/http"
	"time"
)

type createEventRequestDTO struct {
	Title        string    `json:"title" validate:"required,min=2,max=500"`
	Description  *string   `json:"description" validate:"omitempty,min=2"`
	TotalTickets int       `json:"total_tickets" validate:"required,min=1"`
	EventAddress string    `json:"event_address" validate:"required,min=5"`
	EventTime    time.Time `json:"event_time" validate:"required"`
	OwnersUserId int       `json:"owners_user_id" validate:"required"`
}

type createEventResponseDTO eventResponseDTO

func (h *EventsHTTPHandler) CreaterEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseWriter(log, w)

	var request createEventRequestDTO
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate request")
		return
	}

	domain := DomainFromDTO(request)

	eventDomain, err := h.Sevice.CreateEvent(ctx, domain)
	if err != nil {
		responseHandler.ErrorResponse(err, "create event")
		return
	}

	response := createEventResponseDTO(DTOFromDomain(eventDomain))

	responseHandler.JSONResponse(http.StatusCreated, response)
}

func DomainFromDTO(dto createEventRequestDTO) domain.Event {
	return domain.CreateUninitializedEvent(
		dto.Title,
		dto.Description,
		dto.TotalTickets,
		dto.EventAddress,
		dto.EventTime,
		dto.OwnersUserId,
	)
}
