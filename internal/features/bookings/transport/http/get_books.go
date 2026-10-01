package bookings_http_transport

import (
	"booking/internal/core/domain"
	core_logger "booking/internal/core/logger"
	core_http_response "booking/internal/core/transport/http/response"
	"net/http"
)

type getBooksResponseDTO []BookResponseDTO

func (h *BookHTTPHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseWriter(log, w)

	filter, err := domain.GetEventsFilter(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get query params")
		return
	}

	books, err := h.Service.GetBooks(ctx, filter)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get books")
		return
	}

	response := DTOsFromDomains(books)

	responseHandler.JSONResponse(http.StatusOK, response)
}
