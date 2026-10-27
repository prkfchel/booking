package domain

import (
	"fmt"
	"time"

	coreerrors "booking/internal/core/errors"
)

type Booking struct {
	ID         int
	EventID    int
	UserID     int
	CreatedAt  time.Time
	CanceledAt *time.Time
}

func NewBooking(
	id int,
	eventID int,
	userID int,
	createdAt time.Time,
	canceledAt *time.Time,
) Booking {
	return Booking{
		ID:         id,
		EventID:    eventID,
		UserID:     userID,
		CreatedAt:  createdAt,
		CanceledAt: canceledAt,
	}
}

func NewUninitializedBooking(
	eventID int,
	userID int,
) Booking {
	return NewBooking(
		UninitializedID,
		eventID,
		userID,
		UninitializedCreatedAt,
		UninitializedCanceledAt,
	)
}

func (b *Booking) Validate() error {
	if b.EventID < 1 {
		return fmt.Errorf(
			"event id must be positive: %d: %w",
			b.EventID,
			coreerrors.ErrInvalidArgument,
		)
	}

	if b.UserID < 1 {
		return fmt.Errorf(
			"user id must be positive: %d: %w",
			b.UserID,
			coreerrors.ErrInvalidArgument,
		)
	}

	return nil
}

type BookingsFilter struct {
	Pagination

	DateFrom *time.Time
	DateTo   *time.Time
	EventID  *int
	UserID   *int
}

func NewBookingsFilter(
	dateFrom *time.Time,
	dateTo *time.Time,
	eventID *int,
	userID *int,
	limit *int,
	offset *int,
) BookingsFilter {
	return BookingsFilter{
		Pagination: NewPagination(limit, offset),
		DateFrom:   dateFrom,
		DateTo:     dateTo,
		EventID:    eventID,
		UserID:     userID,
	}
}

func (f BookingsFilter) Validate() error {
	if err := f.Pagination.Validate(); err != nil {
		return err
	}

	if err := validateOptionalID("event id", f.EventID); err != nil {
		return err
	}

	if err := validateOptionalID("user id", f.UserID); err != nil {
		return err
	}

	return validateDateRange(f.DateFrom, f.DateTo)
}
