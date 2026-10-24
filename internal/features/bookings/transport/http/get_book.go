package bookings_http_transport

import (
	core_logger "booking/internal/core/logger"
	core_http_response "booking/internal/core/transport/http/response"
	"net/http"
)

type getBookResponseDTO BookResponseDTO

// GetBook      godoc
// @Summary      Получение брони
// @Description  Получение конкретной брони по ее ID
// @Tags         bookings
// @Produce      json
// @Param        id  path string true                          	"ID получаемой брони" Format(uuid)
// @Success      200 {object} getBookResponseDTO                "Бронь успешно найдена"
// @Failure      400 {object} core_http_response.ErrorResponse  "Bad request"
// @Failure      404 {object} core_http_response.ErrorResponse 	"Event not found"
// @Failure      500 {object} core_http_response.ErrorResponse 	"Internal server error"
// @Router       /books/{id} [get]
func (h *BookHTTPHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseWriter(log, w)

	id, err := getBookIDPathValue(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get book  id")
		return
	}

	book, err := h.Service.GetBook(ctx, id)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get book")
		return
	}

	response := getBookResponseDTO(DTOFromDomain(book))

	responseHandler.JSONResponse(http.StatusOK, response)
}
