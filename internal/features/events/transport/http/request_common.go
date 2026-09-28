package events_http_transport

import (
	core_errors "booking/internal/core/errors"
	core_http_request "booking/internal/core/transport/http/request"
	"fmt"
	"net/http"
)

func GetLimitOffsetQueryParams(r *http.Request) (*int, *int, error) {

	const (
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)

	limit, err := core_http_request.GetIntQueryParams(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"failed to get limit query param: %w",
			err,
		)
	}

	offset, err := core_http_request.GetIntQueryParams(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"failed to get offset query param: %w",
			err,
		)
	}

	if limit != nil && *limit < 1 {
		return nil, nil, fmt.Errorf(
			"limit query param must be positive %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if offset != nil && *offset < 0 {
		return nil, nil, fmt.Errorf(
			"offset query param cant be negative: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	return limit, offset, nil

}

func getEventIdPathValue(r *http.Request) (int, error) {
	const eventIdPathValueKey = "id"

	id, err := core_http_request.GetIntPathValuese(r, eventIdPathValueKey)
	if err != nil {
		return 0, fmt.Errorf(
			"faild to get id path value: %w",
			err,
		)
	}

	if id < 1 {
		return 0, fmt.Errorf(
			"event id is must be positive: %v: %w",
			id,
			core_errors.ErrInvalidArgument,
		)
	}

	return id, nil
}
