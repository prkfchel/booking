package events_http_transport

import (
	core_logger "booking/internal/core/logger"
	core_http_response "booking/internal/core/transport/http/response"
	"net/http"
)

type cancelEventResponseDTO eventResponseDTO

func (h *EventsHTTPHandler) CancelEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseWriter(log, w)

	id, err := getEventIdPathValue(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get path value event id")
		return
	}

	event, err := h.Sevice.CancelEvent(ctx, id)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to cancel event")
		return
	}

	response := cancelEventResponseDTO(DTOFromDomain(event))

	responseHandler.JSONResponse(http.StatusOK, response)
}
