package bookings_http_transport

import (
	core_logger "booking/internal/core/logger"
	core_http_response "booking/internal/core/transport/http/response"
	"net/http"
)

type getBookResponseDTO BookResponseDTO

func (h *BookHTTPHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseWriter(log, w)

	id, err := getBookIdPathValue(r)
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
