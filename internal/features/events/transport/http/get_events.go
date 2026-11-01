package events_http_transport

import (
	"booking/internal/core/domain"
	core_logger "booking/internal/core/logger"
	core_http_request "booking/internal/core/transport/http/request"
	core_http_response "booking/internal/core/transport/http/response"
	"fmt"
	"net/http"
)

type getEventsResponseDTO []EventResponseDTO

// GetEvents     godoc
// @Summary      Список событий
// @Description  Просмотр списка событий с опциональной пагинацией и фильтрацией по названию и адресу, времяни проведениря и id организатора
// @Tags         events
// @Produce      json
// @Param        limit  	query int false                    "Размер страницы с событиями"
// @Param        offset 	query int false                    "Смещение страницы с событиями"
// @Param        dateFrom 	query string false                 "Начало временного интервала"
// @Param        dateTo 	query string false                 "Конец временного интервала"
// @Param        ownerId 	query int false                    "ID организатора мероприятия"
// @Param        title		query string false				   "Название мероприятия"
// @Param        address	query string false				   "Адрес мероприятия"
// @Success      200 {object} getEventsResponseDTO             "Успешное получение списка событий"
// @Failure      400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure      500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router       /events [get]
func (h *EventsHTTPHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseWriter(log, w)

	filter, err := GetEventsFilter(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get filter params")
		return
	}

	events, err := h.Service.GetEvents(ctx, filter)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get events")
		return
	}

	response := getEventsResponseDTO(DTOFromDomains(events))

	responseHandler.JSONResponse(http.StatusOK, response)
}

func GetEventsFilter(r *http.Request) (
	domain.EventsFilter,
	error,
) {

	const (
		dateFromQueryParamKey = "dateFrom"
		dateToQueryParamKey   = "dateTo"
		ownerIdQueryParamKey  = "ownerId"
		titleQueryParamKey    = "title"
		addressQueryParamKey  = "address"
	)

	limit, offset, err := core_http_request.GetLimitOffsetQueryParams(r)
	if err != nil {
		return domain.EventsFilter{}, fmt.Errorf(
			"failed to get limit/offset query params: %w",
			err,
		)
	}

	ownersUserId, err := core_http_request.GetIntQueryParams(r, ownerIdQueryParamKey)
	if err != nil {
		return domain.EventsFilter{}, fmt.Errorf(
			"failed to get owners user id: %w",
			err,
		)
	}

	dateFrom, err := core_http_request.GetDateQueryParams(r, dateFromQueryParamKey)
	if err != nil {
		return domain.EventsFilter{}, fmt.Errorf(
			"failed to get dateFrom: %w",
			err,
		)
	}

	dateTo, err := core_http_request.GetDateQueryParams(r, dateToQueryParamKey)
	if err != nil {
		return domain.EventsFilter{}, fmt.Errorf(
			"failed to get dateTo: %w",
			err,
		)
	}

	title := core_http_request.GetStringQueryParams(r, titleQueryParamKey)

	address := core_http_request.GetStringQueryParams(r, addressQueryParamKey)

	filter := domain.NewEventsFilter(
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
