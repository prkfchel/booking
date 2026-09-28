package events_http_transport

import (
	"booking/internal/core/domain"
	core_logger "booking/internal/core/logger"
	core_http_request "booking/internal/core/transport/http/request"
	core_http_response "booking/internal/core/transport/http/response"
	"fmt"
	"net/http"
)

type GetEventsReaponseDTO []eventResponseDTO

func (h *EventsHTTPHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseWriter(log, w)

	filter, err := GetEventsFilter(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get filter params")
		return
	}

	events, err := h.Sevice.GetEvents(ctx, filter)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get events")
		return
	}

	response := GetEventsReaponseDTO(DTOFromDomains(events))

	responseHandler.JSONResponse(http.StatusOK, response)
}

func GetEventsFilter(r *http.Request) (
	domain.EventsQueryFilter,
	error,
) {

	const (
		limitQueryParamKey    = "limit"
		offsetQueryParamKey   = "offset"
		titleQueryParamKey    = "title"
		addressQueryParamKey  = "address"
		dateFromQueryParamKey = "dateFrom"
		dateToQueryParamKey   = "dateTo"
		ownersIdQueryParamKey = "ownersId"
	)

	limit, err := core_http_request.GetIntQueryParams(r, limitQueryParamKey)
	if err != nil {
		return domain.EventsQueryFilter{}, fmt.Errorf(
			"failed to get limit query param: %w",
			err,
		)
	}

	offset, err := core_http_request.GetIntQueryParams(r, offsetQueryParamKey)
	if err != nil {
		return domain.EventsQueryFilter{}, fmt.Errorf(
			"failed to get offset query param: %w",
			err,
		)
	}

	ownersUserId, err := core_http_request.GetIntQueryParams(r, ownersIdQueryParamKey)
	if err != nil {
		return domain.EventsQueryFilter{}, fmt.Errorf(
			"failed to get owners user id: %w",
			err,
		)
	}

	title := core_http_request.GetStringQueryParams(r, titleQueryParamKey)

	address := core_http_request.GetStringQueryParams(r, addressQueryParamKey)

	dateFrom, err := core_http_request.GetDateQueryParams(r, dateFromQueryParamKey)
	if err != nil {
		return domain.EventsQueryFilter{}, fmt.Errorf(
			"failed to get dateFrom: %w",
			err,
		)
	}

	dateTo, err := core_http_request.GetDateQueryParams(r, dateToQueryParamKey)
	if err != nil {
		return domain.EventsQueryFilter{}, fmt.Errorf(
			"failed to get dateTo: %w",
			err,
		)
	}

	filter := domain.NewEventsQueryFilter(
		title,
		address,
		dateFrom,
		dateTo,
		ownersUserId,
		limit,
		offset,
	)

	return filter, nil
}
