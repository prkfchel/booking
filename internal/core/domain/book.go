package domain

import (
	core_errors "booking/internal/core/errors"
	core_http_request "booking/internal/core/transport/http/request"
	"fmt"
	"net/http"
	"time"
)

type Book struct {
	ID         int
	EventID    int
	UserID     int
	CreatedAt  time.Time
	Canceled   bool
	CanceledAt *time.Time
}

func NewBook(
	id int,
	eventID int,
	userID int,
	createdAt time.Time,
	canceled bool,
	canceledAt *time.Time,
) Book {
	return Book{
		ID:         id,
		EventID:    eventID,
		UserID:     userID,
		CreatedAt:  createdAt,
		Canceled:   canceled,
		CanceledAt: canceledAt,
	}
}

func CreateUninitializedBook(
	eventID int,
	userID int,
) Book {
	return NewBook(
		UninitializedID,
		eventID,
		userID,
		UninitializedCreatedAt,
		UninitializedCanceled,
		UninitializedCanceledAt,
	)
}

func (book *Book) Validate() error {
	if book.EventID < 1 {
		return fmt.Errorf(
			"event id must be positive: %d: %w",
			book.EventID,
			core_errors.ErrInvalidArgument,
		)
	}

	if book.UserID < 1 {
		return fmt.Errorf(
			"user id must be positive: %d: %w",
			book.UserID,
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}

type BooksQueryFilter struct {
	DateFrom *time.Time
	DateTo   *time.Time
	EventId  *int
	UserId   *int
	Limit    *int
	Offset   *int
}

func NewBooksQueryFilter(
	dateFrom *time.Time,
	dateTo *time.Time,
	eventId *int,
	userId *int,
	limit *int,
	offset *int,
) BooksQueryFilter {
	return BooksQueryFilter{
		DateFrom: dateFrom,
		DateTo:   dateTo,
		EventId:  eventId,
		UserId:   userId,
		Limit:    limit,
		Offset:   offset,
	}
}

func (filter BooksQueryFilter) Validate() error {
	if filter.EventId != nil {
		if *filter.EventId < 1 {
			return fmt.Errorf(
				"event id must be positive: %w",
				core_errors.ErrInvalidArgument,
			)
		}

	}

	if filter.UserId != nil {
		if *filter.UserId < 1 {
			return fmt.Errorf(
				"user id must be positive: %w",
				core_errors.ErrInvalidArgument,
			)
		}

	}

	if filter.Limit != nil {
		if *filter.Limit < 1 {
			return fmt.Errorf(
				"limit must be positive: %d: %w",
				*filter.Limit,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if filter.Offset != nil {
		if *filter.Offset < 1 {
			return fmt.Errorf(
				"offset must be positive: %d: %w",
				*filter.Offset,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if filter.DateFrom != nil && filter.DateTo != nil {
		dateFrom := *filter.DateFrom
		dateTo := *filter.DateTo
		if !dateFrom.Before(dateTo) {
			return fmt.Errorf(
				"dateFrom='%d' must be bofare dateTo='%d': %w",
				dateFrom,
				dateTo,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	return nil
}

func GetEventsFilter(r *http.Request) (
	BooksQueryFilter,
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
		return BooksQueryFilter{}, fmt.Errorf(
			"failed to get limit query param: %w",
			err,
		)
	}

	offset, err := core_http_request.GetIntQueryParams(r, offsetQueryParamKey)
	if err != nil {
		return BooksQueryFilter{}, fmt.Errorf(
			"failed to get offset query param: %w",
			err,
		)
	}

	eventId, err := core_http_request.GetIntQueryParams(r, eventIDQueryParamKey)
	if err != nil {
		return BooksQueryFilter{}, fmt.Errorf(
			"failed to get event id: %w",
			err,
		)
	}

	userId, err := core_http_request.GetIntQueryParams(r, userIDQueryParamKey)
	if err != nil {
		return BooksQueryFilter{}, fmt.Errorf(
			"failed to get user id: %w",
			err,
		)
	}

	dateFrom, err := core_http_request.GetDateQueryParams(r, dateFromQueryParamKey)
	if err != nil {
		return BooksQueryFilter{}, fmt.Errorf(
			"failed to get dateFrom: %w",
			err,
		)
	}

	dateTo, err := core_http_request.GetDateQueryParams(r, dateToQueryParamKey)
	if err != nil {
		return BooksQueryFilter{}, fmt.Errorf(
			"failed to get dateTo: %w",
			err,
		)
	}

	filter := NewBooksQueryFilter(
		dateFrom,
		dateTo,
		eventId,
		userId,
		limit,
		offset,
	)

	return filter, nil
}
