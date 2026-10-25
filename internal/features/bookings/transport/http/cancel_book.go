package bookings_http_transport

import (
	core_logger "booking/internal/core/logger"
	core_http_response "booking/internal/core/transport/http/response"
	"net/http"
)

type cancelBookResponseDTO BookResponseDTO

// CancelBook   godoc
// @Summary      Отменить бронь
// @Description  Отменить бронь в системе
// @Tags         bookings
// @Produce      json
// @Param        id  path string true                          	"ID отменяемой брони" Format(uuid)
// @Success      201 {object} cancelBookResponseDTO 			"Бронь успешно отменена"
// @Failure      400 {object} core_http_response.ErrorResponse  "Bad request"
// @Failure      500 {object} core_http_response.ErrorResponse  "Internal server error"
// @Router       /books/cancel [post]
func (h *BookHTTPHandler) CancelBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseWriter(log, w)

	id, err := getBookIDPathValue(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get book id")
		return
	}

	book, err := h.Service.CancelBook(ctx, id)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to cancel book")
		return
	}

	response := cancelBookResponseDTO(DTOFromDomain(book))

	responseHandler.JSONResponse(http.StatusOK, response)
}
