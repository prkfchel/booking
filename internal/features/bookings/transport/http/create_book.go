package bookings_http_transport

import (
	"booking/internal/core/domain"
	core_logger "booking/internal/core/logger"
	core_http_request "booking/internal/core/transport/http/request"
	core_http_response "booking/internal/core/transport/http/response"
	"net/http"
)

type createBookRequestDTO struct {
	EventID int `json:"event_id" validate:"required"`
	UserID  int `json:"user_id" validate:"required"`
}

type createBookResponseDTO BookResponseDTO

func (h *BookHTTPHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseWriter(log, w)

	var request createBookRequestDTO

	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate request")
		return
	}

	book := DomainFromDTO(request)

	createdBook, err := h.Service.CreateBook(ctx, book)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create book")
		return
	}

	response := createBookResponseDTO(DTOFromDomain(createdBook))

	responseHandler.JSONResponse(http.StatusCreated, response)
}

func DomainFromDTO(request createBookRequestDTO) domain.Booking {
	return domain.NewUninitializedBooking(
		request.EventID,
		request.UserID,
	)
}
