package events_http_transport

import (
	core_logger "booking/internal/core/logger"
	core_http_response "booking/internal/core/transport/http/response"
	"net/http"
)

type getEventReaponseDTO eventResponseDTO

func (h *EventsHTTPHandler) GetEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseWriter(log, w)

	id, err := getEventIdPathValue(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get event id")
		return
	}

	event, err := h.Sevice.GetEvent(ctx, id)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get event")
		return
	}

	response := getEventReaponseDTO(DTOFromDomain(event))

	responseHandler.JSONResponse(http.StatusOK, response)
}
