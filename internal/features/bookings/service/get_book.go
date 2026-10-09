package bookings_service

import (
	"booking/internal/core/domain"
	"context"
	"fmt"
)

func (s *BookService) GetBook(
	ctx context.Context,
	bookID int,
) (domain.Booking, error) {
	book, err := s.Repository.GetBook(ctx, bookID)
	if err != nil {
		return domain.Booking{}, fmt.Errorf(
			"failed to get book with id='%d' from repository: %w",
			bookID,
			err,
		)
	}

	return book, nil
}
