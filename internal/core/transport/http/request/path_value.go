package core_http_request

import (
	core_errors "booking/internal/core/errors"
	"fmt"
	"net/http"
	"strconv"
)

func GetIntPathValuese(r *http.Request, key string) (int, error) {
	pathValue := r.PathValue(key)
	if pathValue == "" {
		return 0, fmt.Errorf(
			"no key='%s' in path values: %w",
			key,
			core_errors.ErrInvalidArgument,
		)
	}

	val, err := strconv.Atoi(pathValue)
	if err != nil {
		return 0, fmt.Errorf(
			"path value='%s' not a valid integer: %v: %w",
			pathValue,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return val, nil
}
