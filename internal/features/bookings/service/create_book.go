package bookings_service

import (
	"booking/internal/core/domain"
	"context"
	"fmt"
)

func (s *BookService) CreateBook(
	ctx context.Context,
	book domain.Booking,
) (domain.Booking, error) {
	if err := book.Validate(); err != nil {
		return domain.Booking{}, fmt.Errorf(
			"validation error: %w",
			err,
		)
	}

	createdBook, err := s.Repository.CreateBook(ctx, book)

	if err != nil {
		return domain.Booking{}, fmt.Errorf(
			"failed to cvreate book: %w",
			err,
		)
	}

	return createdBook, nil
}
