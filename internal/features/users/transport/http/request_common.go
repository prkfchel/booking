package users_http_transport

import (
	core_errors "booking/internal/core/errors"
	core_http_request "booking/internal/core/transport/http/request"
	"fmt"
	"net/http"
)

func getUserIDPathValue(r *http.Request) (int, error) {
	const userIdPathValueKey = "id"

	id, err := core_http_request.GetIntPathValues(r, userIdPathValueKey)
	if err != nil {
		return 0, fmt.Errorf(
			"failed to get id path value: %w",
			err,
		)
	}

	if id < 1 {
		return 0, fmt.Errorf(
			"user id must be positive: %v: %w",
			id,
			core_errors.ErrInvalidArgument,
		)
	}

	return id, nil
}
