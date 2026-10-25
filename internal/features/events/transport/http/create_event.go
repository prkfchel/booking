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
	Title        string    `json:"title" validate:"required,min=2,max=500" example:"concert"`
	Description  *string   `json:"description" validate:"omitempty,min=2" example:"interesting concert and show"`
	TotalTickets int       `json:"total_tickets" validate:"required,min=1" example:"100"`
	Address      string    `json:"address" validate:"required,min=5" example:"st. Pushkina, b. Kalatushkina"`
	StartsAt     time.Time `json:"starts_at" validate:"required" example:"2026-11-10T19:30:00Z"`
	OwnerUserID  int       `json:"owner_user_id" validate:"required" example:"1"`
}

type createEventResponseDTO EventResponseDTO

// CreaterEvent  godoc
// @Summary      Создать событие
// @Description  Создать новое событие в системе
// @Tags         events
// @Accept       json
// @Produce      json
// @Param        request body createEventRequestDTO true 		"CreateEvent тело запроса"
// @Success      201 {object} createEventResponseDTO 			"Успешно созданное событие"
// @Failure      400 {object} core_http_response.ErrorResponse  "Bad request"
// @Failure      500 {object} core_http_response.ErrorResponse  "Internal server error"
// @Router       /events [post]
func (h *EventsHTTPHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseWriter(log, w)

	var request createEventRequestDTO
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate request")
		return
	}

	event := DomainFromDTO(request)

	createdEvent, err := h.Sevice.CreateEvent(ctx, event)
	if err != nil {
		responseHandler.ErrorResponse(err, "create event")
		return
	}

	response := createEventResponseDTO(DTOFromDomain(createdEvent))

	responseHandler.JSONResponse(http.StatusCreated, response)
}

func DomainFromDTO(dto createEventRequestDTO) domain.Event {
	return domain.NewUninitializedEvent(
		dto.Title,
		dto.Description,
		dto.TotalTickets,
		dto.Address,
		dto.StartsAt,
		dto.OwnerUserID,
	)
}
