package bookings_http_transport

import (
	core_errors "booking/internal/core/errors"
	core_http_request "booking/internal/core/transport/http/request"
	"fmt"
	"net/http"
)

func getBookIdPathValue(r *http.Request) (int, error) {
	const bookIdPathValueKey = "id"

	id, err := core_http_request.GetIntPathValuese(r, bookIdPathValueKey)
	if err != nil {
		return 0, fmt.Errorf(
			"faild to get id path value: %w",
			err,
		)
	}

	if id < 1 {
		return 0, fmt.Errorf(
			"book id is must be positive: %v: %w",
			id,
			core_errors.ErrInvalidArgument,
		)
	}

	return id, nil
}
