package core_http_request

import (
	core_errors "booking/internal/core/errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func GetIntQueryParams(r *http.Request, key string) (*int, error) {
	param := r.URL.Query().Get(key)

	if param == "" {
		return nil, nil
	}

	val, err := strconv.Atoi(param)
	if err != nil {
		return nil, fmt.Errorf(
			"param='%s' by key='%s' not a valid integer: %v: %w",
			param,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return &val, nil
}

func GetStringQueryParams(r *http.Request, key string) *string {
	param := r.URL.Query().Get(key)

	if param == "" {
		return nil
	}

	return &param
}

func GetDateQueryParams(r *http.Request, key string) (*time.Time, error) {
	param := r.URL.Query().Get(key)

	if param == "" {
		return nil, nil
	}

	layout := "2006-01-02"

	date, err := time.Parse(layout, param)
	if err != nil {
		return nil, fmt.Errorf(
			"param='%s' by key='%s' is not valid date: %v: %w",
			param,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return &date, nil
}

func GetLimitOffsetQueryParams(r *http.Request) (*int, *int, error) {

	const (
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)

	limit, err := GetIntQueryParams(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"failed to get limit query param: %w",
			err,
		)
	}

	offset, err := GetIntQueryParams(r, offsetQueryParamKey)
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
