package bookings_service

import (
	"booking/internal/core/domain"
	"context"
	"fmt"
)

func (s *BookService) CancelBook(
	ctx context.Context,
	bookID int,
) (domain.Booking, error) {
	book, err := s.Repository.CancelBook(ctx, bookID)
	if err != nil {
		return domain.Booking{}, fmt.Errorf(
			"failed to cancel book with id='%d' in repository: %w",
			bookID,
			err,
		)
	}

	return book, nil
}
