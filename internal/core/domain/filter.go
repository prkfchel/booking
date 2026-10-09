package domain

import (
	"fmt"
	"time"

	coreerrors "booking/internal/core/errors"
)

// Pagination - общие параметры постраничной выдачи для всех фильтров.
type Pagination struct {
	Limit  *int
	Offset *int
}

func NewPagination(limit, offset *int) Pagination {
	return Pagination{
		Limit:  limit,
		Offset: offset,
	}
}

func (p Pagination) Validate() error {
	if p.Limit != nil && *p.Limit < 1 {
		return fmt.Errorf(
			"limit must be positive: %d: %w",
			*p.Limit,
			coreerrors.ErrInvalidArgument,
		)
	}

	if p.Offset != nil && *p.Offset < 0 {
		return fmt.Errorf(
			"offset must be non-negative: %d: %w",
			*p.Offset,
			coreerrors.ErrInvalidArgument,
		)
	}

	return nil
}

// validateDateRange проверяет, что from строго раньше to (если заданы оба).
func validateDateRange(from, to *time.Time) error {
	if from == nil || to == nil {
		return nil
	}

	if !from.Before(*to) {
		return fmt.Errorf(
			"dateFrom=%s must be before dateTo=%s: %w",
			from.Format(time.RFC3339),
			to.Format(time.RFC3339),
			coreerrors.ErrInvalidArgument,
		)
	}

	return nil
}

// validateOptionalID проверяет опциональный ID: если задан, то должен быть >= 1.
func validateOptionalID(name string, id *int) error {
	if id != nil && *id < 1 {
		return fmt.Errorf(
			"%s must be positive: %d: %w",
			name,
			*id,
			coreerrors.ErrInvalidArgument,
		)
	}

	return nil
}
