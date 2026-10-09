package bookings_service

import (
	"booking/internal/core/domain"
	"context"
	"fmt"
)

func (s *BookService) GetBooks(
	ctx context.Context,
	filter domain.BookingsFilter,
) ([]domain.Booking, error) {
	if err := filter.Validate(); err != nil {
		return nil, fmt.Errorf(
			"validate filter error: %w",
			err,
		)
	}

	books, err := s.Repository.GetBooks(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get books from repository: %w",
			err,
		)
	}

	return books, nil
}
