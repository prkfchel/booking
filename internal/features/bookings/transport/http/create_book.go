package bookings_http_transport

import (
	"booking/internal/core/domain"
	core_logger "booking/internal/core/logger"
	core_http_request "booking/internal/core/transport/http/request"
	core_http_response "booking/internal/core/transport/http/response"
	"net/http"
)

type CreateBookRequestDTO struct {
	EventID int `json:"event_id"`
	UserID  int `json:"user_id"`
}

type CreateBookResponseDTO BookResponseDTO

func (h *BookHTTPHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseWriter(log, w)

	var request CreateBookRequestDTO

	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate request")
		return
	}

	uninitializedBook := DomainFromDTO(request)

	book, err := h.Service.CreateBook(ctx, uninitializedBook)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create book")
		return
	}

	response := CreateBookResponseDTO(DTOFromDomain(book))

	responseHandler.JSONResponse(http.StatusCreated, response)
}

func DomainFromDTO(request CreateBookRequestDTO) domain.Book {
	return domain.CreateUninitializedBook(
		request.EventID,
		request.UserID,
	)
}
