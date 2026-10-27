package bookings_http_transport

import (
	"booking/internal/core/domain"
	core_logger "booking/internal/core/logger"
	core_http_request "booking/internal/core/transport/http/request"
	core_http_response "booking/internal/core/transport/http/response"
	"net/http"
)

type createBookRequestDTO struct {
	EventID int `json:"event_id" validate:"required" example:"1"`
	UserID  int `json:"user_id" validate:"required" example:"1"`
}

type createBookResponseDTO BookResponseDTO

// CreateBook  godoc
// @Summary      Создать бронь
// @Description  Создать новую бронь в системе
// @Tags         bookings
// @Accept       json
// @Produce      json
// @Param        request body createBookRequestDTO true 		"CreateEvent тело запроса"
// @Success      201 {object} createBookResponseDTO 			"Успешно созданная бронь"
// @Failure      400 {object} core_http_response.ErrorResponse  "Bad request"
// @Failure      500 {object} core_http_response.ErrorResponse  "Internal server error"
// @Router       /books [post]
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
