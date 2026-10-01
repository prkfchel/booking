package bookings_service

import (
	"booking/internal/core/domain"
	"context"
	"fmt"
)

func (s *BookService) CancelBook(
	ctx context.Context,
	id int,
) (domain.Book, error) {
	if _, err := s.Repository.GetBook(ctx, id); err != nil {
		return domain.Book{}, fmt.Errorf(
			"failed to get book with id='%d' from repository: %w",
			id,
			err,
		)
	}

	book, err := s.Repository.CancelBook(ctx, id)
	if err != nil {
		return domain.Book{}, fmt.Errorf(
			"failed to cancel book with id='%d' in repository: %w",
			id,
			err,
		)
	}

	return book, nil
}
