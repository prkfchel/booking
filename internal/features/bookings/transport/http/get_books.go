package bookings_http_transport

import (
	"booking/internal/core/domain"
	core_logger "booking/internal/core/logger"
	core_http_request "booking/internal/core/transport/http/request"
	core_http_response "booking/internal/core/transport/http/response"
	"fmt"
	"net/http"
)

type getBooksResponseDTO []BookResponseDTO

// GetEvents     godoc
// @Summary      Список броней
// @Description  Просмотр списка броней с опциональной пагинацией и фильтрацией по пользователю и событию и времяни создания
// @Tags         bookings
// @Produce      json
// @Param        limit  query int false                        "Размер страницы с бронями"
// @Param        offset query int false                        "Смещение страницы с бронями"
// @Param        dateFrom  query time false                    "Начало временного интервала"
// @Param        dateTo query int false                        "Конец временного интервала"
// @Param        userId query int false                        "ID владельца брони"
// @Param        eventId query int false                       "ID события"
// @Success      200 {object} getBooksResponseDTO             "Успешное получение списка броней"
// @Failure      400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure      500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router       /books [get]
func (h *BookHTTPHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseWriter(log, w)

	filter, err := GetBookingFilter(r)
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

func GetBookingFilter(r *http.Request) (
	domain.BookingsFilter,
	error,
) {

	const (
		eventIDQueryParamKey  = "eventId"
		userIDQueryParamKey   = "userId"
		dateFromQueryParamKey = "dateFrom"
		dateToQueryParamKey   = "dateTo"
	)

	limit, offset, err := core_http_request.GetLimitOffsetQueryParams(r)
	if err != nil {
		return domain.BookingsFilter{}, fmt.Errorf(
			"failed to get limit/offset query params: %w",
			err,
		)
	}

	eventId, err := core_http_request.GetIntQueryParams(r, eventIDQueryParamKey)
	if err != nil {
		return domain.BookingsFilter{}, fmt.Errorf(
			"failed to get event id: %w",
			err,
		)
	}

	userId, err := core_http_request.GetIntQueryParams(r, userIDQueryParamKey)
	if err != nil {
		return domain.BookingsFilter{}, fmt.Errorf(
			"failed to get user id: %w",
			err,
		)
	}

	dateFrom, err := core_http_request.GetDateQueryParams(r, dateFromQueryParamKey)
	if err != nil {
		return domain.BookingsFilter{}, fmt.Errorf(
			"failed to get dateFrom: %w",
			err,
		)
	}

	dateTo, err := core_http_request.GetDateQueryParams(r, dateToQueryParamKey)
	if err != nil {
		return domain.BookingsFilter{}, fmt.Errorf(
			"failed to get dateTo: %w",
			err,
		)
	}

	filter := domain.NewBookingsFilter(
		dateFrom,
		dateTo,
		eventId,
		userId,
		limit,
		offset,
	)

	return filter, nil
}
