package events_http_transport

import (
	"booking/internal/core/domain"
	core_errors "booking/internal/core/errors"
	core_logger "booking/internal/core/logger"
	core_http_request "booking/internal/core/transport/http/request"
	core_http_response "booking/internal/core/transport/http/response"
	core_http_type "booking/internal/core/transport/http/type"
	"fmt"
	"net/http"
	"time"
)

type requestPatchEvent struct {
	Title        core_http_type.Nullable[string]    `json:"title"`
	Description  core_http_type.Nullable[string]    `json:"description"`
	TotalTickets core_http_type.Nullable[int]       `json:"total_tickets"`
	EventAddress core_http_type.Nullable[string]    `json:"event_address"`
	EventTime    core_http_type.Nullable[time.Time] `json:"event_time"`
	Canceled     core_http_type.Nullable[bool]      `json:"canceled"`
}

type patchEventResponseDTO eventResponseDTO

func (h *EventsHTTPHandler) PatchEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseWriter(log, w)

	id, err := getEventIdPathValue(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get event id path value")
		return
	}

	var request requestPatchEvent
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate")
		return
	}

	patchDomain := PatchEventDoaminFromDTO(request)

	eventDomain, err := h.Sevice.PatchEvent(ctx, id, patchDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch event")
		return
	}

	response := patchEventResponseDTO(DTOFromDomain(eventDomain))

	responseHandler.JSONResponse(http.StatusOK, response)

}

func PatchEventDoaminFromDTO(request requestPatchEvent) domain.PatchEvent {
	return domain.NewPatchEvent(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.TotalTickets.ToDomain(),
		request.EventAddress.ToDomain(),
		request.EventTime.ToDomain(),
		request.Canceled.ToDomain(),
	)
}

func (request *requestPatchEvent) Validate() error {
	if request.Title.Set {
		titleLen := len([]rune(*request.Title.Value))

		if titleLen < 2 || titleLen > 500 {
			return fmt.Errorf(
				"title len must be between 2 and 500: %d: %w",
				titleLen,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if request.Description.Set {
		if request.Description.Value != nil {
			descriptionLen := len([]rune(*request.Description.Value))
			if descriptionLen < 2 {
				return fmt.Errorf(
					"description len must be more then 2: %d: %w",
					descriptionLen,
					core_errors.ErrInvalidArgument,
				)
			}
		}
	}

	if request.TotalTickets.Set {
		if *request.TotalTickets.Value < 1 {
			return fmt.Errorf(
				"total tickets must be positive number: %d: %w",
				*request.TotalTickets.Value,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if request.EventAddress.Set {
		eventAddresLen := len([]rune(*request.EventAddress.Value))

		if eventAddresLen < 5 {
			return fmt.Errorf(
				"event address len must be more then 5: %d: %w",
				eventAddresLen,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if request.EventTime.Set {
		if request.EventTime.Value.IsZero() {
			return fmt.Errorf(
				"evet time must be: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	return nil
}
