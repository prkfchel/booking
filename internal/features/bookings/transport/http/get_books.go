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

func (h *BookHTTPHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseWriter(log, w)

	filter, err := GetEventsFilter(r)
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

func GetEventsFilter(r *http.Request) (
	domain.BookingsFilter,
	error,
) {

	const (
		limitQueryParamKey    = "limit"
		offsetQueryParamKey   = "offset"
		eventIDQueryParamKey  = "eventId"
		userIDQueryParamKey   = "userId"
		dateFromQueryParamKey = "dateFrom"
		dateToQueryParamKey   = "dateTo"
	)

	limit, err := core_http_request.GetIntQueryParams(r, limitQueryParamKey)
	if err != nil {
		return domain.BookingsFilter{}, fmt.Errorf(
			"failed to get limit query param: %w",
			err,
		)
	}

	offset, err := core_http_request.GetIntQueryParams(r, offsetQueryParamKey)
	if err != nil {
		return domain.BookingsFilter{}, fmt.Errorf(
			"failed to get offset query param: %w",
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
