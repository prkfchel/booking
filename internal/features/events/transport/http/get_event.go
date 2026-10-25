package events_http_transport

import (
	core_logger "booking/internal/core/logger"
	core_http_response "booking/internal/core/transport/http/response"
	"net/http"
)

type getEventReaponseDTO EventResponseDTO

// GetEvent      godoc
// @Summary      Получение события
// @Description  Получение конкретного события по его ID
// @Tags         events
// @Produce      json
// @Param        id  path string true                          	"ID получаемого события" Format(uuid)
// @Success      200 {object} getEventReaponseDTO                "событие успешно найдено"
// @Failure      400 {object} core_http_response.ErrorResponse  "Bad request"
// @Failure      404 {object} core_http_response.ErrorResponse 	"Event not found"
// @Failure      500 {object} core_http_response.ErrorResponse 	"Internal server error"
// @Router       /events/{id} [get]
func (h *EventsHTTPHandler) GetEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseWriter(log, w)

	id, err := getEventIDPathValue(r)
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
